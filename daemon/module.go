package daemon

import (
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/config"
	
	"github.com/rhizome-chain/tendermint-daemon/daemon/worker"
	"github.com/rhizome-chain/tendermint-daemon/types"
)

// 존재 이유가 약함
type Module interface {
	Name() string
	GetFactory(name string) worker.Factory
	Factories() (facs []worker.Factory)
	Init(config *config.Config)
	GetConfig() types.ModuleConfig
	BeforeDaemonStarting(cmd *cobra.Command, dm *Daemon)
	AfterDaemonStarted(dm *Daemon)
}

type BaseModule struct {
}

func (b BaseModule) GetFactory(name string) worker.Factory {
	return nil
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
