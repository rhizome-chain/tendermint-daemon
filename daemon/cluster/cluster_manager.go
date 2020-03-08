package cluster

import (
	"github.com/rhizome-chain/tendermint-daemon/daemon/common"
)

type BaseManager struct {
	common.Context
	cluster *Cluster
	dao     Repository
	running bool
}

var _ Manager = (*BaseManager)(nil)

func NewBaseManager(context common.Context) Manager {
	cluster := newCluster(context.GetConfig().ChainID)
	dao := NewRepository(context.GetConfig(), context, context.GetClient())
	
	manager := &BaseManager{
		Context: context,
		// config:  config,
		// logger:  logger,
		cluster: cluster,
		dao:     dao,
	}
	
	localMemb := Member{
		NodeID:    context.GetConfig().NodeID,
		Name:      context.GetConfig().NodeName,
		RPCAddr:   context.GetTMConfig().RPC.ListenAddress,
		APIAddr:   context.GetConfig().APIAddr,
		heartbeat: context.LastBlockIndex(),
		leader:    false,
		alive:     true,
		local:     true,
	}
	
	cluster.putMember(&localMemb)
	cluster.localMember = &localMemb
	
	return manager
}

// IsLeaderNode : returns whether this kernel is leader.
func (manager *BaseManager) IsLeaderNode() bool {
	return manager.cluster.localMember.IsLeader()
}

// GetCluster get cluster
func (manager *BaseManager) GetCluster() *Cluster {
	return manager.cluster
}

func (manager *BaseManager) Start() {
	panic("BaseManager must be implemented.")
}
