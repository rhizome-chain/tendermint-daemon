package commands

import (
	"fmt"
	"path/filepath"
	
	"github.com/rhizome-chain/tendermint-daemon/daemon/common"
	
	"github.com/spf13/viper"
	
	"github.com/spf13/cobra"
	cfg "github.com/tendermint/tendermint/config"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
	
	"github.com/rhizome-chain/tendermint-daemon/daemon"
)

// InitFilesCmd initialises a fresh Tendermint Core instance.

func AddInitCommand(cmd *cobra.Command, daemonProvider *daemon.BaseProvider) {
	// Create Init
	cmd.AddCommand(NewInitCmd(daemonProvider))
}

func AddFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("force-rewrite", false, "force rewrite")
	
	cmd.Flags().String("chain-id", config.ChainID(), "Chain ID")
	cmd.Flags().String("node-name", config.Moniker, "Node Name")
	cmd.Flags().String("rpc.laddr", config.RPC.ListenAddress, "rpc laddr")
	cmd.Flags().Bool("p2p.seed_mode", config.P2P.SeedMode, "p2p seed_mode")
	cmd.Flags().String("p2p.persistent_peers", config.P2P.PersistentPeers, "p2p persistent_peers")
	cmd.Flags().String("p2p.seeds", config.P2P.Seeds, "p2p seeds")
	cmd.Flags().String("p2p.laddr", config.P2P.ListenAddress, "p2p laddr")
	cmd.Flags().Bool("p2p.allow_duplicate_ip", config.P2P.AllowDuplicateIP, "p2p persistent_peers")
	cmd.Flags().Int("mempool.size", config.Mempool.Size, "mempool.size")
	cmd.Flags().Int64("mempool.max_txs_bytes", config.Mempool.MaxTxsBytes, "mempool.max_txs_bytes")
	cmd.Flags().Int("mempool.max_tx_bytes", config.Mempool.MaxTxBytes, "mempool.max_tx_bytes")
	cmd.Flags().String("consensus.timeout_commit", "1s", "consensus timeout_commit")
	cmd.Flags().String("instrumentation.prometheus_listen_addr", config.Instrumentation.PrometheusListenAddr, "instrumentation.prometheus_listen_addr")
	
}

func NewInitCmd(daemonProvider *daemon.BaseProvider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize Dist-Daemon on Tendermint",
		RunE: func(cmd *cobra.Command, args []string) error {
			return initFilesWithConfig(config, daemonProvider)
		},
	}
	AddFlags(cmd)
	daemonProvider.AddFlags(cmd)
	return cmd
}

func initFilesWithConfig(config *cfg.Config, daemonProvider *daemon.BaseProvider) error {
	// private validator
	privValKeyFile := config.PrivValidatorKeyFile()
	privValStateFile := config.PrivValidatorStateFile()
	var pv *privval.FilePV
	if tmos.FileExists(privValKeyFile) {
		pv = privval.LoadFilePV(privValKeyFile, privValStateFile)
		logger.Info("Found private validator", "keyFile", privValKeyFile,
			"stateFile", privValStateFile)
	} else {
		pv = privval.GenFilePV(privValKeyFile, privValStateFile)
		pv.Save()
		logger.Info("Generated private validator", "keyFile", privValKeyFile,
			"stateFile", privValStateFile)
	}
	
	nodeKeyFile := config.NodeKeyFile()
	if tmos.FileExists(nodeKeyFile) {
		logger.Info("Found node key", "path", nodeKeyFile)
	} else {
		if _, err := p2p.LoadOrGenNodeKey(nodeKeyFile); err != nil {
			return err
		}
		logger.Info("Generated node key", "path", nodeKeyFile)
	}
	
	// genesis file
	genFile := config.GenesisFile()
	if tmos.FileExists(genFile) {
		logger.Info("Found genesis file", "path", genFile)
	} else {
		chainID := viper.GetString("chain-id")
		
		if len(chainID) == 0 {
			chainID = fmt.Sprintf("test-chain-%v", tmrand.Str(6))
		}
		genDoc := types.GenesisDoc{
			ChainID:         chainID,
			GenesisTime:     tmtime.Now(),
			ConsensusParams: types.DefaultConsensusParams(),
		}
		key := pv.GetPubKey()
		genDoc.Validators = []types.GenesisValidator{{
			Address: key.Address(),
			PubKey:  key,
			Power:   10,
		}}
		
		if err := genDoc.SaveAs(genFile); err != nil {
			return err
		}
		logger.Info("Generated genesis file", "path", genFile)
	}
	
	confFilePath := filepath.Join(config.RootDir, "config", "config.toml")
	
	if viper.GetBool("force-rewrite") || !tmos.FileExists(confFilePath) {
		nodeName := viper.GetString("node-name")
		config.BaseConfig.Moniker = nodeName
		cfg.WriteConfigFile(confFilePath, config)
		logger.Info("Write config.toml ", "path", confFilePath)
	}
	
	daemonConfig := &common.DaemonConfig{
		ChainID:  config.ChainID(),
		NodeName: config.Moniker,
	}
	
	viper.Unmarshal(daemonConfig)
	
	daemonProvider.InitFiles(config, daemonConfig)
	return nil
}
