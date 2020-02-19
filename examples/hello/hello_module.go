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

var _ daemon.Module = (*Module)(nil)


func (e Module) GetDefaultConfig() types.ModuleConfig {
	config := &types.EmptyModuleConfig{}
	return config
}

func (e Module) Factories() (facs []worker.Factory) {
	return []worker.Factory{&Factory{}}
}

func (e Module) AddFlags(cmd *cobra.Command) {
	// DO NOTHING
}


func (e Module) BeforeDaemonStarting(cmd *cobra.Command, dm *daemon.Daemon, moduleConfig types.ModuleConfig) {
	// DO Nothing
}

func (e Module) AfterDaemonStarted(dm *daemon.Daemon) {
	// DO Nothing
}


func (e Module) InitFile(config *config.Config) {
	// DO Nothing
}
