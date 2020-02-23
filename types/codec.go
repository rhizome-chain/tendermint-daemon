package types

import (
	"github.com/tendermint/go-amino"
)

type Codec = amino.Codec

// NewCodec
func NewCodec() *Codec {
	return amino.NewCodec()
}

// BasicCdc is the basic codec
var BasicCdc = NewCodec()

func init() {
	RegisterCodec(BasicCdc)
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *Codec) {
	cdc.RegisterConcrete(TxMsg{}, "daemon/txmsg", nil)
	cdc.RegisterConcrete(ViewMsg{}, "daemon/viewmsg", nil)
}
