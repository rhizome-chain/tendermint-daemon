package cluster

import (
	"github.com/rhizome-chain/tendermint-daemon/daemon/common"
	"github.com/rhizome-chain/tendermint-daemon/types"
)

const (
	MemberChangedEventPath = types.EventPath("member-changed")
	LeaderChangedEventPath = types.EventPath("leader-changed")
)

type MemberChangedEvent struct {
	common.DaemonEvent
	IsLeader       bool
	AliveMemberIDs []string
	AliveMembers   []*Member
}

func (event MemberChangedEvent) Path() types.EventPath { return MemberChangedEventPath }

type LeaderChangedEvent struct {
	common.DaemonEvent
	IsLeader bool
	Leader   *Member
}

func (event LeaderChangedEvent) Path() types.EventPath { return LeaderChangedEventPath }
