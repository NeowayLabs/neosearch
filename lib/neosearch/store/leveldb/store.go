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

// LVDB is the leveldb interface exposed by NeoSearch
type LVDB struct {
	debug   bool
	isBatch bool
	dataDir string

	opts       *levigo.Options
	db         *levigo.DB
	writeBatch *levigo.WriteBatch

	onceWriter sync.Once
	defWriter  store.KVWriter
}

// LVDBConstructor build the constructor
func LVDBConstructor(config store.KVConfig) (store.KVStore, error) {
	return NewLVDB(config)
}

// Registry the leveldb module
func init() {
	store.RegisterKVStore(KVName, LVDBConstructor)
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
	if err != nil {
		return err
	}

	if lvdb.debug {
		fmt.Printf("Database '%s' open: %s\n", fullPath, err)
	}
	return nil
}

// IsOpen returns true if database is open
func (lvdb *LVDB) IsOpen() bool {
	return lvdb.db != nil
}

// Close the database
func (lvdb *LVDB) Close() error {
	if lvdb.writeBatch != nil {
		lvdb.writeBatch.Close()
		lvdb.writeBatch = nil
		lvdb.isBatch = false
	}

	if lvdb.db != nil {
		// levigo close implementation does not returns error
		lvdb.db.Close()
		lvdb.db = nil
	}

	if lvdb.opts != nil {
		lvdb.opts.Close()
		lvdb.opts = nil
	}
	return nil
}

// Reader returns a LVDBReader singleton instance
func (lvdb *LVDB) Reader() store.KVReader {
	return newReader(lvdb)
}

// Writer returns the singleton writer
func (lvdb *LVDB) Writer() store.KVWriter {
	lvdb.onceWriter.Do(func() {
		lvdb.defWriter = newWriter(lvdb)
	})
	return lvdb.defWriter
}
