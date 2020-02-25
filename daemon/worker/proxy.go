package worker

import (
	"fmt"
	
	"github.com/rhizome-chain/tendermint-daemon/daemon/job"
)

// Proxy for worker
type Proxy interface {
	GetJob() job.Job
	GetCheckpoint(checkpoint interface{}) error
	GetData(topic string, rowID string) (data []byte, err error)
	GetObject(topic string, rowID string, data interface{}) error
	GetDataList(topic string, handler DataHandler) error
	GetDataListRange(topic string, from string, end string, handler DataHandler) error
	DeleteData(topic string, rowID string) error
	SubscribeTx(topic string, from string, handler DataHandler) (CancelSubs, error)
	CollectAndSubscribe(topic string, from string, handler DataHandler) (CancelSubs, error)
}

type ProxyProvider interface {
	NewWorkerProxy(jobID string) (Proxy, error)
}

type BaseProxy struct {
	job    job.Job
	helper *Helper
}

func NewBaseProxy(job job.Job, helper *Helper) *BaseProxy {
	return &BaseProxy{
		job:    job,
		helper: helper,
	}
}

func (b BaseProxy) GetJob() job.Job {
	return b.job
}

func (b BaseProxy) GetCheckpoint(checkpoint interface{}) error {
	return b.helper.GetCheckpoint(checkpoint)
}

func (b BaseProxy) GetData(topic string, rowID string) (data []byte, err error) {
	return b.helper.GetData(topic, rowID)
}

func (b BaseProxy) GetObject(topic string, rowID string, ptr interface{}) error {
	return b.helper.GetObject(topic, rowID, ptr)
}

func (b BaseProxy) GetDataList(topic string, handler DataHandler) error {
	return b.helper.GetDataList(topic, handler)
}

func (b BaseProxy) GetDataListRange(topic string, from string, end string, handler DataHandler) error {
	return b.helper.GetDataListRange(topic, from, end, handler)
}

func (b BaseProxy) DeleteData(topic string, rowID string) error {
	return b.helper.DeleteData(topic, rowID)
}

func (b BaseProxy) SubscribeTx(topic string, from string, handler DataHandler) (CancelSubs, error) {
	return b.helper.SubscribeTx(topic, from, handler)
}

func (b BaseProxy) CollectAndSubscribe(topic string, from string, handler DataHandler) (CancelSubs, error) {
	var lastRow string
	err := b.GetDataListRange(topic, from, "", func(jobID string, topic string, rowID string, value []byte) bool {
		lastRow = rowID
		handler(jobID, topic, rowID, value)
		return true
	})
	
	if err != nil {
		return nil, err
	}
	b.helper.Info(fmt.Sprintf("[Proxy:%s] collect from %s to %s", b.job.ID, from, lastRow))
	
	cancel, err := b.SubscribeTx("in", lastRow+"!", handler)
	return cancel, err
}

var _ Proxy = (*BaseProxy)(nil)
