// +build leveldb

// Package store defines the interface for the KV store technology
package store

import (
	"fmt"

	"github.com/jmhodges/levigo"
)

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

// Registry the leveldb module
func init() {
	initFn := func(config *KVConfig) (*KVStore, error) {
		if config.Debug {
			fmt.Println("Initializing leveldb backend store")
		}

		store, err := NewLVDB(config)
		return &store, err
	}

	err := SetDefault(KVName, &initFn)

	if err != nil {
		fmt.Println("Failed to initialize leveldb backend")
	}
}

// NewLVDB creates a new leveldb instance
func NewLVDB(config *KVConfig) (KVStore, error) {
	lvdb := LVDB{
		Config: config,
	}

	lvdb.Setup()

	return &lvdb, nil
}

// Setup the leveldb instance
func (lvdb *LVDB) Setup() error {
	if lvdb.Config.Debug {
		fmt.Println("Setup leveldb")
	}

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

	// We avoid some cycles by not checking the last '/'
	fullPath := lvdb.Config.DataDir + "/" + path
	lvdb._db, err = levigo.Open(fullPath, lvdb._opts)

	if lvdb.Config.Debug {
		fmt.Printf("Database '%s' opened: %s\n", path, err)
	}

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
	data, err := lvdb.Get(key)

	if err != nil {
		return err
	}

	fmt.Println("Data: ", data)

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

// Close the database
func (lvdb *LVDB) Close() {
	lvdb._db.Close()
}

// LVDBConstructor build the constructor
func LVDBConstructor(config *KVConfig) (*KVStore, error) {
	store, err := NewLVDB(config)
	return &store, err
}
