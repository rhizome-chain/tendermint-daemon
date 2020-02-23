package types

import (
	"time"
)

type ViewType uint8
type TxType uint8

var (
	LastKeySuffix      = "~~~"
	LastKeySuffixBytes = []byte(LastKeySuffix)
)

const (
	Has     = ViewType(0)
	GetOne  = ViewType(1)
	GetMany = ViewType(2)
	GetKeys = ViewType(3)
)

const (
	TxSet        = TxType(11)
	TxSetSync    = TxType(12)
	TxDelete     = TxType(21)
	TxDeleteSync = TxType(22)
	TxCommit     = TxType(90)
)

func TxTypeString(typ TxType) string {
	switch typ {
	case TxDelete:
		return "Delete"
	case TxDeleteSync:
		return "Delete"
	case TxSetSync:
		return "Set"
	case TxSet:
		return "Set"
	}
	
	return ""
}

// const (
// 	MethodPutHeartbeat = "daemon/heartbeat"
// )

type TxMsg struct {
	Type   TxType
	Space  string
	Path   string
	Key    []byte
	Data   []byte
	TxHash int64
}

func (msg *TxMsg) GetEventKey() string {
	return msg.Space + "." + msg.Path
}

func (msg *TxMsg) GetTypeString() string {
	return TxTypeString(msg.Type)
}

// to avoid 'tx already exists in cache' problem
func (msg *TxMsg) SetTxHash() {
	msg.TxHash = time.Now().UnixNano()
}

func NewTxMsg(typ TxType, space string, path string, key string, data []byte) *TxMsg {
	msg := TxMsg{Type: typ, Space: space, Path: path, Key: []byte(key), Data: data}
	return &msg
}

// func EncodeTx(typ TxType, space string, path string, key string, data []byte) ([]byte, error) {
// 	msg := TxMsg{Type: typ, Space: space, Path: path, Key: []byte(key), Data: data}
// 	return BasicCdc.MarshalBinaryBare(msg)
// }

func EncodeTxMsg(msg *TxMsg) ([]byte, error) {
	return BasicCdc.MarshalBinaryBare(msg)
}

func DecodeTxMsg(msgBytes []byte) (*TxMsg, error) {
	msg := &TxMsg{}
	err := BasicCdc.UnmarshalBinaryBare(msgBytes, msg)
	return msg, err
}

type ViewMsg struct {
	Type  ViewType
	Space string
	Path  string
	Start []byte
	End   []byte
}

func NewViewMsgHas(space string, path string, key string) *ViewMsg {
	msg := ViewMsg{Type: Has, Space: space, Path: path, Start: []byte(key)}
	return &msg
}

func NewViewMsgOne(space string, path string, key string) *ViewMsg {
	msg := ViewMsg{Type: GetOne, Space: space, Path: path, Start: []byte(key)}
	return &msg
}

func NewViewMsgMany(space string, path string, start string, end string) *ViewMsg {
	return newViewMsgMany(GetMany, space, path, start, end)
}

func NewViewMsgKeys(space string, path string, start string, end string) *ViewMsg {
	return newViewMsgMany(GetKeys, space, path, start, end)
}

func newViewMsgMany(typ ViewType, space string, path string, start string, end string) *ViewMsg {
	var endBytes []byte
	if len(end) > 0 {
		endBytes = []byte(end)
	}
	msg := ViewMsg{Type: typ, Space: space, Path: path, Start: []byte(start), End: endBytes}
	return &msg
}

func EncodeViewMsg(msg *ViewMsg) ([]byte, error) {
	return BasicCdc.MarshalBinaryBare(msg)
}

func DecodeViewMsg(msgBytes []byte) (*ViewMsg, error) {
	msg := &ViewMsg{}
	err := BasicCdc.UnmarshalBinaryBare(msgBytes, msg)
	return msg, err
}
