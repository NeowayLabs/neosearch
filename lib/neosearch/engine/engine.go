package engine

import (
	"errors"

	"github.com/NeowayLabs/neosearch/lib/neosearch/cache"
	"github.com/NeowayLabs/neosearch/lib/neosearch/store"
	"github.com/NeowayLabs/neosearch/lib/neosearch/utils"
)

const (
	// OpenCacheSize is the default value for the maximum number of
	// open database files. This value can be override by
	// NGConfig.OpenCacheSize.
	OpenCacheSize int = 100

	// BatchSize is the default size of cached operations before
	// a write batch occurs. You can override this value with
	// NGConfig.BatchSize.
	BatchSize int = 5000
)

// NGConfig configure the Engine
type NGConfig struct {
	KVCfg *store.KVConfig
	// OpenCacheSize adjust the length of maximum number of
	// open indices. This is a LRU cache, the least used
	// database open will be closed when needed.
	OpenCacheSize int

	// Default batch write size
	BatchSize int
}

// Engine type
type Engine struct {
	stores cache.Cache
	config NGConfig
}

const (
	_ = iota
	TypeNil
	TypeUint
	TypeInt
	TypeFloat
	TypeString
	TypeDate
	TypeBool
	TypeBinary // TODO: TBD
)

// New creates a new Engine instance
// Engine is the generic interface to access database/index files.
// You can execute commands directly to database using Execute method
// acquire direct iterators using the Store interface.
func New(config NGConfig) *Engine {
	if config.OpenCacheSize == 0 {
		config.OpenCacheSize = OpenCacheSize
	}

	if config.BatchSize == 0 {
		config.BatchSize = BatchSize
	}

	ng := &Engine{
		config: config,
		stores: cache.NewLRUCache(config.OpenCacheSize),
	}

	ng.stores.OnRemove(func(key string, value interface{}) {
		storekv, ok := value.(store.KVStore)

		if !ok {
			panic("Unexpected value in cache")
		}

		if storekv.IsOpen() {
			storekv.Close()
		}
	})

	return ng
}

// Open the index and cache then for future uses
func (ng *Engine) open(indexName, databaseName string) (store.KVStore, error) {
	var (
		err     error
		storekv store.KVStore
		ok      bool
		value   interface{}
	)

	value, ok = ng.stores.Get(indexName + "." + databaseName)

	if ok == false || value == nil {
		storeConstructor := store.KVStoreConstructorByName("leveldb")
		if storeConstructor == nil {
			return nil, errors.New("Unknown storage type")
		}

		storekv, err = storeConstructor(ng.config.KVCfg)
		if err != nil {
			return nil, err
		}

		ng.stores.Add(indexName+"."+databaseName, storekv)
		err = storekv.Open(indexName, databaseName)

		return storekv, err
	}

	storekv, ok = value.(store.KVStore)

	if ok {
		return storekv, nil
	}

	return nil, errors.New("Failed to convert cache entry to KVStore")
}

// Execute the given command
func (ng *Engine) Execute(cmd Command) ([]byte, error) {
	var err error

	store, err := ng.GetStore(cmd.Index, cmd.Database)

	if ng.config.KVCfg.Debug {
		cmd.Println()
	}

	if err != nil {
		return nil, err
	}

	switch cmd.Command {
	case "batch":
		store.Writer().StartBatch()
		return nil, nil
	case "flushbatch":
		err = store.Writer().FlushBatch()
		return nil, err
	case "set":
		err = store.Writer().Set(cmd.Key, cmd.Value)
		return nil, err
	case "get":
		return store.Reader().Get(cmd.Key)
	case "mergeset":
		v := utils.BytesToUint64(cmd.Value)
		return nil, store.Writer().MergeSet(cmd.Key, v)
	case "delete":
		err = store.Writer().Delete(cmd.Key)
		return nil, err
	}

	return nil, errors.New("Failed to execute command.")
}

// GetStore returns a instance of KVStore for the given index name
// If the given index name isn't open, then this method will open
// and cache the index for next use.
func (ng *Engine) GetStore(indexName string, databaseName string) (store.KVStore, error) {
	var (
		err     error
		storekv store.KVStore
	)

	storekv, err = ng.open(indexName, databaseName)

	if err != nil {
		return nil, err
	}

	return storekv, nil
}

// Close all of the open databases
func (ng *Engine) Close() {
	// Clean will un-ref and Close the databases
	ng.stores.Clean()
}
