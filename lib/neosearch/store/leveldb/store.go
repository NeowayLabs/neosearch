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
func LVDBConstructor(config store.KVConfig) (store.KVStore, error) {
	store, err := NewLVDB(config)
	return store, err
}

// Registry the leveldb module
func init() {
	store.RegisterKVStore(KVName, LVDBConstructor)
}

// LVDB is the leveldb interface exposed by NeoSearch
type LVDB struct {
	debug   bool
	isBatch bool
	dataDir string

	opts         *levigo.Options
	db           *levigo.DB
	readOptions  *levigo.ReadOptions
	writeOptions *levigo.WriteOptions
	writeBatch   *levigo.WriteBatch

	onceWriter sync.Once
	onceReader sync.Once
	defReader  store.KVReader
	defWriter  store.KVWriter
}

// NewLVDB creates a new leveldb instance
func NewLVDB(config store.KVConfig) (*LVDB, error) {
	lvdb := LVDB{
		debug:   false,
		isBatch: false,
	}

	lvdb.setup(config)

	return &lvdb, nil
}

// Setup the leveldb instance
func (lvdb *LVDB) setup(config store.KVConfig) {
	debug, ok := config["debug"].(bool)
	if ok {
		lvdb.debug = debug
	}

	if debug {
		fmt.Println("Setup leveldb")
	}

	dataDir, ok := config["dataDir"].(string)
	if ok {
		lvdb.dataDir = dataDir
	} else {
		lvdb.dataDir = "/tmp"
	}

	lvdb.opts = levigo.NewOptions()

	enableCache, ok := config["enableCache"].(bool)
	if ok && enableCache {
		cacheSize, _ := config["cacheSize"].(int)
		if cacheSize == 0 && enableCache {
			cacheSize = 3 << 30
		}
		lvdb.opts.SetCache(levigo.NewLRUCache(cacheSize))
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
	fullPath := (lvdb.dataDir + string(filepath.Separator) +
		indexName + string(filepath.Separator) + databaseName)

	lvdb.db, err = levigo.Open(fullPath, lvdb.opts)

	if err == nil && lvdb.debug {
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
	lvdb.onceReader.Do(func() {
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
	lvdb.onceWriter.Do(func() {
		lvdb.defWriter = &LVDBWriter{
			store: lvdb,
		}
	})
	return lvdb.defWriter
}
