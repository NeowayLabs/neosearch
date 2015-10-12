// +build leveldb

package leveldb

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/NeowayLabs/neosearch/lib/neosearch/store"
	"github.com/jmhodges/levigo"
)

// KVName is the name of leveldb data store
const KVName = "leveldb"

// LVDBConstructor build the constructor
func LVDBConstructor(config *store.KVConfig) (store.KVStore, error) {
	store, err := NewLVDB(config)
	return store, err
}

var (
	onceWriter sync.Once
	onceReader sync.Once
)

// Registry the leveldb module
func init() {
	store.RegisterKVStore(KVName, LVDBConstructor)
}

// LVDB is the leveldb interface exposed by NeoSearch
type LVDB struct {
	Config       *store.KVConfig
	isBatch      bool
	opts         *levigo.Options
	db           *levigo.DB
	readOptions  *levigo.ReadOptions
	writeOptions *levigo.WriteOptions
	writeBatch   *levigo.WriteBatch

	defReader store.KVReader
	defWriter store.KVWriter
}

// NewLVDB creates a new leveldb instance
func NewLVDB(config *store.KVConfig) (*LVDB, error) {
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

	if !store.ValidateDatabaseName(databaseName) {
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
func (lvdb *LVDB) Reader() store.KVReader {
	onceReader.Do(func() {
		lvdb.defReader = lvdb.NewReader()
	})

	return lvdb.defReader
}

// NewReader returns a new reader
func (lvdb *LVDB) NewReader() store.KVReader {
	return &LVDBReader{
		store: lvdb,
	}
}

// Writer returns the singleton writer
func (lvdb *LVDB) Writer() store.KVWriter {
	onceReader.Do(func() {
		lvdb.defWriter = &LVDBWriter{
			store: lvdb,
		}
	})

	return lvdb.defWriter
}
