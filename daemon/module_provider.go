package daemon

import (
	"github.com/spf13/cobra"
	
	"github.com/tendermint/tendermint/config"
	
	"github.com/rhizome-chain/tendermint-daemon/daemon/common"
	"github.com/rhizome-chain/tendermint-daemon/types"
)

type ModuleProvider interface {
	NewModule(tmCfg *config.Config, config common.DaemonConfig) Module
	GetDefaultConfig() types.ModuleConfig
	AddFlags(cmd *cobra.Command)
	InitFile(config *config.Config)
}

type BaseModuleProvider struct {
}

func (b *BaseModuleProvider) NewModule(tmCfg *config.Config, config common.DaemonConfig) Module {
	// To be implemented
	return nil
}

func (b *BaseModuleProvider) GetDefaultConfig() types.ModuleConfig {
	return &types.EmptyModuleConfig{}
}

func (b *BaseModuleProvider) AddFlags(cmd *cobra.Command) {
	// To be implemented
}

func (b *BaseModuleProvider) InitFile(config *config.Config) {
	// To be implemented
}

var _ ModuleProvider = (*BaseModuleProvider)(nil)
