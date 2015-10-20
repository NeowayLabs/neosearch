// +build leveldb

package leveldb

import (
	"github.com/NeowayLabs/neosearch/lib/neosearch/store"
	"github.com/jmhodges/levigo"
)

// LVDBReader is the readonly view of database. It implements the KVReader interface
type LVDBReader struct {
	store    *LVDB
	snapshot *levigo.Snapshot
}

// newReader returns a new reader
func newReader(store *LVDB) store.KVReader {
	return &LVDBReader{
		store:    store,
		snapshot: store.db.NewSnapshot(),
	}
}

// Get returns the value of the given key
func (reader *LVDBReader) Get(key []byte) ([]byte, error) {
	options := defaultReadOptions()
	options.SetSnapshot(reader.snapshot)
	b, err := reader.store.db.Get(options, key)
	options.Close()
	return b, err
}

// GetIterator returns a new KVIterator
func (reader *LVDBReader) GetIterator() store.KVIterator {
	options := defaultReadOptions()
	options.SetSnapshot(reader.snapshot)
	it := reader.store.db.NewIterator(options)
	options.Close()
	return &LVDBIterator{it}
}

func (reader *LVDBReader) Close() error {
	reader.store.db.ReleaseSnapshot(reader.snapshot)
	return nil
}
