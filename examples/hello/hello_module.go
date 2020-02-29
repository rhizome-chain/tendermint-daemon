package hello

import (
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/config"
	
	"github.com/rhizome-chain/tendermint-daemon/daemon"
	"github.com/rhizome-chain/tendermint-daemon/daemon/worker"
	"github.com/rhizome-chain/tendermint-daemon/types"
)


const (
	SpaceHello = "hello"
)


type Module struct {
	factory *Factory
	proxyFactory *ProxyFactory
}

func (e *Module) GetFactory(name string) worker.Factory {
	if name == FactoryName {
		return e.factory
	}
	return nil
}

func (e *Module) Name() string {
	return "hello"
}

func (e *Module) GetConfig() types.ModuleConfig {
	config := &types.EmptyModuleConfig{}
	return config
}

func (e *Module) Init(config *config.Config) {
	e.factory = &Factory{}
	e.proxyFactory = &ProxyFactory{}
}

func (e *Module) BeforeDaemonStarting(cmd *cobra.Command, dm *daemon.Daemon) {
	// DO NOTHING
}

func (e *Module) AfterDaemonStarted(dm *daemon.Daemon) {
	// DO Nothing
}

func (e *Module) Factories() (facs []worker.Factory) {
	return []worker.Factory{e.factory,e.proxyFactory}
}

var _ daemon.Module = (*Module)(nil)
