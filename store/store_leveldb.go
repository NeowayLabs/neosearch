// +build leveldb

// Package store defines the interface for the KV store technology
package store

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"sort"

	"bitbucket.org/i4k/neosearch/utils"

	"github.com/jmhodges/levigo"
)

// KVName is the name of leveldb data store
const KVName = "leveldb"

// LVDBConstructor build the constructor
func LVDBConstructor(config *KVConfig) (*KVStore, error) {
	store, err := NewLVDB(config)
	return &store, err
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

	if lvdb.Config.EnableCache {
		lvdb._opts.SetCache(levigo.NewLRUCache(lvdb.Config.CacheSize))
	}

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

// MergeSet add value to a ordered set of integers stored in key. If value
// is already on the key, than the set will be skipped.
func (lvdb *LVDB) MergeSet(key []byte, value uint64) error {
	var (
		valueBytes bytes.Buffer
		values     []uint64
		err        error
	)

	dataGob, err := lvdb.Get(key)

	if err != nil {
		return err
	}

	if len(dataGob) > 0 {
		valueBytes.Write(dataGob)
		decoder := gob.NewDecoder(&valueBytes)

		err = decoder.Decode(&values)

		if err != nil {
			return err
		}
	}

	exists := false
	// only testing, we need a efficient solution
	for i := 0; i < len(values); i++ {
		if values[i] == value {
			exists = true
			break
		}
	}

	if !exists {
		values = append(values, value)
		sort.Sort(utils.Uint64Slice(values))

		encoder := gob.NewEncoder(&valueBytes)
		if err = encoder.Encode(values); err == nil {
			return lvdb.Set(key, valueBytes.Bytes())
		}

		return err
	}

	// value already in the set
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

// GetIterator returns a new KVIterator
func (lvdb *LVDB) GetIterator() KVIterator {
	var ro = lvdb._readOptions

	ro.SetFillCache(false)
	it := lvdb._db.NewIterator(ro)
	return it
}
