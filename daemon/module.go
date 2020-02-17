package daemon

import (
	"github.com/rhizome-chain/tendermint-daemon/daemon/worker"
	"github.com/rhizome-chain/tendermint-daemon/types"
	"github.com/spf13/cobra"
)

type Module interface {
	GetDefaultConfig() types.ModuleConfig
	Factories() (facs []worker.Factory)
	AddFlags(cmd *cobra.Command)
	BeforeDaemonStarting(cmd *cobra.Command, dm *Daemon, moduleConfig types.ModuleConfig)
	AfterDaemonStarted(dm *Daemon)
}
