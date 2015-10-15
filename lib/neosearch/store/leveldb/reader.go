// +build leveldb

package leveldb

import (
	"github.com/NeowayLabs/neosearch/lib/neosearch/store"
	"github.com/jmhodges/levigo"
)

// LVDBReader is the readonly view of database. It implements the KVReader interface
type LVDBReader struct {
	store *LVDB
}

// newReader returns a new reader
func newReader(lvdb *LVDB) *LVDBReader {
	return &LVDBReader{
		store: lvdb,
	}
}

// Get returns the value of the given key
func (reader *LVDBReader) Get(key []byte) ([]byte, error) {
	return reader.store.db.Get(reader.store.readOptions, key)
}

// GetCustom is the same as Get but enables override default read options
func (reader *LVDBReader) GetCustom(opt *levigo.ReadOptions, key []byte) ([]byte, error) {
	return reader.store.db.Get(opt, key)
}

// GetIterator returns a new KVIterator
func (reader *LVDBReader) GetIterator() store.KVIterator {
	var ro = reader.store.readOptions

	ro.SetFillCache(false)
	it := reader.store.db.NewIterator(ro)
	return LVDBIterator{it}
}

func (reader *LVDBReader) Close() error {
	return nil
}
