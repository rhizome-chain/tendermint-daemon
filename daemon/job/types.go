package job

import (
	"encoding/json"
	
	"github.com/rhizome-chain/tendermint-daemon/types"
	
	"github.com/google/uuid"
)

type Repository interface {
	PutJob(job Job) error
	PutJobIfNotExist(job Job) error
	RemoveJob(jobID string) error
	RemoveAllJobs() error
	GetJob(jobID string) (job Job, err error)
	ContainsJob(jobID string) bool
	GetAllJobIDs() (jobIDs []string, err error)
	GetAllJobs() (jobs map[string]Job, err error)
	GetMemberJobIDs(membID string) (jobIDs []string, err error)
	GetAllMemberJobIDs() (membJobMap map[string][]string, err error)
	PutMemberJobIDs(membID string, jobIDs []string) (err error)
	GetMemberJobs(membID string) (jobs []Job, err error)
	Commit() error
}

func init() {
	types.BasicCdc.RegisterConcrete(Job{}, "daemon/job", nil)
}

// Job job data structure
type Job struct {
	FactoryName string
	ID          string
	Data        []byte
}

// NewJob create new job with uuid
func New(factoryName string, data []byte) Job {
	uuid := uuid.New()
	return Job{FactoryName: factoryName, ID: uuid.String(), Data: data}
}

// NewWithID create new job with pi
func NewWithID(factory string, jobID string, data []byte) Job {
	return Job{FactoryName: factory, ID: jobID, Data: data}
}

// GetAsString Get data as string
func (job *Job) GetAsString() string {
	return string(job.Data)
}

// GetAsObject Get data as interface
func (job *Job) GetAsObject(obj interface{}) error {
	err := json.Unmarshal(job.Data, &obj)
	return err
}
