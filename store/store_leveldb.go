// +build leveldb

// Package store defines the interface for the KV store technology
package store

import "github.com/jmhodges/levigo"

// KVName is the name of leveldb data store
const KVName = "leveldb"

// LVDB is the leveldb interface exposed by NeoSearch
type LVDB struct {
	Config        *KVConfig
	_opts         *levigo.Options
	_db           *levigo.DB
	_readOptions  *levigo.ReadOptions
	_writeOptions *levigo.WriteOptions
}

// NewLVDB creates a new leveldb instance
func NewLVDB(config *KVConfig) (KVStore, error) {
	return &LVDB{
		Config: config,
	}, nil
}

// Setup the leveldb instance
func (lvdb *LVDB) Setup() error {
	lvdb._opts = levigo.NewOptions()
	lvdb._opts.SetCache(levigo.NewLRUCache(3 << 30))
	lvdb._opts.SetCreateIfMissing(true)

	lvdb._readOptions = levigo.NewReadOptions()
	lvdb._writeOptions = levigo.NewWriteOptions()

	return nil
}

// Open the database
func (lvdb *LVDB) Open(path string) error {
	var err error
	lvdb._db, err = levigo.Open(path, lvdb._opts)

	return err
}

// Set put or update the key with the given value
func (lvdb *LVDB) Set(key []byte, value []byte) error {
	return lvdb._db.Put(lvdb._writeOptions, key, value)
}

// SetCustom is the same as Set but enables override default write options
func (lvdb *LVDB) SetCustom(opt *levigo.WriteOptions, key []byte, value []byte) error {
	return lvdb._db.Put(opt, key, value)
}

// MergeSet isn't implemented yet
func (lvdb *LVDB) MergeSet(key []byte, value uint64) error {
	return nil
}

// Get returns the value of the given key
func (lvdb *LVDB) Get(key []byte) ([]byte, error) {
	return lvdb._db.Get(lvdb._readOptions, key)
}

// GetCustom is the same as Get but enables override default read options
func (lvdb *LVDB) GetCustom(opt *levigo.ReadOptions, key []byte) ([]byte, error) {
	return lvdb._db.Get(opt, key)
}

// Delete remove the given key
func (lvdb *LVDB) Delete(key []byte) error {
	return lvdb._db.Delete(lvdb._writeOptions, key)
}

// DeleteCustom is the same as Delete but enables override default write options
func (lvdb *LVDB) DeleteCustom(opt *levigo.WriteOptions, key []byte) error {
	return lvdb._db.Delete(opt, key)
}

// LVDBConstructor build the constructor
func LVDBConstructor(config *KVConfig) (*KVStore, error) {
	store, err := NewLVDB(config)
	return &store, err
}

func init() {
	//	error := SetDefault(KVName, func(config *KVConfig) (*KVStore, error) {
	//		store, err := NewLVDB(config)
	//		return &store, err
	//	})
}
