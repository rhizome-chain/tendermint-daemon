package tm

import (
	"fmt"
	"os"
	
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
)

type Provider func(*cfg.Config, log.Logger) (tmNode *node.Node, dapp *DaemonApp, err error)

type NodeProvider struct {
}

func (provider *NodeProvider) NewNode(config *cfg.Config, logger log.Logger) (tmNode *node.Node, dapp *DaemonApp, err error) {
	// Generate node PrivKey
	nodeKey, err := p2p.LoadOrGenNodeKey(config.NodeKeyFile())
	if err != nil {
		return nil, nil, err
	}
	
	// Convert old PrivValidator if it exists.
	oldPrivVal := config.OldPrivValidatorFile()
	newPrivValKey := config.PrivValidatorKeyFile()
	newPrivValState := config.PrivValidatorStateFile()
	if _, err := os.Stat(oldPrivVal); !os.IsNotExist(err) {
		oldPV, err := privval.LoadOldFilePV(oldPrivVal)
		if err != nil {
			return nil, nil, fmt.Errorf("error reading OldPrivValidator from %v: %v", oldPrivVal, err)
		}
		logger.Info("Upgrading PrivValidator file",
			"old", oldPrivVal,
			"newKey", newPrivValKey,
			"newState", newPrivValState,
		)
		oldPV.Upgrade(newPrivValKey, newPrivValState)
	}
	
	dapp = NewDaemonApplication(config, logger)
	
	tmNode, err = node.NewNode(config,
		privval.LoadOrGenFilePV(newPrivValKey, newPrivValState),
		nodeKey,
		proxy.NewLocalClientCreator(dapp),
		node.DefaultGenesisDocProviderFunc(config),
		node.DefaultDBProvider,
		node.DefaultMetricsProvider(config.Instrumentation),
		logger,
	)
	
	return tmNode, dapp, err
}
