package cluster

import (
	"fmt"
	
	"github.com/rhizome-chain/tendermint-daemon/daemon/common"
	tmevents "github.com/rhizome-chain/tendermint-daemon/tm/events"
	"github.com/rhizome-chain/tendermint-daemon/types"
)

// Manager1 is a cluster manager which checks member's health by Heartbeat TX
type Manager1 struct {
	*BaseManager
}

var _ Manager = (*Manager1)(nil)

func NewManager1(context common.Context) Manager {
	manager := &Manager1{
		BaseManager: NewBaseManager(context).(*BaseManager),
	}
	
	return manager
}

func (manager *Manager1) Start() {
	manager.running = true
	
	err := manager.dao.PutMember(manager.cluster.localMember)
	
	if err != nil {
		manager.Error("Cannot PutMember.", err)
		panic(err)
	}
	
	err = manager.dao.PutHeartbeat(manager.GetNodeID(), manager.Context.LastBlockIndex())
	
	if err != nil {
		manager.Error("Cannot send Heartbeat.", err)
		panic(err)
	}
	
	tmevents.SubscribeBlockEvent(tmevents.BeginBlockEventPath, "cluster-heartbeat", func(event types.Event) {
		// fmt.Println("heartbeat ", event)
		blockEvent := event.(tmevents.BeginBlockEvent)
		
		err := manager.dao.PutHeartbeat(manager.GetNodeID(),blockEvent.Height)
		if err != nil {
			manager.Error("Cannot send Heartbeat.", err)
		}
	})
	
	tmevents.SubscribeBlockEvent(tmevents.EndBlockEventPath, "cluster-checkMembers", func(event types.Event) {
		changed := false
		blockEvent := event.(tmevents.EndBlockEvent)
		err := manager.dao.GetHeartbeats(func(nodeid string, blockHeight int64) {
			c := manager.handleHeartbeat(nodeid, blockHeight, blockEvent.Height)
			if c {
				changed = true
			}
		})
		
		if err != nil {
			manager.Error("[FATAL] Cannot check heartbeats.", err)
		}
		
		manager.checkLeader(changed)
		
		if changed {
			manager.onMemberChanged()
		}
		
	})
	
	manager.Info("[INFO-Cluster] Start Cluster Manager.")
}

// returns true if member state changed
func (manager *Manager1) handleHeartbeat(nodeid string, heartbeat int64, blockHeight int64) (changed bool) {
	member := manager.cluster.GetMember(nodeid)
	if member == nil {
		memb, err := manager.dao.GetMember(nodeid)
		if err != nil {
			manager.Error(fmt.Sprintf("[FATAL] Cannot find member [%s]", nodeid), err)
			return false
		}
		member = memb
		manager.cluster.putMember(member)
	}
	
	oldAlive := member.IsAlive()
	
	if member.IsLocal() {
		member.SetAlive(true)
	} else  {
		gap := blockHeight - heartbeat
		if gap > int64(manager.GetConfig().AliveThresholdBlocks) {
			manager.Info(fmt.Sprintf("Member[%s:%s] haven't sent heartbeat for %d blocks.", member.Name, member.NodeID, gap))
			member.SetAlive(false)
		} else {
			member.SetAlive(true)
		}
	}
	
	member.SetHeartbeat(heartbeat)
	
	changed = oldAlive != member.IsAlive()
	
	return changed
}

func (manager *Manager1) checkLeader(memberChanged bool) {
	if manager.cluster.leader != nil && manager.cluster.leader.IsLocal() {
		return
	}
	
	oldLeader := manager.cluster.leader
	
	if oldLeader != nil && !memberChanged {
		return
	}
	
	leaderID, err := manager.dao.GetLeader()
	if err != nil {
		manager.Error("Get Leader ", err)
	}
	
	manager.Info("[INFO-Cluster] Leader is", "leaderID", leaderID)
	
	if oldLeader != nil {
		if oldLeader.NodeID == leaderID {
			if oldLeader.IsAlive() {
				return
			}
			manager.Info("[INFO-Cluster] Old leader is dead. ", "old_node_ID", oldLeader.NodeID)
			oldLeader.SetLeader(false)
			manager.cluster.leader = nil
		} else {
			oldLeader.SetLeader(false)
			manager.cluster.leader = nil
		}
	}
	
	var leader *Member
	
	if len(leaderID) > 0 {
		leader = manager.cluster.GetMember(leaderID)
		if leader == nil {
			manager.Info("[INFO-Cluster] Old leader is missing. ", "leaderID", leaderID)
			leader = manager.electLeader()
		} else if !leader.IsAlive() {
			manager.Info("[INFO-Cluster] Old leader is dead. ", "leaderID", leaderID)
			leader.SetLeader(false)
			leader = manager.electLeader()
		}
	} else {
		leader = manager.electLeader()
	}
	
	leader.SetLeader(true)
	
	manager.cluster.leader = leader
	
	manager.onLeaderChanged(leader)
}

func (manager *Manager1) electLeader() *Member {
	members := manager.cluster.GetSortedMembers()
	
	// fmt.Println("****** electLeader:: len(members) ", len(members))
	
	for _, id := range members {
		memb := manager.cluster.GetMember(id)
		// fmt.Println("    ****** electLeader:: member ", id, memb)
		if memb.IsAlive() {
			if memb.IsLocal() {
				manager.dao.PutLeader(id)
			}
			return memb
		}
	}
	//
	// local := manager.cluster.Local()
	// manager.dao.PutLeader(local.NodeID)
	manager.Error("No Leader elected.")
	return nil
}

func (manager *Manager1) onLeaderChanged(leader *Member) {
	if manager.cluster.localMember == leader {
		manager.Info("[INFO-Cluster] Leader changed. I'm the leader", "leader",leader.NodeID )
		manager.cluster.localMember.SetLeader(true)
	} else {
		manager.Info("[INFO-Cluster] Leader is set", "leader", leader.NodeID)
	}
	
	common.PublishDaemonEvent(LeaderChangedEvent{
		IsLeader: manager.cluster.localMember.IsLeader(),
		Leader:   leader,
	})
}

func (manager *Manager1) onMemberChanged() {
	manager.Info("[INFO-Cluster] Members changed.", "members",
		manager.cluster.GetAliveMemberIDs())
	
	common.PublishDaemonEvent(MemberChangedEvent{
		IsLeader:       manager.cluster.localMember.IsLeader(),
		AliveMemberIDs: manager.cluster.GetAliveMemberIDs(),
		AliveMembers:   manager.cluster.GetAliveMembers(),
	})
}
