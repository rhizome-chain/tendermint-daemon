package daemon

import (
	"errors"
	"fmt"
	"time"
	
	"github.com/rhizome-chain/tendermint-daemon/daemon/cluster"
	"github.com/rhizome-chain/tendermint-daemon/daemon/common"
	"github.com/rhizome-chain/tendermint-daemon/daemon/job"
	"github.com/rhizome-chain/tendermint-daemon/daemon/worker"
	
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	
	"github.com/rhizome-chain/tendermint-daemon/types"
)

type Daemon struct {
	id                    string
	tmNode                *node.Node
	tmCfg                 *cfg.Config
	logger                log.Logger
	client                types.Client
	context               common.Context
	config                common.DaemonConfig
	modules               map[string]Module
	spaceRegistry         types.SpaceRegistry
	clusterManager        cluster.Manager
	jobManager            *job.Manager
	workerManager         *worker.Manager
	jobOrganizer          job.Organizer
	beforeStartingHandler func(dm *Daemon)
}

var _ worker.ProxyProvider = (*Daemon)(nil)

func NewDaemon(tmCfg *cfg.Config, logger log.Logger, tmNode *node.Node, config common.DaemonConfig,
	spaceRegistry types.SpaceRegistry) (dm *Daemon) {
	ctx := common.NewContext(tmCfg, logger, tmNode, config)
	dm = &Daemon{
		context: ctx,
		config:  config,
		tmCfg:   tmCfg,
		logger:  logger,
		tmNode:  tmNode,
		modules: make(map[string]Module),
	}
	
	dm.client = ctx.GetClient()
	dm.id = string(dm.tmNode.NodeInfo().ID())
	dm.spaceRegistry = spaceRegistry
	spaceRegistry.RegisterSpace(common.SpaceDaemon)
	
	dm.clusterManager = cluster.NewManager2(ctx)
	ctx.SetClusterState(dm.clusterManager.GetCluster())
	dm.jobManager = job.NewManager(ctx)
	dm.workerManager = worker.NewManager(ctx, spaceRegistry, dm)
	
	return dm
}

func (dm *Daemon) GetModule(name string) Module {
	module, ok := dm.modules[name]
	if !ok {
		panic(fmt.Sprintf("Module named '%s' is not registered.", name))
	}
	return module
}

func (dm *Daemon) AddModule(module Module) {
	_, ok := dm.modules[module.Name()]
	if ok {
		panic(fmt.Sprintf("Module named '%s' is already registered.", module.Name()))
	}
	dm.modules[module.Name()] = module
}

func (dm *Daemon) GetTMConfig() cfg.Config {
	return *dm.tmCfg
}

func (dm *Daemon) GetDaemonConfig() common.DaemonConfig {
	return dm.config
}

func (dm *Daemon) GetContext() common.Context {
	return dm.context
}

func (dm *Daemon) SetJobOrganizer(organizer job.Organizer) {
	dm.jobOrganizer = organizer
}

// RegisterWorkerFactory register worker.Factory
func (dm *Daemon) RegisterWorkerFactory(factory worker.Factory) {
	dm.workerManager.RegisterWorkerFactory(factory)
}

func (dm *Daemon) Start() {
	go func() {
		dm.waitReady()
		
		dm.logger.Info("[Daemon] Starting Daemon...", "node_id", dm.tmNode.NodeInfo().ID())
		
		if dm.GetContext().IsValidator() {
			dm.logger.Info("[Daemon] I'm a Validator")
		}
		
		if dm.beforeStartingHandler != nil {
			dm.beforeStartingHandler(dm)
		}
		
		dm.clusterManager.Start()
		dm.jobManager.Start()
		dm.workerManager.Start()
		
		if dm.jobOrganizer == nil {
			dm.jobOrganizer = job.NewSimpleOrganizer(dm.logger)
		}
		
		// common.StartDaemonEventBus()
		
		common.SubscribeDaemonEvent(cluster.MemberChangedEventPath,
			"daemon-onMemberChanged",
			dm.onMemberChanged)
		
		common.SubscribeDaemonEvent(job.MemberJobsChangedEventPath,
			"daemon-onMemberJobsChanged",
			dm.onMemberJobsChanged)
		
		common.SubscribeDaemonEvent(job.JobsChangedEventPath,
			"daemon-onJobsChanged",
			dm.onJobsChanged)
		
		for _, module := range dm.modules {
			module.AfterDaemonStarted(dm)
		}
		
		dm.checkJobsAndAllocate()
	}()
}

