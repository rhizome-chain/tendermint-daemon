package cluster

import (
	"github.com/rhizome-chain/tendermint-daemon/daemon/common"
	tmevents "github.com/rhizome-chain/tendermint-daemon/tm/events"
	"github.com/rhizome-chain/tendermint-daemon/types"
)

// Manager2 Single Validator Cluster Manager is a cluster manager which checks member's health by TM Peer status
type Manager2 struct {
	*BaseManager
}

var _ Manager = (*Manager2)(nil)

func NewManager2(context common.Context) Manager {
	manager := &Manager2{
		BaseManager: NewBaseManager(context).(*BaseManager),
	}
	
	return manager
}

func (manager *Manager2) Start() {
	manager.running = true
	
	err := manager.dao.PutMember(manager.cluster.localMember)
	
	if err != nil {
		manager.Error("Cannot PutMember.", err)
		panic(err)
	}
	
	tmevents.SubscribeBlockEvent(tmevents.EndBlockEventPath, "cluster-checkMembers", func(event types.Event) {
		peerIDs := manager.GetClient().GetPeerIDs()
		changed := manager.updateClusterStatus(peerIDs)
		manager.checkLeader(peerIDs)
		if manager.IsLeaderNode() && changed {
			manager.invokeMemberChanged()
		}
	})
	
	manager.Info("[INFO-Cluster] Start Cluster Manager.")
}

func (manager *Manager2) updateClusterStatus(peerIDs []string) (changed bool) {
	for _, id := range peerIDs {
		member:= manager.cluster.GetMember(id)
		if member == nil {
			member, err := manager.dao.GetMember(id)
			if err != nil {
				manager.Error("[Cluster] Get Member", err)
			} else {
				manager.cluster.putMember(member)
				changed = true
			}
		}
	}
	
	for id,member:=range manager.cluster.members{
		if member.IsLocal() {
			continue
		}
		alive := false
		for _,pid := range peerIDs {
			if id == pid {
				alive = true
				break
			}
		}
		if member.IsAlive() != alive {
			changed = true
		}
		member.SetAlive(alive)
	}
	
	return changed
}


func (manager *Manager2) checkLeader(peerIDs []string) {
	if manager.cluster.IsLeader() {
		return
	}
	
	oldLeader := manager.cluster.leader
	
	if oldLeader != nil && oldLeader.IsAlive(){
		return
	}
	
	leaderID, err := manager.dao.GetLeader()
	
	if err != nil {
		manager.Error("[Cluster] Get Leader ID ", err)
	}
	
	if len(leaderID) == 0{
		if manager.IsValidator() {
			manager.setLocalLeader()
		} else {
			manager.Error("[Cluster] Leader is not elected.")
			return
		}
	} else {
		if manager.cluster.localMember.NodeID == leaderID {
			manager.cluster.localMember.leader = true
			manager.cluster.leader = manager.cluster.localMember
			manager.Info("[Cluster] I'm the leader.")
			return
		} else {
			leader := manager.cluster.GetMember(leaderID)
			if leader == nil {
				manager.Error("[Cluster] Leader is lost", err)
				if manager.IsValidator() {
					manager.setLocalLeader()
				}
			} else {
				leader.SetLeader(true)
				manager.cluster.leader = leader
				manager.Info("[Cluster] Leader is set." , "leader",leader)
			}
		}
	}
}

func (manager *Manager2) setLocalLeader() {
	manager.dao.PutLeader(manager.cluster.Local().NodeID)
	manager.Info("[Cluster] I'm the leader.")
}


func (manager *Manager2) invokeMemberChanged() {
	manager.Info("[INFO-Cluster] Members changed.", "members",
		manager.cluster.GetAliveMemberIDs())
	
	common.PublishDaemonEvent(MemberChangedEvent{
		IsLeader:       manager.cluster.localMember.IsLeader(),
		AliveMemberIDs: manager.cluster.GetAliveMemberIDs(),
		AliveMembers:   manager.cluster.GetAliveMembers(),
	})
}

