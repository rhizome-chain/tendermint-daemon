package worker

import (
	"github.com/rhizome-chain/tendermint-daemon/daemon/job"
)


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

func (b BaseProxy) DeleteData(topic string, rowID string) error {
	return b.helper.DeleteData(topic, rowID)
}

var _ Proxy = (*BaseProxy)(nil)
