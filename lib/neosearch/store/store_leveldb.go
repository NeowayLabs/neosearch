// +build leveldb

// Package store defines the interface for the KV store technology
package store

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/jmhodges/levigo"
)

// KVName is the name of leveldb data store
const KVName = "leveldb"

// LVDBConstructor build the constructor
func LVDBConstructor(config *KVConfig) (KVStore, error) {
	store, err := NewLVDB(config)
	return store, err
}

var (
	muGetWriter *sync.Mutex
	muGetReader *sync.Mutex
)

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

	muGetWriter = &sync.Mutex{}
	muGetReader = &sync.Mutex{}
}

// LVDB is the leveldb interface exposed by NeoSearch
type LVDB struct {
	Config       *KVConfig
	isBatch      bool
	opts         *levigo.Options
	db           *levigo.DB
	readOptions  *levigo.ReadOptions
	writeOptions *levigo.WriteOptions
	writeBatch   *levigo.WriteBatch

	defReader KVReader
	defWriter KVWriter
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

	lvdb.opts = levigo.NewOptions()

	if lvdb.Config.EnableCache {
		lvdb.opts.SetCache(levigo.NewLRUCache(lvdb.Config.CacheSize))
	}

	lvdb.opts.SetCreateIfMissing(true)

	// TODO: export this configuration options
	lvdb.readOptions = levigo.NewReadOptions()
	lvdb.writeOptions = levigo.NewWriteOptions()
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

	lvdb.db, err = levigo.Open(fullPath, lvdb.opts)

	if err == nil && lvdb.Config.Debug {
		fmt.Printf("Database '%s' open: %s\n", fullPath, err)
	}

	return err
}

// IsOpen returns true if database is open
func (lvdb *LVDB) IsOpen() bool {
	return lvdb.db != nil
}

// Close the database
func (lvdb *LVDB) Close() error {
	if lvdb.db != nil {
		// levigo close implementation does not returns error
		lvdb.db.Close()
		lvdb.db = nil
	}

	if lvdb.writeBatch != nil {
		lvdb.writeBatch.Close()
		lvdb.writeBatch = nil
		lvdb.isBatch = false
	}

	return nil
}

// Reader returns a LVDBReader singleton instance
func (lvdb *LVDB) Reader() KVReader {
	muGetReader.Lock()
	defer muGetReader.Unlock()

	if lvdb.defReader == nil {
		lvdb.defReader = lvdb.NewReader()
	}

	return lvdb.defReader
}

// NewReader returns a new reader
func (lvdb *LVDB) NewReader() KVReader {
	return &LVDBReader{
		store: lvdb,
	}
}

// Writer returns the singleton writer
func (lvdb *LVDB) Writer() KVWriter {
	muGetWriter.Lock()
	defer muGetWriter.Unlock()

	if lvdb.defWriter == nil {
		lvdb.defWriter = &LVDBWriter{
			store: lvdb,
		}
	}

	return lvdb.defWriter
}
