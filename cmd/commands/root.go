package commands

import (
	"fmt"
	"os"
	"path/filepath"
	
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tmos "github.com/tendermint/tendermint/libs/os"
	
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	tmflags "github.com/tendermint/tendermint/libs/cli/flags"
	"github.com/tendermint/tendermint/libs/log"
	
	dmlog "github.com/rhizome-chain/tendermint-daemon/tm/log"
)

var (
	config = cfg.DefaultConfig()
	logger = dmlog.NewTMLogger(log.NewSyncWriter(os.Stdout))
)

func init() {
	registerFlagsRootCmd(RootCmd)
}

func registerFlagsRootCmd(cmd *cobra.Command) {
	cmd.PersistentFlags().String("log_level", config.LogLevel, "Log level")
}

// ParseConfig retrieves the default environment configuration,
// sets up the Tendermint root and ensures that the root exists
func ParseConfig() (*cfg.Config, error) {
	conf := cfg.DefaultConfig()
	err := viper.Unmarshal(conf)
	if err != nil {
		return nil, err
	}
	conf.SetRoot(conf.RootDir)
	EnsureRoot(conf.RootDir)
	if err = conf.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("error in config file: %v", err)
	}
	return conf, err
}

// RootCmd is the root command for Tendermint core.
var RootCmd = &cobra.Command{
	Use:   "bcb",
	Short: "Rhizome-Chain on Tendermint Core (BFT Consensus)",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if cmd.Name() == VersionCmd.Name() {
			return nil
		}
		config, err = ParseConfig()
		if err != nil {
			return err
		}
		if config.LogFormat == cfg.LogFormatJSON {
			logger = log.NewTMJSONLogger(log.NewSyncWriter(os.Stdout))
		}
		logger, err = tmflags.ParseLogLevel(config.LogLevel, logger, cfg.DefaultLogLevel())
		if err != nil {
			return err
		}
		if viper.GetBool(cli.TraceFlag) {
			logger = log.NewTracingLogger(logger)
		}
		logger = logger.With("module", "main")
		return nil
	},
}

func InitRootCommand() *cobra.Command {
	rootCmd := RootCmd
	rootCmd.AddCommand(
		ResetAllCmd,
		ResetPrivValidatorCmd,
		ShowValidatorCmd,
		ShowNodeIDCmd,
		VersionCmd,
	)
	return rootCmd
}

func EnsureRoot(rootDir string) {
	if err := tmos.EnsureDir(rootDir, cfg.DefaultDirPerm); err != nil {
		panic(err.Error())
	}
	if err := tmos.EnsureDir(filepath.Join(rootDir, "config"), cfg.DefaultDirPerm); err != nil {
		panic(err.Error())
	}
	if err := tmos.EnsureDir(filepath.Join(rootDir, "data"), cfg.DefaultDirPerm); err != nil {
		panic(err.Error())
	}
}
