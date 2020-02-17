package daemon

import (
	"github.com/rhizome-chain/tendermint-daemon/daemon/common"
	"github.com/rhizome-chain/tendermint-daemon/tm"
	"github.com/spf13/cobra"
	
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
)

type Provider func(tmCfg *cfg.Config, logger log.Logger, tmNode *node.Node, daemonApp *tm.DaemonApp, config common.DaemonConfig) *Daemon

type BaseProvider struct {
	modules []Module
}

func (provider *BaseProvider) AddModule(module Module) {
	if provider.modules == nil {
		provider.modules = []Module{}
	}
	provider.modules = append(provider.modules, module)
}

func (provider *BaseProvider) AddFlags(cmd *cobra.Command) {
	if provider.modules != nil {
		for _, module := range provider.modules {
			module.AddFlags(cmd)
		}
	}
}


func (provider *BaseProvider) NewDaemon(cmd *cobra.Command, tmCfg *cfg.Config, logger log.Logger, tmNode *node.Node, daemonApp *tm.DaemonApp, config common.DaemonConfig) *Daemon {
	dm := NewDaemon(tmCfg, logger, tmNode, config, daemonApp)
	if provider.modules != nil {
		for _, module := range provider.modules {
			for _, fac := range module.Factories() {
				dm.RegisterWorkerFactory(fac)
			}
		}
		
		dm.BeforeStartingHandler = func(dm *Daemon) {
			for _, module := range provider.modules {
				moduleConfig := module.GetDefaultConfig()
				module.BeforeDaemonStarting(cmd, dm, moduleConfig)
			}
		}
		
		dm.AfterStartedHandler = func(dm *Daemon) {
			for _, module := range provider.modules {
				module.AfterDaemonStarted(dm)
			}
		}
	}
	
	
	return dm
}
