package goleveldb

import (
	"github.com/NeowayLabs/neosearch/lib/neosearch/store"
	"github.com/syndtr/goleveldb/leveldb"
)

type LVDBReader struct {
	store    *LVDB
	snapshot *leveldb.Snapshot
}

func newReader(store *LVDB) *LVDBReader {
	snapshot, _ := store.db.GetSnapshot()
	return &LVDBReader{
		store:    store,
		snapshot: snapshot,
	}
}

// Get returns the value of the given key
func (reader *LVDBReader) Get(key []byte) ([]byte, error) {
	options := defaultReadOptions()
	b, err := reader.snapshot.Get(key, options)
	if err == leveldb.ErrNotFound {
		return nil, nil
	}
	return b, err
}

func (reader *LVDBReader) GetIterator() store.KVIterator {
	return newIteratorWithSnapshot(reader.store, reader.snapshot)
}

func (reader *LVDBReader) Close() error {
	reader.snapshot.Release()
	return nil
}
