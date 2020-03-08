package worker

import (
	"fmt"
	
	"github.com/google/uuid"
	
	"github.com/rhizome-chain/tendermint-daemon/tm/events"
	
	"github.com/tendermint/tendermint/libs/log"
	
	"github.com/rhizome-chain/tendermint-daemon/daemon/common"
	"github.com/rhizome-chain/tendermint-daemon/types"
)

const (
	PathCheckpoint = "chkpnt"
	
	// PatternPathData : data/{jobID}/{topic}
	PatternPathJobTopicData = "%s/%s"
)

// workerDao kv store model for cluster
type workerDao struct {
	config common.DaemonConfig
	logger log.Logger
	client types.Client
}

var _ Repository = (*workerDao)(nil)

func NewRepository(config common.DaemonConfig, logger log.Logger, client types.Client) Repository {
	dao := &workerDao{config: config, logger: logger, client: client}
	return dao
}

// PutCheckpoint marshal checkpoint to json
func (dao *workerDao) PutCheckpoint(jobID string, checkpoint interface{}) error {
	jobIDsBytes, err := dao.client.MarshalJson(checkpoint)
	if err != nil {
		return err
	}
	msg := types.NewTxMsg(types.TxSet, common.SpaceDaemon, PathCheckpoint, jobID, jobIDsBytes)
	err = dao.client.BroadcastTxSync(msg)
	if err != nil {
		dao.logger.Error(fmt.Sprintf("[Worker - %s] PutCheckpoint %v",jobID, checkpoint), err)
	}
	return err
}

// GetCheckpoint  unmarshal checkpoint from json
func (dao *workerDao) GetCheckpoint(jobID string, checkpoint interface{}) error {
	msg := types.NewViewMsgOne(common.SpaceDaemon, PathCheckpoint, jobID)
	bytes, err := dao.client.Query(msg)
	if err != nil {
		return err
	}
	err = dao.client.UnmarshalJson(bytes, &checkpoint)
	return err
}

//func (dao *workerDao) CurrentBlockNumber() (block int64) {
//	return dao.client.CurrentBlockNumber()
//}

// PutData put data to {space}
func (dao *workerDao) PutData(space string, jobID string, topic string, rowID string, data []byte) error {
	fullPath := makeDataPath(jobID, topic)
	msg := types.NewTxMsg(types.TxSet, space, fullPath, rowID, data)
	return dao.client.BroadcastTxAsync(msg)
}

// PutObject data type must be registered to Codec
func (dao *workerDao) PutObject(space string, jobID string, topic string, rowID string, data interface{}) error {
	dataBytes, err := dao.client.MarshalObject(data)
	if err != nil {
		dao.logger.Error("[ERROR-WorkerDao] PutObject marshal", err)
	}
	// fmt.Println("&&&&& PutObject :: key=", key, ", data=", data)
	return dao.PutData(space, jobID, topic, rowID, dataBytes)
}

// PutData put data to {space}
func (dao *workerDao) PutDataSync(space string, jobID string, topic string, rowID string, data []byte) error {
	fullPath := makeDataPath(jobID, topic)
	msg := types.NewTxMsg(types.TxSetSync, space, fullPath, rowID, data)
	return dao.client.BroadcastTxSync(msg)
}

// PutObject data type must be registered to Codec
func (dao *workerDao) PutObjectSync(space string, jobID string, topic string, rowID string, data interface{}) error {
	dataBytes, err := dao.client.MarshalObject(data)
	if err != nil {
		dao.logger.Error("[ERROR-WorkerDao] PutObjectSync marshal", err)
	}
	// fmt.Println("&&&&& PutObject :: key=", key, ", data=", data)
	return dao.PutDataSync(space, jobID, topic, rowID, dataBytes)
}

// GetData ..
func (dao *workerDao) GetData(space string, jobID string, topic string, rowID string) (data []byte, err error) {
	fullPath := makeDataPath(jobID, topic)
	msg := types.NewViewMsgOne(space, fullPath, rowID)
	bytes, err := dao.client.Query(msg)
	if err != nil {
		return nil, err
	}
	return bytes, err
}

// GetObject data type must be registered to Codec
func (dao *workerDao) GetObject(space string, jobID string, topic string, rowID string, data interface{}) error {
	bytes, err := dao.GetData(space, jobID, topic, rowID)
	if err != nil {
		return err
	}
	return dao.client.UnmarshalObject(bytes, data)
}

// DeleteData ..
func (dao *workerDao) DeleteData(space string, jobID string, topic string, rowID string) error {
	fullPath := makeDataPath(jobID, topic)
	msg := types.NewTxMsg(types.TxDelete, space, fullPath, rowID, nil)
	err := dao.client.BroadcastTxSync(msg)
	if err != nil {
		dao.logger.Error("[ERROR-WorkerDao] DeleteData ", err)
	}
	
	return err
}

// DeleteData ..
func (dao *workerDao) DeleteDataByPrefix(space string, jobID string, topic string, prefix string) error {
	fullPath := makeDataPath(jobID, topic)
	msg := types.NewTxMsg(types.TxDeleteByPrefix, space, fullPath, prefix, nil)
	err := dao.client.BroadcastTxSync(msg)
	if err != nil {
		dao.logger.Error("[ERROR-WorkerDao] DeleteDataByPrefix ", err)
	}
	
	return err
}

