package store

import (
	"bytes"
	"errors"
	"fmt"
	
	"github.com/rhizome-chain/tendermint-daemon/types"
	
	"sync"
	
	dbm "github.com/tendermint/tm-db"
)

const (
	pathKeySeparator = byte('!')
)

var (
	endBytes = []byte("~~~~")
)

type Path []byte

// ValueOf
func valueOf(keyStr string) Path {
	return Path([]byte(keyStr))
}

func (path Path) String() string {
	return string(path)
}

func (path Path) makeKeyString(key string) []byte {
	return path.makeKey([]byte(key))
}

func (path Path) makeKey(key []byte) []byte {
	newKey := append(path, pathKeySeparator)
	if key != nil {
		newKey = append(newKey, key...)
	}
	return newKey
}

func (path Path) extractKey(key []byte) []byte {
	if bytes.HasPrefix(key, path) {
		return key[len(path)+1:]
	}
	return key
}

type Store struct {
	path Path
	db   dbm.DB
}

var _ types.Store = (*Store)(nil)

func NewStore(path string, db dbm.DB) *Store {
	return &Store{path: valueOf(path), db: db}
}

func (store *Store) GetPath() string {
	return store.path.String()
}

func (store *Store) Get(key []byte) ([]byte, error) {
	keyBytes := store.path.makeKey(key)
	return store.db.Get(keyBytes)
}

func (store *Store) Has(key []byte) (bool, error) {
	keyBytes := store.path.makeKey(key)
	return store.db.Has(keyBytes)
}

func (store *Store) Set(key []byte, value []byte) error {
	keyBytes := store.path.makeKey(key)
	return store.db.Set(keyBytes, value)
}

func (store *Store) SetSync(key []byte, value []byte) error {
	keyBytes := store.path.makeKey(key)
	return store.db.SetSync(keyBytes, value)
}

func (store *Store) Delete(key []byte) error {
	keyBytes := store.path.makeKey(key)
	return store.db.Delete(keyBytes)
}

func (store *Store) DeleteSync(key []byte) error {
	keyBytes := store.path.makeKey(key)
	return store.db.DeleteSync(keyBytes)
}

func (store *Store) Iterator(start, end []byte) (dbm.Iterator, error) {
	startBytes := store.path.makeKey(start)
	var endBts []byte
	if len(end) > 0 {
		endBts = store.path.makeKey(end)
	} else {
		endBts = store.path.makeKey(endBytes)
	}
	
	// fmt.Println("[Store] Iterator ", string(startBytes), string(endBts))
	return store.db.Iterator(startBytes, endBts)
}

func (store *Store) GetMany(start, end []byte) (kvArrayBytes []byte, err error) {
	iterator, err := store.Iterator(start, end)
	// s, e := iterator.Domain()
	// fmt.Println("Store # GetMany: start=", string(s), ", end=", string(e))
	
	if err != nil {
		return nil, err
	}
	
	kvArray := []types.KeyValue{}
	var kv types.KeyValue
	
	for iterator.Valid() {
		key := store.path.extractKey(iterator.Key())
		kv = types.KeyValue{Key: key, Value: iterator.Value()}
		kvArray = append(kvArray, kv)
		
		// fmt.Println("Store # GetMany:", string(iterator.Key()), string(iterator.Value()))
		iterator.Next()
	}
	
	kvArrayBytes, err = types.BasicCdc.MarshalBinaryBare(kvArray)
	
	return kvArrayBytes, err
}

func (store *Store) GetKeys(start, end []byte) (keyArrayBytes []byte, err error) {
	iterator, err := store.Iterator(start, end)
	// fmt.Println("[Store] GetKeys ", start, end)
	
	if err != nil {
		return nil, err
	}
	
	keyArray := []string{}
	
	for iterator.Valid() {
		key := string(store.path.extractKey(iterator.Key()))
		keyArray = append(keyArray, key)
		iterator.Next()
	}
	
	keyArrayBytes, err = types.BasicCdc.MarshalBinaryBare(keyArray)
	
	return keyArrayBytes, err
}

type Registry struct {
	sync.Mutex
	pathStores map[string]*Store
	db         dbm.DB
}

func NewRegistry(db dbm.DB) *Registry {
	return &Registry{pathStores: make(map[string]*Store), db: db}
}

func (reg *Registry) DB() dbm.DB {
	return reg.db
}

func (reg *Registry) RegisterStore(path string) error {
	reg.Lock()
	defer reg.Unlock()
	
	_, ok := reg.pathStores[path]
	if ok {
		return errors.New(fmt.Sprintf("Store [%s] is already registered.", path))
	}
	
	store := &Store{path: valueOf(path), db: reg.db}
	reg.pathStores[path] = store
	return nil
}

func (reg *Registry) GetStore(path string) (*Store, error) {
	reg.Lock()
	defer reg.Unlock()
	
	store, ok := reg.pathStores[path]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Cannot find Store [%s].", path))
	}
	
	return store, nil
}

func (reg *Registry) GetOrMakeStore(path string) *Store {
	reg.Lock()
	defer reg.Unlock()
	
	store, ok := reg.pathStores[path]
	if !ok {
		store = &Store{path: valueOf(path), db: reg.db}
		reg.pathStores[path] = store
	}
	return store
}
