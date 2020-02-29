package tm

import (
	"fmt"
	"github.com/rhizome-chain/tendermint-daemon/tm/events"
	
	cfg "github.com/tendermint/tendermint/config"
	
	"github.com/tendermint/tendermint/abci/example/code"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	
	"github.com/rhizome-chain/tendermint-daemon/types"
)

type DaemonApp struct {
	*BaseApplication
}

var _ abcitypes.Application = (*DaemonApp)(nil)

func NewDaemonApplication(config *cfg.Config, logger log.Logger) (dapp *DaemonApp) {
	baseapp := NewBaseApplication(config, logger)
	dapp = &DaemonApp{BaseApplication: baseapp}
	
	return dapp
}

func (app *DaemonApp) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
	if app.isValidatorTx(req.Tx) {
		return app.execValidatorTx(req.Tx)
	}
	
	msg, err := types.DecodeTxMsg(req.Tx)
	
	if err != nil {
		app.logger.Error("[DMA]DeliverTx", err)
		return abcitypes.ResponseDeliverTx{Code: code.CodeTypeEncodingError}
	} else {
		store := app.getSpaceStoreAny(msg.Space, msg.Path)
		
		switch msg.Type {
		case types.TxSet:
			store.Set(msg.Key, msg.Data)
		case types.TxSetSync:
			store.SetSync(msg.Key, msg.Data)
		case types.TxDelete:
			store.Delete(msg.Key)
		case types.TxDeleteSync:
			store.DeleteSync(msg.Key)
		case types.TxCommit:
			// DO NOTHING
		default:
			app.logger.Error("[DMA] Unknown TxMsg ", "type", msg.Type)
			return abcitypes.ResponseDeliverTx{Code: code.CodeTypeUnknownError}
		}
		
		app.IncreaseTxSize()
		events.PublishTxEvent(*msg)
		return abcitypes.ResponseDeliverTx{Code: code.CodeTypeOK}
	}
}

func (app *DaemonApp) Query(reqQuery abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery) {
	msg, err := types.DecodeViewMsg(reqQuery.Data)
	if err != nil {
		app.logger.Error("[DMA] DecodeViewMsg ", err)
		resQuery.Log = err.Error()
		return resQuery
	}
	
	resQuery.Index = -1
	resQuery.Key = reqQuery.Data
	
	// store, err := app.getSpace(msg.Space).GetStore(msg.Path)
	//
	// if err != nil {
	// 	app.logger.Error("[DMA] Unknown Store "+reqQuery.Path, err)
	// 	resQuery.Log = err.Error()
	// 	return resQuery
	// }
	
	store := app.getSpaceStoreAny(msg.Space, msg.Path)
	
	if msg.Type == types.Has {
		ok, err := store.Has(msg.Start)
		
		if err != nil {
			app.logger.Error("[DMA] Query Has "+string(msg.Start), err)
			resQuery.Log = err.Error()
			return resQuery
		}
		
		if ok {
			resQuery.Log = "exists"
		} else {
			resQuery.Log = "does not exist"
		}
		
		bytes, err := types.BasicCdc.MarshalBinaryBare(ok)
		
		if err != nil {
			app.logger.Error("[DMA] Query Has Unmarshal", err)
			resQuery.Log = err.Error()
			return resQuery
		}
		
		resQuery.Value = bytes
	} else if msg.Type == types.GetOne {
		bytes, err := store.Get(msg.Start)
		if err != nil {
			app.logger.Error("[DMA] Query GetOne "+string(msg.Start), err)
			resQuery.Log = err.Error()
			return resQuery
		}
		
		if bytes == nil {
			resQuery.Log = "does not exist"
		} else {
			resQuery.Log = "exists"
		}
		
		resQuery.Value = bytes
	} else if msg.Type == types.GetMany {
		bytes, err := store.GetMany(msg.Start, msg.End)
		
		if err != nil {
			app.logger.Error(fmt.Sprintf("[DMA] Query GetMany %s - %s ", msg.Start, msg.End), err)
			resQuery.Log = err.Error()
			return resQuery
		}
		
		if bytes == nil {
			resQuery.Log = "does not exist"
		} else {
			resQuery.Log = "exists"
		}
		
		resQuery.Value = bytes
	} else if msg.Type == types.GetKeys {
		bytes, err := store.GetKeys(msg.Start, msg.End)
		
		if err != nil {
			app.logger.Error(fmt.Sprintf("[DMA] Query GetKeys %s - %s ", msg.Start, msg.End), err)
			resQuery.Log = err.Error()
			return resQuery
		}
		
		if bytes == nil {
			resQuery.Log = "does not exist"
		} else {
			resQuery.Log = "exists"
		}
		
		resQuery.Value = bytes
	} else {
		app.logger.Error(fmt.Sprintf("[DMA] Unknown ViewType %d ", msg.Type), err)
	}
	
	return resQuery
}
