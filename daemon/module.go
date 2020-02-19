package daemon

import (
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/config"
	
	"github.com/rhizome-chain/tendermint-daemon/daemon/worker"
	"github.com/rhizome-chain/tendermint-daemon/types"
)

type Module interface {
	GetDefaultConfig() types.ModuleConfig
	Factories() (facs []worker.Factory)
	AddFlags(cmd *cobra.Command)
	BeforeDaemonStarting(cmd *cobra.Command, dm *Daemon, moduleConfig types.ModuleConfig)
	AfterDaemonStarted(dm *Daemon)
	InitFile(config *config.Config)
}

type BaseModule struct {

}

func (b BaseModule) GetDefaultConfig() types.ModuleConfig {
	return &types.EmptyModuleConfig{}
}

func (b BaseModule) Factories() (facs []worker.Factory) {
	return []worker.Factory{}
}

func (b BaseModule) AddFlags(cmd *cobra.Command) {
	// To be implemented
}

func (b BaseModule) BeforeDaemonStarting(cmd *cobra.Command, dm *Daemon, moduleConfig types.ModuleConfig) {
	// To be implemented
}

func (b BaseModule) AfterDaemonStarted(dm *Daemon) {
	// To be implemented
}

func (b BaseModule) InitFile(config *config.Config) {
	// To be implemented
}

var _ Module = (*BaseModule)(nil)
