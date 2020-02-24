package daemon

import (
	"fmt"
	"path/filepath"
	
	"github.com/spf13/cobra"
	
	"github.com/rhizome-chain/tendermint-daemon/daemon/common"
	"github.com/rhizome-chain/tendermint-daemon/tm"
	
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
)

type Provider func(tmCfg *cfg.Config, logger log.Logger, tmNode *node.Node, daemonApp *tm.DaemonApp, config common.DaemonConfig) *Daemon

type BaseProvider struct {
	moduleProviders []ModuleProvider
}

func (provider *BaseProvider) AddModuleProvider(moduleProvider ModuleProvider) {
	if provider.moduleProviders == nil {
		provider.moduleProviders = []ModuleProvider{}
	}
	provider.moduleProviders = append(provider.moduleProviders, moduleProvider)
}

func (provider *BaseProvider) AddFlags(cmd *cobra.Command) {
	common.AddDaemonFlags(cmd)
	
	if provider.moduleProviders != nil {
		for _, provider := range provider.moduleProviders {
			provider.AddFlags(cmd)
		}
	}
}

func (provider *BaseProvider) InitFiles(config *cfg.Config, daemonConfig *common.DaemonConfig) {
	confFilePath := filepath.Join(config.RootDir, "config", "daemon.toml")
	common.WriteConfigFile(confFilePath, daemonConfig)
	if provider.moduleProviders != nil {
		for _, provider := range provider.moduleProviders {
			provider.InitFile(config)
		}
	}
}

func (provider *BaseProvider) NewDaemon(cmd *cobra.Command, tmCfg *cfg.Config, logger log.Logger,
	tmNode *node.Node, daemonApp *tm.DaemonApp, config common.DaemonConfig) *Daemon {
	
	dm := NewDaemon(tmCfg, logger, tmNode, config, daemonApp)
	
	if provider.moduleProviders != nil {
		for _, moduleProvider := range provider.moduleProviders {
			module := moduleProvider.NewModule(tmCfg, config)
			
			module.Init(tmCfg)
			for _, fac := range module.Factories() {
				dm.workerManager.RegisterWorkerFactory(fac)
			}
			
			dm.AddModule(module)
			dm.logger.Info(fmt.Sprintf("Init Module[%s].", module.Name()), "config", module.GetConfig())
		}
	}
	
	if dm.modules != nil {
		dm.beforeStartingHandler = func(dm *Daemon) {
			for _, module := range dm.modules {
				module.BeforeDaemonStarting(cmd, dm)
			}
		}
	}
	
	return dm
}