func (dm *Daemon) waitReady() {
	threshold := time.Second * 3
	for time.Now().Sub(dm.tmNode.ConsensusState().GetState().LastBlockTime) > threshold {
		time.Sleep(200 * time.Millisecond)
	}
}

func (dm *Daemon) ID() string { return dm.id }

func (dm *Daemon) GetClient() types.Client { return dm.client }

func (dm *Daemon) GetCluster() *cluster.Cluster { return dm.clusterManager.GetCluster() }

func (dm *Daemon) IsLeaderNode() bool { return dm.clusterManager.IsLeaderNode() }

func (dm *Daemon) GetJobRepository() job.Repository { return dm.jobManager.GetRepository() }

func (dm *Daemon) NewWorkerProxy(jobID string) (worker.Proxy, error) {
	job, err := dm.GetJobRepository().GetJob(jobID)
	
	if err != nil {
		if types.IsNoDataError(err){
			err = errors.New("Cannot find Job " + jobID)
		}
		return nil, err
	}
	
	return dm.workerManager.NewWorkerProxy(job)
}

// member's jobs changed event handler
func (dm *Daemon) onMemberChanged(event types.Event) {
	dm.logger.Debug(" - [Daemon] onMemberChanged :", event)
	memberChangedEvent := event.(cluster.MemberChangedEvent)
	
	if memberChangedEvent.IsLeader {
		dm.logger.Info("[Daemon] members changed", "members", memberChangedEvent.AliveMemberIDs)
		dm.distributeJobs(memberChangedEvent.AliveMemberIDs)
	}
}

// member's jobs changed event handler
func (dm *Daemon) onMemberJobsChanged(event types.Event) {
	dm.logger.Debug("-[Daemon] onMemberJobsChanged :", event)
	
	memberJobsChangedEvent := event.(job.MemberJobsChangedEvent)
	
	if memberJobsChangedEvent.NodeID != dm.ID() {
		return
	}
	
	dm.logger.Info("[Daemon] member's jobs changed", "nodeID", memberJobsChangedEvent.NodeID,
		"jobs", memberJobsChangedEvent.JobIDs)
	
	dm.checkJobsAndAllocate()
}

func (dm *Daemon) checkJobsAndAllocate() {
	jobs, err := dm.jobManager.GetRepository().GetMemberJobs(dm.ID())
	
	if err != nil && !types.IsNoDataError(err) {
		dm.logger.Error(fmt.Sprintf("[Daemon] cannot get %s's jobs", dm.ID()))
		return
	}
	
	dm.workerManager.SetJobs(jobs)
}

func (dm *Daemon) onJobsChanged(event types.Event) {
	jobsChangedEvent := event.(job.JobsChangedEvent)
	if dm.jobOrganizer == nil {
		dm.logger.Info("[Daemon] JobOrganizer is not set.")
	}
	if dm.IsLeaderNode() {
		dm.logger.Info(" - [Daemon] onJobsChanged :", "blockHeight", jobsChangedEvent.BlockHeight)
		aliveMembers := dm.clusterManager.GetCluster().GetAliveMemberIDs()
		dm.distributeJobs(aliveMembers)
	}
}

func (dm *Daemon) distributeJobs(aliveMembers []string) {
	allJobs, err := dm.jobManager.GetRepository().GetAllJobs()
	if err != nil && !types.IsNoDataError(err) {
		dm.logger.Error("[Daemon] distributeJobs - GetAllJobs ", err)
		return
	}
	membJobMap, err := dm.jobManager.GetRepository().GetAllMemberJobIDs()
	if err != nil && !types.IsNoDataError(err) {
		dm.logger.Error("[Daemon] distributeJobs - GetAllMemberJobIDs ", err)
		return
	}
	
	newMembJobs, err := dm.jobOrganizer.Distribute(allJobs, aliveMembers, membJobMap)
	
	if err != nil {
		dm.logger.Error("[Daemon] distributeJobs : ", "new members' jobs", err)
		return
	}
	
	dm.logger.Info("[Daemon] distributeJobs : ", "new members' jobs", newMembJobs)
	
	for nodeid, jobs := range newMembJobs {
		err = dm.jobManager.GetRepository().PutMemberJobIDs(nodeid, jobs)
		if err != nil {
			dm.logger.Error(fmt.Sprintf("[Daemon] Put Member(%s) jobs. %s", nodeid, jobs), err)
		} else {
			dm.logger.Info(fmt.Sprintf("[Daemon] Put Member(%s) jobs. %s", nodeid, jobs))
		}
	}
	
}
