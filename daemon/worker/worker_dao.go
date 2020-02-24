package worker

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	
	"github.com/tendermint/tendermint/libs/log"
	
	"github.com/rhizome-chain/tendermint-daemon/daemon/common"
	"github.com/rhizome-chain/tendermint-daemon/types"
)

const (
	PathCheckpoint = "chkpnt"
	
	// PatternPathData : data/{jobID}/{topic}
	PatternPathJobTopicData = "data/%s/%s"
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
	return dao.client.BroadcastTxSync(msg)
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

// GetDataWithTopic ..
func (dao *workerDao) GetDataWithTopic(space string, jobID string, topic string, handler DataHandler) error {
	fullPath := makeDataPath(jobID, topic)
	msg := types.NewViewMsgMany(space, fullPath, "", "")
	
	err := dao.client.GetMany(msg, func(key []byte, value []byte) bool {
		fmt.Println("key=", string(key), "value=", string(value))
		// handler(jobID,topic,rowID, value)
		return true
	})
	
	return err
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
