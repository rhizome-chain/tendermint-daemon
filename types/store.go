package types

import dbm "github.com/tendermint/tm-db"

type KeyValue struct {
	Key []byte
	Value []byte
}

type Store interface {
	GetPath() string
	Get(key []byte) ([]byte, error)
	Has(key []byte) (bool, error)
	Set(key []byte, value []byte) error
	SetSync(key []byte, value []byte) error
	Delete(key []byte) error
	DeleteSync(key []byte) error
	Iterator(start, end []byte) (dbm.Iterator, error)
	GetMany(start, end []byte) (kvArrayBytes[]byte, err error)
}

type SpaceRegistry interface {
	RegisterSpace(name string)
	RegisterSpaceIfNotExist(name string)
}

