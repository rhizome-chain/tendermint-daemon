package client

import (
	"encoding/json"
	"errors"
	"time"
	
	"github.com/tendermint/tendermint/mempool"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/bytes"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/rpc/core"
	rpctypes "github.com/tendermint/tendermint/rpc/lib/types"
	tmtypes "github.com/tendermint/tendermint/types"
	
	"github.com/rhizome-chain/tendermint-daemon/tm/tmcom"
	"github.com/rhizome-chain/tendermint-daemon/types"
)

func IsErrMempoolIsFull(err error) bool {
	_, ok := err.(mempool.ErrMempoolIsFull)
	return ok
}

func IsErrTxInCache(err error) bool {
	return err == mempool.ErrTxInCache
}

type TMClient struct {
	config *cfg.Config
	logger log.Logger
	codec  *types.Codec
	ctx    *rpctypes.Context
}

var _ types.Client = (*TMClient)(nil)

func NewClient(config *cfg.Config, logger log.Logger, codec *types.Codec) types.Client {
	return &TMClient{config, logger, codec, &rpctypes.Context{}}
}

func (client *TMClient) broadcastTx(funcTx func() (*ctypes.ResultBroadcastTx, error)) (err error) {
	_, err = funcTx()
	if err != nil {
		for IsErrMempoolIsFull(err) {
			client.logger.Info("[TMClient] Wait 3sec... ", "err", err)
			time.Sleep(3 * time.Second)
			_, err = funcTx()
		}
		
		return err
	}
	return err
}

func (client *TMClient) BroadcastTxSync(msg *types.TxMsg) (err error) {
	msgBytes, err := types.EncodeTxMsg(msg)
	if err != nil {
		client.logger.Error("[TMClient] BroadcastTxSync : marshal", err)
		return err
	}
	
	err = client.broadcastTx(func() (*ctypes.ResultBroadcastTx, error) {
		return core.BroadcastTxSync(client.ctx, msgBytes)
	})
	
	if err != nil {
		client.logger.Error("[TMClient] BroadcastTxSync ", err)
	}
	return err
}

func (client *TMClient) BroadcastTxAsync(msg *types.TxMsg) (err error) {
	msgBytes, err := types.EncodeTxMsg(msg)
	if err != nil {
		client.logger.Error("[TMClient] BroadcastTxAsync : EncodeTxMsg", err)
		return err
	}
	
	err = client.broadcastTx(func() (*ctypes.ResultBroadcastTx, error) {
		return core.BroadcastTxAsync(client.ctx, msgBytes)
	})
	
	if err != nil {
		client.logger.Error("[TMClient] BroadcastTxAsync", err)
	}
	return err
}

func (client *TMClient) BroadcastTxCommit(msg *types.TxMsg) (err error) {
	msgBytes, err := types.EncodeTxMsg(msg)
	if err != nil {
		client.logger.Error("[TMClient] BroadcastTxCommit : EncodeTxMsg", err)
		return err
	}
	
	_, err = core.BroadcastTxCommit(client.ctx, msgBytes)
	
	if err != nil && !IsErrTxInCache(err) {
		client.logger.Error("[TMClient] BroadcastTxCommit", err)
	}
	return err
}

func (client *TMClient) Commit() (err error) {
	// _, err = core.Commit(client.ctx, nil)
	// if err != nil {
	// 	client.logger.Error("[TMClient] Commit : ", err)
	// }
	msgBytes, err := types.EncodeTxMsg(&types.TxMsg{Type: types.TxCommit, Space: tmcom.SpaceDaemonState})
	_, err = core.BroadcastTxCommit(client.ctx, msgBytes)
	if err != nil {
		client.logger.Error("[TMClient] BroadcastTxCommit", err)
		return err
	}
	return err
}

func (client *TMClient) Has(msg *types.ViewMsg) (ok bool, err error) {
	if msg.Type != types.Has {
		return ok, errors.New("[TMClient] Has needs ViewType Has")
	}
	
	data, err := client.Query(msg)
	
	if err != nil {
		client.logger.Error("[TMClient] Has : ", err)
		return ok, err
	}
	
	err = client.UnmarshalObject(data, &ok)
	
	if err != nil {
		client.logger.Error("[TMClient] Has Unmarshal : ", err)
		return ok, err
	}
	return ok, err
}

func (client *TMClient) Query(msg *types.ViewMsg) (data []byte, err error) {
	msgBytes, err := types.EncodeViewMsg(msg)
	if err != nil {
		client.logger.Error("[TMClient] Query : EncodeViewMsg", err)
		return data, err
	}
	
	res, err := core.ABCIQuery(client.ctx, msg.Path, bytes.HexBytes(msgBytes), 0, false)
	
	if err != nil {
		client.logger.Error("[TMClient] Query : ABCIQuery", err)
	}
	
	return res.Response.Value, err
}

