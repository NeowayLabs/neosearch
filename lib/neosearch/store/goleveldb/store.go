package goleveldb

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/NeowayLabs/neosearch/lib/neosearch/store"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// KVName is the name of goleveldb data store
const KVName = "goleveldb"

type LVDB struct {
	debug   bool
	isBatch bool
	dataDir string

	opts       *opt.Options
	db         *leveldb.DB
	writeBatch *leveldb.Batch

	onceWriter sync.Once
	defWriter  *LVDBWriter
}

func LVDBConstructor(config store.KVConfig) (store.KVStore, error) {
	return NewLVDB(config)
}

func init() {
	store.RegisterKVStore(KVName, LVDBConstructor)
}

func NewLVDB(config store.KVConfig) (*LVDB, error) {
	lvdb := LVDB{
		debug:   false,
		isBatch: false,
	}

	lvdb.setup(config)

	return &lvdb, nil
}

func (lvdb *LVDB) setup(config store.KVConfig) {
	opts := &opt.Options{}

	debug, ok := config["debug"].(bool)
	if ok {
		lvdb.debug = debug
	}

	if debug {
		fmt.Println("Setup goleveldb")
	}

	dataDir, ok := config["dataDir"].(string)
	if ok {
		lvdb.dataDir = dataDir
	} else {
		lvdb.dataDir = "/tmp"
	}

	ro, ok := config["read_only"].(bool)
	if ok {
		opts.ReadOnly = ro
	}

	cim, ok := config["create_if_missing"].(bool)
	if ok {
		opts.ErrorIfMissing = !cim
	}

	eie, ok := config["error_if_exists"].(bool)
	if ok {
		opts.ErrorIfExist = eie
	}

	wbs, ok := config["write_buffer_size"].(float64)
	if ok {
		opts.WriteBuffer = int(wbs)
	}

	bs, ok := config["block_size"].(float64)
	if ok {
		opts.BlockSize = int(bs)
	}

	bri, ok := config["block_restart_interval"].(float64)
	if ok {
		opts.BlockRestartInterval = int(bri)
	}

	lcc, ok := config["lru_cache_capacity"].(float64)
	if ok {
		opts.BlockCacheCapacity = int(lcc)
	}

	bfbpk, ok := config["bloom_filter_bits_per_key"].(float64)
	if ok {
		bf := filter.NewBloomFilter(int(bfbpk))
		opts.Filter = bf
	}

	lvdb.opts = opts
}

func (lvdb *LVDB) Open(indexName, databaseName string) error {
	var err error

	if !store.ValidateDatabaseName(databaseName) {
		return fmt.Errorf("Invalid name: %s", databaseName)
	}

	// index should exists
	fullPath := (lvdb.dataDir + string(filepath.Separator) +
		indexName + string(filepath.Separator) + databaseName)

	lvdb.db, err = leveldb.OpenFile(fullPath, lvdb.opts)
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
	if lvdb.db != nil {
		lvdb.db.Close()
		lvdb.db = nil
	}

	if lvdb.writeBatch != nil {
		lvdb.writeBatch.Reset()
		lvdb.writeBatch = nil
		lvdb.isBatch = false
	}
	return nil
}

// Reader returns a LVDBReader instance
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
