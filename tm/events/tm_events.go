package events

import (
	"time"
	
	"github.com/rhizome-chain/tendermint-daemon/types"
)

const (
	EventScopeBlock = types.EventScope("tm-block")
	EventScopeTx    = types.EventScope("tm-tx")
	
	EndBlockEventPath   = types.EventPath("end-block")
	BeginBlockEventPath = types.EventPath("begin-block")
	CommitEventPath = types.EventPath("commit")
)

var (
	blockEventBus = types.RegisterEventBus(EventScopeBlock)
	txEventBus    = types.RegisterEventBus(EventScopeTx)
)

// func StartTMEventBus() {
// 	blockEventBus.Start()
// 	txEventBus.Start()
// }

func PublishBlockEvent(event BlockEvent) {
	blockEventBus.Publish(event)
}

func SubscribeBlockEvent(path types.EventPath, name string, handler types.EventHandler) error {
	return blockEventBus.Subscribe(path, name, handler)
}

func UnsubscribeBlockEvent(path types.EventPath, name string) {
	blockEventBus.Unsubscribe(path, name)
}

type BlockEvent types.Event

type BeginBlockEvent struct {
	BlockEvent
	Height int64
	Time   time.Time
}

func (event BeginBlockEvent) Path() types.EventPath { return BeginBlockEventPath }

type EndBlockEvent struct {
	BlockEvent
	Height int64
	Size   int
}

func (event EndBlockEvent) Path() types.EventPath { return EndBlockEventPath }


type CommitEvent struct {
	BlockEvent
	Height int64
	Size   int64
	AppHash []byte
}

func (event CommitEvent) Path() types.EventPath { return CommitEventPath }




type TxEvent struct {
	path types.EventPath
	Type types.TxType
	Key  []byte
}

type TxEventHandler func(event TxEvent)

func (event TxEvent) Path() types.EventPath { return event.path }

func MakeTxEventPath(space string, path string, prefix string)  types.EventPath {
	eventPath := types.EventPath(space + "!" + path + "/" + prefix)
	return eventPath
}

func NewTxEvent(msg types.TxMsg) (event TxEvent) {
	event = TxEvent{
		path: MakeTxEventPath(msg.Space, msg.Path, string(msg.Key)),
		Type: msg.Type,
		Key:  msg.Key,
	}
	
	return event
}


func PublishTxEvent(msg types.TxMsg) {
	event := NewTxEvent(msg)
	txEventBus.Publish(event)
}

func SubscribeTxEvent(eventPath types.EventPath, name string, handler TxEventHandler) error {
	wrap := func(event types.Event){
		txEvt := event.(TxEvent)
		handler(txEvt)
	}
	
	return txEventBus.Subscribe(eventPath, name, wrap)
}

func UnsubscribeTxEvent(eventPath types.EventPath , name string) {
	txEventBus.Unsubscribe(eventPath, name)
}
