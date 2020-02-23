package job

import (
	"github.com/rhizome-chain/tendermint-daemon/daemon/common"
	"github.com/rhizome-chain/tendermint-daemon/types"
)

const (
	JobsChangedEventPath       = types.EventPath("jobs-changed")
	MemberJobsChangedEventPath = types.EventPath("msm_jobs-changed")
)

type MemberJobsChangedEvent struct {
	common.DaemonEvent
	NodeID string
	JobIDs []string
}

func (event MemberJobsChangedEvent) Path() types.EventPath { return MemberJobsChangedEventPath }

type JobsChangedEvent struct {
	common.DaemonEvent
	BlockHeight int64
	// JobIDs []string
}

func (event JobsChangedEvent) Path() types.EventPath { return JobsChangedEventPath }
