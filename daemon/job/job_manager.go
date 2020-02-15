package job

import (
	"github.com/rhizome-chain/tendermint-daemon/daemon/common"
	tmevents "github.com/rhizome-chain/tendermint-daemon/tm/events"
	"github.com/rhizome-chain/tendermint-daemon/types"
)

// Manager manager for jobs
type Manager struct {
	common.Context
	nodeID string
	dao    Repository
	//logger log.Logger
}

// NewManager ..
func NewManager(context common.Context) *Manager {
	dao := NewRepository(context.GetConfig(), context, context.GetClient())
	manager := Manager{Context: context, dao: dao, nodeID: context.GetNodeID()}
	return &manager
}

func (manager *Manager) GetRepository() Repository {
	return manager.dao
}

func (manager *Manager) Start() {
	
	var jobsChanged bool
	
	jobsEvtPath := tmevents.MakeTxEventPath(common.SpaceDaemon, PathJobs, "")
	tmevents.SubscribeTxEvent(jobsEvtPath, "job-detectJobChanged", func(event tmevents.TxEvent) {
		jobsChanged = true
	})
	
	
	tmevents.SubscribeBlockEvent(tmevents.EndBlockEventPath, "job-checkJobsChanged", func(event types.Event) {
		if jobsChanged {
			endBlockEvent := event.(tmevents.EndBlockEvent)
			
			manager.Info("[JobManager] Jobs changed.", "height",endBlockEvent.Height )
			
			common.PublishDaemonEvent(JobsChangedEvent{
				BlockHeight:endBlockEvent.Height,
			})
			
			jobsChanged = false
		}
	})
	
	
	memJobsEvtPath := tmevents.MakeTxEventPath(common.SpaceDaemon, PathMemberJobs, manager.nodeID)
	
	tmevents.SubscribeTxEvent(memJobsEvtPath, "job-memJobsChanged", func(event tmevents.TxEvent) {
		nodeID := string(event.Key)
		
		jobIDs, err := manager.dao.GetMemberJobIDs(nodeID)
		if err != nil  && !types.IsNoDataError(err){
			manager.Error("[JobManager GetAllJobIDs", err)
		}
		common.PublishDaemonEvent(MemberJobsChangedEvent{
			NodeID: nodeID,
			JobIDs: jobIDs,
		})
	})
}
