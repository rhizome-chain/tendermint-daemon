package cluster

import (
	"fmt"
	"github.com/rhizome-chain/tendermint-daemon/daemon/common"
	tmevents "github.com/rhizome-chain/tendermint-daemon/tm/events"
	"github.com/rhizome-chain/tendermint-daemon/types"
	"time"
)

type Manager struct {
	common.Context
	cluster *Cluster
	dao     Repository
	running bool
}

func NewManager(context common.Context) *Manager {
	//config common.DaemonConfig, logger log.Logger, client types.Client
	cluster := newCluster(context.GetConfig().ChainID)
	dao := NewRepository(context.GetConfig(), context, context.GetClient())
	
	manager := &Manager{
		Context: context,
		//config:  config,
		//logger:  logger,
		cluster: cluster,
		dao:     dao,
	}
	
	localMemb := Member{
		NodeID:    context.GetConfig().NodeID,
		Name:      context.GetConfig().NodeName,
		heartbeat: time.Now(),
		leader:    false,
		alive:     true,
		local:     true,
	}
	
	cluster.putMember(&localMemb)
	cluster.localMember = &localMemb
	
	return manager
}

func (manager *Manager) Start() {
	manager.running = true
	
	err := manager.dao.PutMember(manager.cluster.localMember)
	
	if err != nil {
		manager.Error("Cannot PutMember.", err)
		panic(err)
	}
	
	err = manager.dao.PutHeartbeat(manager.GetNodeID())
	
	if err != nil {
		manager.Error("Cannot send Heartbeat.", err)
		panic(err)
	}
	
	tmevents.SubscribeBlockEvent(tmevents.BeginBlockEventPath, "cluster-heartbeat", func(event types.Event) {
		// fmt.Println("heartbeat ", event)
		err := manager.dao.PutHeartbeat(manager.GetNodeID())
		if err != nil {
			manager.Error("Cannot send Heartbeat.", err)
		}
	})
	
	tmevents.SubscribeBlockEvent(tmevents.EndBlockEventPath, "cluster-checkMembers", func(event types.Event) {
		changed := false
		err := manager.dao.GetHeartbeats(func(nodeid string, tm time.Time) {
			c := manager.handleHeartbeat(nodeid, tm)
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
func (manager *Manager) handleHeartbeat(nodeid string, tm time.Time) (changed bool) {
	member := manager.cluster.GetMember(nodeid)
	if member == nil {
		memb, err := manager.dao.GetMember(nodeid)
		if err != nil {
			manager.Error(fmt.Sprintf("[FATAL] Cannot find member [%s]", nodeid), err)
			return false
		}
		member = &memb
		manager.cluster.putMember(member)
	}
	
	oldAlive := member.IsAlive()
	
	if member.IsLocal() {
		member.SetAlive(true)
	} else if member.Heartbeat().IsZero() || tm.Equal(member.Heartbeat()) {
		gap := time.Now().Sub(tm).Seconds()
		if gap > float64(manager.GetConfig().AliveThresholdSeconds) {
			manager.Info(fmt.Sprintf("Member[%s] haven't sent heartbeat for %f seconds.", member.NodeID, gap))
			member.SetAlive(false)
		} else {
			member.SetAlive(true)
		}
	} else {
		member.SetAlive(true)
	}
	
	member.SetHeartbeat(tm)
	
	changed = oldAlive != member.IsAlive()
	
	return changed
}

func (manager *Manager) checkLeader(memberChanged bool) {
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

func (manager *Manager) electLeader() *Member {
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

// IsLeaderNode : returns whether this kernel is leader.
func (manager *Manager) IsLeaderNode() bool {
	return manager.cluster.localMember.IsLeader()
}

func (manager *Manager) onLeaderChanged(leader *Member) {
	if manager.cluster.localMember == leader {
		manager.Info("[INFO-Cluster] Leader changed. I'm the leader")
		manager.cluster.localMember.SetLeader(true)
	} else {
		manager.Info("[INFO-Cluster] Leader is set", "leader",
			manager.cluster.leader.NodeID)
	}
	
	common.PublishDaemonEvent(LeaderChangedEvent{
		IsLeader: manager.cluster.localMember.IsLeader(),
		Leader:   leader,
	})
}

func (manager *Manager) onMemberChanged() {
	manager.Info("[INFO-Cluster] Members changed.", "members",
		manager.cluster.GetAliveMemberIDs())
	
	common.PublishDaemonEvent(MemberChangedEvent{
		IsLeader:       manager.cluster.localMember.IsLeader(),
		AliveMemberIDs: manager.cluster.GetAliveMemberIDs(),
		AliveMembers:   manager.cluster.GetAliveMembers(),
	})
}

// GetCluster get cluster
func (manager *Manager) GetCluster() *Cluster {
	return manager.cluster
}
