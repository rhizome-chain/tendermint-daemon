package store

import (
	"testing"
	
	dbm "github.com/tendermint/tm-db"
)

func TestStore_DeleteByPrefix(t *testing.T) {
	db, err := dbm.NewGoLevelDB("test", "./")
	if err != nil {
		t.Error(err)
	}
	
	store := &Store{path: valueOf("test"), db: db}
	
	store.Set([]byte("a1/aaa"), []byte("aaa"))
	store.Set([]byte("a1/bbb"), []byte("bbb"))
	store.Set([]byte("a1/ccc"), []byte("ccc"))
	store.Set([]byte("a2/aaa"), []byte("aaa"))
	
	err = store.DeleteByPrefix([]byte("a1"))
	if err != nil {
		t.Error(err)
	}
	
	val, _ := store.Get([]byte("a2/aaa"))
	if "aaa" != string(val)  {
		t.Error("Irrelevant item Deleted")
	}
	
	val, _ = store.Get([]byte("a1/ccc"))
	if val != nil  {
		t.Error("Item is not Deleted")
	}
}
