package worker

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

const (
	SpaceDefaultWorker = "daemon-data"
)

type DataHandler func(jobID string, topic string, rowID string, value []byte) bool

type CancelSubs func()

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
	GetDataWithTopicRange(space string, jobID string, topic string, from string, end string, handler DataHandler) error
	// SubscribeTx async subscription
	SubscribeTx(space string, jobID string, topic string, from string, handler DataHandler) (CancelSubs, error)
	
	
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

type DataPath struct {
	JobID string
	Topic string
}

func (path DataPath) String() string {
	return makeDataPath(path.JobID, path.Topic)
}

func makeDataPath(jobID string, topic string) string {
	return fmt.Sprintf(PatternPathJobTopicData, jobID, topic)
}

func ParseDataPathBytes(path []byte) (dataPath DataPath, err error) {
	paths := bytes.Split(path, []byte("/"))
	if len(paths) != 2 {
		return dataPath, errors.New(fmt.Sprintf("Illegal DataPath format %s", path))
	}
	dataPath = DataPath{JobID: string(paths[0]), Topic: string(paths[1])}
	return dataPath, err
}

func ParseDataPath(path string) (dataPath DataPath, err error) {
	paths := strings.Split(path, "/")
	if len(paths) != 2 {
		return dataPath, errors.New(fmt.Sprintf("Illegal DataPath format %s", path))
	}
	dataPath = DataPath{JobID: paths[0], Topic: paths[1]}
	return dataPath, err
}