func (client *TMClient) GetObject(msg *types.ViewMsg, obj interface{}) (err error) {
	if msg.Type != types.GetOne {
		return errors.New("[TMClient] GetObject needs ViewType GetOne")
	}
	data, err := client.Query(msg)
	
	if err != nil {
		client.logger.Error("[TMClient] GetObject : ", err)
		return err
	}
	
	if len(data) == 0 {
		return types.NewNoDataError()
	}
	
	err = client.UnmarshalObject(data, obj)
	
	if err != nil {
		client.logger.Error("[TMClient] GetObject Unmarshal : ", err)
		return err
	}
	return err
}

func (client *TMClient) GetMany(msg *types.ViewMsg, handler types.KVHandler) (err error) {
	if msg.Type != types.GetMany {
		return errors.New("[TMClient] GetMany needs ViewType GetMany")
	}
	data, err := client.Query(msg)
	
	if err != nil {
		client.logger.Error("[TMClient] GetMany : ", err)
		return err
	}
	
	if len(data) == 0 {
		return types.NewNoDataError()
	}
	
	kvArray := []types.KeyValue{}
	err = client.UnmarshalObject(data, &kvArray)
	
	if err != nil {
		client.logger.Error("[TMClient] GetMany Unmarshal : ", err)
		return err
	}
	
	for _, kv := range kvArray {
		if !handler(kv.Key, kv.Value) {
			break
		}
	}
	return err
}

func (client *TMClient) GetKeys(msg *types.ViewMsg) (keys []string, err error) {
	if msg.Type != types.GetKeys {
		return nil, errors.New("[TMClient] GetKeys needs ViewType GetKeys")
	}
	
	data, err := client.Query(msg)
	
	if err != nil {
		client.logger.Error("[TMClient] GetKeys : ", err)
		return nil, err
	}
	
	if len(data) == 0 {
		return nil, types.NewNoDataError()
	}
	
	keys = []string{}
	err = client.UnmarshalObject(data, &keys)
	
	if err != nil {
		client.logger.Error("[TMClient] GetKeys Unmarshal : ", err)
		return nil, err
	}
	
	return keys, err
}

// func (client *TMClient) Subscribe(msg *types.ViewMsg) (keys []string, err error){
// 	core.Subscribe(client.ctx, )
// }

func (client *TMClient) UnmarshalObject(bz []byte, ptr interface{}) error {
	if len(bz) == 0 {
		return types.NewNoDataError()
	}
	
	err := client.codec.UnmarshalBinaryBare(bz, ptr)
	
	if err != nil {
		err := errors.New("[TMClient] UnmarshalObject : " + err.Error())
		return err
	}
	return err
}

func (client *TMClient) MarshalObject(ptr interface{}) (bytes []byte, err error) {
	bytes, err = client.codec.MarshalBinaryBare(ptr)
	
	if err != nil {
		err := errors.New("[TMClient] MarshalObject : " + err.Error())
		return nil, err
	}
	return bytes, err
}

func (client *TMClient) UnmarshalJson(bz []byte, ptr interface{}) error {
	if len(bz) == 0 {
		return types.NewNoDataError()
	}
	
	err := json.Unmarshal(bz, ptr)
	
	if err != nil {
		err := errors.New("[TMClient] UnmarshalJson : " + err.Error())
		return err
	}
	return err
}

func (client *TMClient) MarshalJson(ptr interface{}) (bytes []byte, err error) {
	bytes, err = json.Marshal(ptr)
	
	if err != nil {
		err := errors.New("[TMClient] MarshalJson : " + err.Error())
		return nil, err
	}
	return bytes, err
}

func (client *TMClient) CurrentBlockNumber() (block int64) {
	blockRes, err := core.Block(client.ctx, nil)
	if err != nil {
		client.logger.Error("[TMClient] CurrentBlockNumber", err)
		return 0
	}
	return blockRes.Block.Height
}

func (client *TMClient) GetValidators() (validators []*tmtypes.Validator) {
	valRes, err := core.Validators(client.ctx,nil,0,100)
	if err != nil {
		client.logger.Error("[TMClient] GetValidators", err)
	} else {
		validators = valRes.Validators
	}
	return validators
}

func (client *TMClient) GetPeerIDs() (peerIDs []string) {
	result, err := core.NetInfo(client.ctx)
	
	if err != nil {
		client.logger.Error("[TMClient] GetValidators", err)
		return peerIDs
	}
	
	peerIDs = []string{}
	
	for _,p := range result.Peers{
		peerIDs = append(peerIDs, string(p.NodeInfo.ID()))
	}
	return peerIDs
}


