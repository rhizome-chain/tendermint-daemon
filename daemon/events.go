package daemon

import (
	"github.com/rhizome-chain/tendermint-daemon/daemon/common"
	"github.com/rhizome-chain/tendermint-daemon/types"
)

const (
	StartedEventPath = types.EventPath("daemon-started")
)

type StartedEvent struct {
	common.DaemonEvent
}

func (event StartedEvent) Path() types.EventPath { return StartedEventPath }