// GetDataWithTopic ..
func (dao *workerDao) GetDataWithTopic(space string, jobID string, topic string, handler DataHandler) error {
	fullPath := makeDataPath(jobID, topic)
	msg := types.NewViewMsgMany(space, fullPath, "", "")
	
	err := dao.client.GetMany(msg, func(key []byte, value []byte) bool {
		// fmt.Println("key=", string(key), "value=", string(value))
		handler(jobID, topic, string(key), value)
		return true
	})
	
	return err
}

// GetDataWithTopic ..
func (dao *workerDao) GetDataWithTopicRange(space string, jobID string, topic string, from string, end string, handler DataHandler) error {
	fullPath := makeDataPath(jobID, topic)
	
	if len(end) >0 && from > end {
		dao.logger.Error(fmt.Sprintf("[WorkerDao] GetDataWithTopicRange from is larger than end : from=%s, end=%s",from,end))
		dao.logger.Error(fmt.Sprintf("          - GetDataWithTopicRange : set from=end"))
		from = end
	}
	
	msg := types.NewViewMsgMany(space, fullPath, from, end)
	
	err := dao.client.GetMany(msg, func(key []byte, value []byte) bool {
		//fmt.Println("GetDataWithTopicRange key=", string(key), "value=", string(value))
		handler(jobID, topic, string(key), value)
		return true
	})
	
	if err != nil && types.IsNoDataError(err){
		return nil
	}
	
	return err
}

type CancelTxSubs struct {
	eventPath types.EventPath
	name      string
}

// SubscribeTx async subscription
func (dao *workerDao) SubscribeTx(space string, jobID string, topic string, from string, handler DataHandler) (CancelSubs, error) {
	cancel, rowID, err := dao.innerSubscribeTx(space, jobID, topic)
	if err != nil {
		return nil, err
	}
	
	lastRow := <-rowID
	
	fmt.Println(" - WorkerDao : SubscribeTx GetDataWithTopicRange :: " , " from=",from, ",lastRow=", lastRow)
	
	err = dao.GetDataWithTopicRange(space, jobID, topic, from, lastRow, handler)
	if err != nil {
		return nil, err
	}
	
	type delegate struct {
		running bool
	}
	
	d := &delegate{running: true}
	
	go func(d *delegate){
		for d.running {
			row := <-rowID
			data, err := dao.GetData(space, jobID, topic, row)
			if err != nil {
				dao.logger.Error("SubscribeTx - get data ", "rowID", row)
			} else {
				handler(jobID, topic, row, data)
			}
		}
	}(d)
	
	
	return func() {
		cancel()
		d.running = false
		dao.logger.Info(fmt.Sprintf("Cancel subscribing [worker:%s]",jobID ))
	}, nil
}

// innerSubscribeTx ..
func (dao *workerDao) innerSubscribeTx(space string, jobID string, topic string) (CancelSubs, chan string, error) {
	evtPath := events.MakeTxEventPath(space, jobID, topic)
	name := uuid.New().String()
	
	var rowID = make(chan string)
	err := events.SubscribeTxEvent(evtPath, name, func(event events.TxEvent) {
		rowID <- string(event.Key)
	})
	
	if err != nil {
		dao.logger.Error(fmt.Sprintf("WorkerDAO[%s] innerSubscribeTx",jobID), err)
		return nil, nil, err
	}
	
	return func() {
		events.UnsubscribeTxEvent(evtPath, name)
		dao.logger.Info(fmt.Sprintf("Unsubscribe [worker:%s] path:%s",jobID, evtPath))
	}, rowID, nil
}

// PutDataFullPath ..
func (dao *workerDao) PutDataFullPath(space string, fullPath string, data []byte) error {
	msg := types.NewTxMsg(types.TxSet, space, fullPath, "", data)
	return dao.client.BroadcastTxAsync(msg)
}

// PutObjectFullPath ..
func (dao *workerDao) PutObjectFullPath(space string, fullPath string, data interface{}) error {
	dataBytes, err := dao.client.MarshalObject(data)
	if err != nil {
		dao.logger.Error("[ERROR-WorkerDao] PutObjectFullPath marshal", err)
	}
	return dao.PutDataFullPath(space, fullPath, dataBytes)
}

// GetData ..
func (dao *workerDao) GetDataFullPath(space string, fullPath string) (data []byte, err error) {
	msg := types.NewViewMsgOne(space, fullPath, "")
	bytes, err := dao.client.Query(msg)
	if err != nil {
		return nil, err
	}
	return bytes, err
}

// GetObjectFullPath data type must be registered to Codec
func (dao *workerDao) GetObjectFullPath(space string, fullPath string, data interface{}) error {
	bytes, err := dao.GetDataFullPath(space, fullPath)
	if err != nil {
		return err
	}
	return dao.client.UnmarshalObject(bytes, data)
}

// DeleteDataFullPath ..
func (dao *workerDao) DeleteDataFullPath(space string, fullPath string) error {
	msg := types.NewTxMsg(types.TxDelete, space, fullPath, "", nil)
	err := dao.client.BroadcastTxSync(msg)
	if err != nil {
		dao.logger.Error("[ERROR-WorkerDao] DeleteData ", err)
	}
	
	return err
}
