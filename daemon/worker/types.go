package worker

import (
	"errors"
	"fmt"
	"sync"
	
	"github.com/tendermint/tendermint/libs/log"
	
	"github.com/rhizome-chain/tendermint-daemon/daemon/common"
	"github.com/rhizome-chain/tendermint-daemon/daemon/job"
)

const (
	SpaceDefaultWorker = "daemon-data"
)

type DataPath struct {
	JobID string
	Topic string
}

type DataHandler func(jobID string, topic string, rowID string, value []byte) bool

type Repository interface {
	PutCheckpoint(jobID string, checkpoint interface{}) error
	GetCheckpoint(jobID string, checkpoint interface{}) error
	PutData(space string, jobID string, topic string, rowID string, data []byte) error
	PutObject(space string, jobID string, topic string, rowID string, data interface{}) error
	PutDataSync(space string, jobID string, topic string, rowID string, data []byte) error
	PutObjectSync(space string, jobID string, topic string, rowID string, data interface{}) error
	GetData(space string, jobID string, topic string, rowID string) (data []byte, err error)
	GetObject(space string, jobID string, topic string, rowID string, data interface{}) error
	DeleteData(space string, jobID string, topic string, rowID string) error
	GetDataWithTopic(space string, jobID string, topic string, handler DataHandler) error
	
	PutDataFullPath(space string, path string, data []byte) error
	PutObjectFullPath(space string, path string, data interface{}) error
	GetDataFullPath(space string, path string) (data []byte, err error)
	GetObjectFullPath(space string, path string, data interface{}) (err error)
	DeleteDataFullPath(space string, path string) error
}

// Worker ..
type Worker interface {
	ID() string
	Start() error
	Stop() error
	IsStarted() bool
}

// Factory ..
type Factory interface {
	Name() string
	Space() string
	NewWorker(helper *Helper) (Worker, error)
}

// Helper ..
type Helper struct {
	space  string
	config common.DaemonConfig
	logger log.Logger
	job    job.Job
	dao    Repository
}

// NewHelper ..
func NewHelper(space string, config common.DaemonConfig, logger log.Logger, job job.Job, dao Repository) *Helper {
	helper := Helper{space: space, config: config, logger: logger, job: job, dao: dao}
	return &helper
}

// CreateChildHelper ...
// func (helper *Helper) CreateChildHelper(subID string, job *job.Job) *Helper {
// 	helper2 := Helper{cluster: helper.cluster, kernelid: helper.kernelid, id: helper.id + "-" + subid, job: job, kv: helper.kv}
// 	helper2.dao = helper.dao
// 	return &helper2
// }

// ID get worker's id
func (helper *Helper) ID() string {
	return helper.job.ID
}

// ID get worker's space
func (helper *Helper) Space() string {
	return helper.space
}

// Job get worker's Job
func (helper *Helper) Job() job.Job {
	return helper.job
}

// Config get node's config
func (helper *Helper) Config() common.DaemonConfig {
	return helper.config
}

// Config get worker's Repository
func (helper *Helper) GetRepository() Repository {
	return helper.dao
}

// Info log info
func (helper *Helper) Info(msg string, keyvals ...interface{}) {
	helper.logger.Info(msg, keyvals...)
}

// Info log info
func (helper *Helper) Debug(msg string, keyvals ...interface{}) {
	helper.logger.Debug(msg, keyvals...)
}

// Info log info
func (helper *Helper) Error(msg string, keyvals ...interface{}) {
	helper.logger.Error(msg, keyvals...)
}

// PutCheckpoint ..
func (helper *Helper) PutCheckpoint(checkpoint interface{}) error {
	return helper.dao.PutCheckpoint(helper.ID(), checkpoint)
}

// GetCheckpoint ..
func (helper *Helper) GetCheckpoint(checkpoint interface{}) error {
	return helper.dao.GetCheckpoint(helper.ID(), checkpoint)
}

// PutData ..
func (helper *Helper) PutData(topic string, rowID string, data []byte) error {
	return helper.dao.PutData(helper.space, helper.ID(), topic, rowID, data)
}

// PutDataFullPath ..
func (helper *Helper) PutDataFullPath(fullPath string, data []byte) error {
	return helper.dao.PutDataFullPath(helper.space, fullPath, data)
}

// PutObject ..
func (helper *Helper) PutObject(topic string, rowID string, data interface{}) error {
	return helper.dao.PutObject(helper.space, helper.ID(), topic, rowID, data)
}

// GetData ..
func (helper *Helper) GetData(topic string, rowID string) (data []byte, err error) {
	return helper.dao.GetData(helper.space, helper.ID(), topic, rowID)
}

// GetObject ..
func (helper *Helper) GetObject(topic string, rowID string, ptr interface{}) error {
	return helper.dao.GetObject(helper.space, helper.ID(), topic, rowID, ptr)
}

// GetDataList ..
func (helper *Helper) GetDataList(topic string, handler DataHandler) error {
	return helper.dao.GetDataWithTopic(helper.space, helper.ID(), topic, handler)
}

// DeleteData ..
func (helper *Helper) DeleteData(topic string, rowID string) error {
	return helper.dao.DeleteData(helper.space, helper.ID(), topic, rowID)
}

type factoryRegistry struct {
	sync.Mutex
	workerFactories map[string]Factory
}

// NewAbstractWorkerFactory create AbstractWorkerFactory
func NewFactoryRegistry() (factory *factoryRegistry) {
	factory = &factoryRegistry{}
	factory.workerFactories = make(map[string]Factory)
	return factory
}

// AddFactory add worker factory
func (registry *factoryRegistry) RegisterFactory(factory Factory) error {
	registry.Lock()
	defer registry.Unlock()
	
	_, ok := registry.workerFactories[factory.Name()]
	if ok {
		return errors.New(fmt.Sprintf("Worker factory[%s] is already registered.", factory.Name()))
	}
	registry.workerFactories[factory.Name()] = factory
	return nil
}

// GetFactory get worker factory
func (registry *factoryRegistry) GetFactory(name string) (factory Factory, err error) {
	registry.Lock()
	defer registry.Unlock()
	factory = registry.workerFactories[name]
	if factory == nil {
		err = errors.New("Factory not found for " + name)
	}
	return factory, err
}

// // NewWorker implements worker.Factory.NewWorker
// func (registry *factoryRegistry) NewWorker(helper *Helper) (worker Worker, err error) {
// 	job := helper.Job()
// 	if len(job.FactoryName) == 0 {
// 		err = errors.New("job must have FactoryName")
// 		return nil, err
// 	}
//
// 	factory, err := registry.GetFactory(job.FactoryName)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return factory.NewWorker(helper)
//
// }
