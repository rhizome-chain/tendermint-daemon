package hello

import (
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/config"
	
	"github.com/rhizome-chain/tendermint-daemon/daemon"
	"github.com/rhizome-chain/tendermint-daemon/daemon/worker"
	"github.com/rhizome-chain/tendermint-daemon/types"
)

type Module struct {
}

func (e *Module) Name() string {
	return "hello"
}

func (e *Module) GetConfig() types.ModuleConfig {
	config := &types.EmptyModuleConfig{}
	return config
}

func (e *Module) Init(config *config.Config) {
	// DO NOTHING
}

func (e *Module) BeforeDaemonStarting(cmd *cobra.Command, dm *daemon.Daemon) {
	// DO NOTHING
}

func (e *Module) AfterDaemonStarted(dm *daemon.Daemon) {
	// DO Nothing
}

func (e *Module) Factories() (facs []worker.Factory) {
	return []worker.Factory{&Factory{}}
}

var _ daemon.Module = (*Module)(nil)
