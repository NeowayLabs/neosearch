// +build leveldb

// Package store defines the interface for the KV store technology
package store

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"path/filepath"

	"github.com/NeowayLabs/neosearch/lib/neosearch/utils"
	"github.com/jmhodges/levigo"
)

// KVName is the name of leveldb data store
const KVName = "leveldb"

// LVDBConstructor build the constructor
func LVDBConstructor(config *KVConfig) (KVStore, error) {
	store, err := NewLVDB(config)
	return store, err
}

// Registry the leveldb module
func init() {
	initFn := func(config *KVConfig) (KVStore, error) {
		if config.Debug {
			fmt.Println("Initializing leveldb backend store")
		}

		return NewLVDB(config)
	}

	err := SetDefault(KVName, initFn)

	if err != nil {
		fmt.Println("Failed to initialize leveldb backend")
	}
}

// LVDB is the leveldb interface exposed by NeoSearch
type LVDB struct {
	Config        *KVConfig
	isBatch       bool
	_opts         *levigo.Options
	_db           *levigo.DB
	_readOptions  *levigo.ReadOptions
	_writeOptions *levigo.WriteOptions
	_writeBatch   *levigo.WriteBatch
}

// NewLVDB creates a new leveldb instance
func NewLVDB(config *KVConfig) (*LVDB, error) {
	lvdb := LVDB{
		Config: config,
	}

	lvdb.setup()

	return &lvdb, nil
}

// Setup the leveldb instance
func (lvdb *LVDB) setup() {
	if lvdb.Config.Debug {
		fmt.Println("Setup leveldb")
	}

	lvdb._opts = levigo.NewOptions()

	if lvdb.Config.EnableCache {
		lvdb._opts.SetCache(levigo.NewLRUCache(lvdb.Config.CacheSize))
	}

	lvdb._opts.SetCreateIfMissing(true)

	// TODO: export this configuration options
	lvdb._readOptions = levigo.NewReadOptions()
	lvdb._writeOptions = levigo.NewWriteOptions()
}

// Open the database
func (lvdb *LVDB) Open(indexName, databaseName string) error {
	var err error

	if !validateDatabaseName(databaseName) {
		return fmt.Errorf("Invalid name: %s", databaseName)
	}

	// index should exists
	fullPath := (lvdb.Config.DataDir + string(filepath.Separator) +
		indexName + string(filepath.Separator) + databaseName)

	lvdb._db, err = levigo.Open(fullPath, lvdb._opts)

	if err == nil && lvdb.Config.Debug {
		fmt.Printf("Database '%s' open: %s\n", fullPath, err)
	}

	return err
}

// IsOpen returns true if database is open
func (lvdb *LVDB) IsOpen() bool {
	return lvdb._db != nil
}

// Set put or update the key with the given value
func (lvdb *LVDB) Set(key []byte, value []byte) error {
	if lvdb.isBatch {
		// isBatch == true, we can safely access _writeBatch pointer
		lvdb._writeBatch.Put(key, value)
		return nil
	}

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
		buf      *bytes.Buffer
		err      error
		v        uint64
		i        uint64
		inserted bool
	)

	data, err := lvdb.Get(key)

	if err != nil {
		return err
	}

	if lvdb.Config.Debug {
		fmt.Printf("[INFO] %d ids == %d GB of ids\n", len(data)/8, len(data)/(1000*1000*1000))
        }                

	buf = new(bytes.Buffer)
	lenBytes := uint64(len(data))

	// O(n)
	for i = 0; i < lenBytes; i += 8 {
		v = utils.BytesToUint64(data[i : i+8])

		// returns if value is already stored
		if v == value {
			return nil
		}

		if value < v {
			err = binary.Write(buf, binary.BigEndian, value)
			if err != nil {
				return err
			}

			inserted = true
		}

		err = binary.Write(buf, binary.BigEndian, v)
		if err != nil {
			return err
		}
	}

	if lenBytes == 0 || !inserted {
		err = binary.Write(buf, binary.BigEndian, value)
		if err != nil {
			return err
		}
	}

	return lvdb.Set(key, buf.Bytes())
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
	if lvdb.isBatch {
		lvdb._writeBatch.Delete(key)
		return nil
	}

	return lvdb._db.Delete(lvdb._writeOptions, key)
}

// DeleteCustom is the same as Delete but enables override default write options
func (lvdb *LVDB) DeleteCustom(opt *levigo.WriteOptions, key []byte) error {
	return lvdb._db.Delete(opt, key)
}

// StartBatch start a new batch write processing
func (lvdb *LVDB) StartBatch() {
	if lvdb._writeBatch == nil {
		lvdb._writeBatch = levigo.NewWriteBatch()
	} else {
		lvdb._writeBatch.Clear()
	}

	lvdb.isBatch = true
}

// IsBatch returns true if LVDB is in batch mode
func (lvdb *LVDB) IsBatch() bool {
	return lvdb.isBatch
}

// FlushBatch writes the batch to disk
func (lvdb *LVDB) FlushBatch() error {
	var err error
	if lvdb._writeBatch != nil {
		err = lvdb._db.Write(lvdb._writeOptions, lvdb._writeBatch)
		// After flush, release the writeBatch for future uses
		lvdb._writeBatch.Clear()
		lvdb.isBatch = false
	}

	return err
}

// Close the database
func (lvdb *LVDB) Close() {
	if lvdb._db != nil {
		lvdb._db.Close()
		lvdb._db = nil
	}

	if lvdb._writeBatch != nil {
		lvdb._writeBatch.Close()
		lvdb._writeBatch = nil
		lvdb.isBatch = false
	}
}

// GetIterator returns a new KVIterator
func (lvdb *LVDB) GetIterator() KVIterator {
	var ro = lvdb._readOptions

	ro.SetFillCache(false)
	it := lvdb._db.NewIterator(ro)
	return it
}
