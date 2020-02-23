package cluster

import (
	"fmt"
	"sort"
	"strings"
	"time"
	
	"github.com/rhizome-chain/tendermint-daemon/daemon/common"
	
	"github.com/rhizome-chain/tendermint-daemon/types"
)

func init() {
	types.BasicCdc.RegisterConcrete(Member{}, "daemon/member", nil)
}

type Repository interface {
	PutMember(member *Member) (err error)
	GetMember(nodeID string) (member Member, err error)
	HasMember(nodeID string) (ok bool)
	PutLeader(leader string) (err error)
	GetLeader() (leader string, err error)
	GetAllMembers() (members []*Member, err error)
	PutHeartbeat(nodeID string) (err error)
	GetHeartbeats(handler func(nodeid string, time time.Time)) (err error)
}

// Cluster ..
type Cluster struct {
	name        string
	membIDs     []string
	members     map[string]*Member
	localMember *Member
	leader      *Member
}

var _ common.ClusterState = (*Cluster)(nil)

func newCluster(name string) *Cluster {
	cluster := Cluster{name: name}
	cluster.membIDs = []string{}
	cluster.members = make(map[string]*Member)
	
	return &cluster
}

// Name get name
func (cluster *Cluster) Name() string {
	return cluster.name
}

func (cluster *Cluster) putMember(memb *Member) {
	if _, ok := cluster.members[memb.NodeID]; !ok {
		cluster.membIDs = append(cluster.membIDs, memb.NodeID)
		sort.Strings(cluster.membIDs)
	}
	// fmt.Println("********* putMember :: ", cluster.members, memb.ID, memb)
	cluster.members[memb.NodeID] = memb
}

func (cluster *Cluster) removeMember(id string) {
	index := -1
	for i, mid := range cluster.membIDs {
		if mid == id {
			index = i
			break
		}
	}
	
	if index > -1 {
		if index <= len(cluster.membIDs)-1 {
			copy(cluster.membIDs[index:], cluster.membIDs[index+1:])
		}
		cluster.membIDs = cluster.membIDs[0 : len(cluster.membIDs)-1]
	}
	
	delete(cluster.members, id)
}

// GetMember get member with given name
func (cluster *Cluster) GetMember(id string) *Member {
	memb := cluster.members[id]
	return memb
}

// GetSortedMembers get all member ids
func (cluster *Cluster) GetSortedMembers() []string {
	return cluster.membIDs
}

// GetAliveMembers get active members
func (cluster *Cluster) GetAliveMembers() []*Member {
	membs := []*Member{}
	for _, memb := range cluster.members {
		if memb.IsAlive() {
			membs = append(membs, memb)
		}
	}
	return membs
}

// GetAliveMemberIDs get active member IDs
func (cluster *Cluster) GetAliveMemberIDs() []string {
	membs := []string{}
	for id, memb := range cluster.members {
		if memb.IsAlive() {
			membs = append(membs, id)
		}
	}
	return membs
}

// Leader get Leader
func (cluster *Cluster) Leader() *Member {
	return cluster.leader
}

// Local get localMember
func (cluster *Cluster) Local() *Member {
	return cluster.localMember
}

// Local get localMember
func (cluster *Cluster) IsLeader() bool {
	return cluster.localMember.IsLeader()
}

// Member member info
type Member struct {
	NodeID    string    `json:"nodeid"`
	Name      string    `json:"name"`
	heartbeat time.Time `json:"heartbeat"`
	leader    bool      // transient field
	alive     bool      // transient field
	local     bool      // transient field
}

var (
	NilMember = Member{}
)

func NewMember(name string, nodeid string) *Member {
	return &Member{NodeID: nodeid, Name: name}
}

// IsLeader return whether member is leader
func (m *Member) IsLeader() bool {
	return m.leader
}

// SetLeader Set member as leader
func (m *Member) SetLeader(leader bool) {
	m.leader = leader
}

// IsAlive return whether member is alive
func (m *Member) IsAlive() bool {
	return m.alive
}

// SetAlive Set member alive
func (m *Member) SetAlive(alive bool) {
	m.alive = alive
}

// IsLocal return whether member is alive
func (m *Member) IsLocal() bool {
	return m.local
}

// SetLocal Set member alive
func (m *Member) SetLocal(local bool) {
	m.local = local
}

// SetLocal Set member alive
func (m *Member) SetHeartbeat(time time.Time) {
	m.heartbeat = time
}

// SetLocal Set member alive
func (m *Member) Heartbeat() time.Time {
	return m.heartbeat
}

// String implement fmt.Stringer
func (m *Member) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Member[%s:%s] %s alive=%t, leader=%t`,
		m.Name, m.NodeID, m.heartbeat, m.alive, m.leader))
}
