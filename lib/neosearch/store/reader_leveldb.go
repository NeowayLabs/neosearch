// +build leveldb

package store

import "github.com/jmhodges/levigo"

// LVDBReader is the readonly view of database. It implements the KVReader interface
type LVDBReader struct {
	store *LVDB
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
func (reader *LVDBReader) GetIterator() KVIterator {
	var ro = reader.store.readOptions

	ro.SetFillCache(false)
	it := reader.store.db.NewIterator(ro)
	wrapper := LVDBIterator{it}
	return wrapper
}
