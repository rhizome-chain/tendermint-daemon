package daemon

import (
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/config"
	
	"github.com/rhizome-chain/tendermint-daemon/daemon/worker"
	"github.com/rhizome-chain/tendermint-daemon/types"
)

type Module interface {
	Name() string
	Factories() (facs []worker.Factory)
	Init(config *config.Config)
	GetConfig() types.ModuleConfig
	BeforeDaemonStarting(cmd *cobra.Command, dm *Daemon)
	AfterDaemonStarted(dm *Daemon)
}

type BaseModule struct {
}

func (b BaseModule) Name() string {
	return "base-module"
}

func (b BaseModule) Init(config *config.Config) {
	// To be implemented
}

func (b BaseModule) GetConfig() types.ModuleConfig {
	return &types.EmptyModuleConfig{}
}

func (b BaseModule) Factories() (facs []worker.Factory) {
	return []worker.Factory{}
}

func (b BaseModule) BeforeDaemonStarting(cmd *cobra.Command, dm *Daemon) {
	// To be implemented
}

func (b BaseModule) AfterDaemonStarted(dm *Daemon) {
	// To be implemented
}

var _ Module = (*BaseModule)(nil)
