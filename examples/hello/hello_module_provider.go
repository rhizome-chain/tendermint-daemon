package hello

import (
	"github.com/spf13/cobra"
	
	"github.com/tendermint/tendermint/config"
	
	"github.com/rhizome-chain/tendermint-daemon/daemon"
	"github.com/rhizome-chain/tendermint-daemon/daemon/common"
	"github.com/rhizome-chain/tendermint-daemon/types"
)

type ModuleProvider struct {
}

var _ daemon.ModuleProvider = (*ModuleProvider)(nil)

func (b *ModuleProvider) NewModule(tmCfg *config.Config, config common.DaemonConfig) daemon.Module {
	return &Module{}
}

func (b *ModuleProvider) GetDefaultConfig() types.ModuleConfig {
	return &types.EmptyModuleConfig{}
}

func (b *ModuleProvider) AddFlags(cmd *cobra.Command) {
	// Do Nothing
}

func (b *ModuleProvider) InitFile(config *config.Config) {
	// Do Nothing
}
