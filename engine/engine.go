package engine

import (
	"errors"
	"fmt"

	"github.com/NeowayLabs/neosearch/cache"
	"github.com/NeowayLabs/neosearch/store"
	"github.com/NeowayLabs/neosearch/utils"
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
	TypeUint   = iota + 1
	TypeInt    = iota
	TypeFloat  = iota
	TypeString = iota
	TypeBinary = iota // TODO: TBD
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
func (ng *Engine) open(name string) (store.KVStore, error) {
	var (
		err     error
		storekv store.KVStore
		ok      bool
		value   interface{}
	)

	value, ok = ng.stores.Get(name)

	if ok == false || value == nil {
		storekv, err = store.New(ng.config.KVCfg)

		if err != nil {
			return nil, err
		}

		ng.stores.Add(name, storekv)
		err = storekv.Open(name)

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

	store, err := ng.GetStore(cmd.Index)

	if ng.config.KVCfg.Debug {
		fmt.Println(cmd)
	}

	if err != nil {
		return nil, err
	}

	switch cmd.Command {
	case "batch":
		store.StartBatch()
		return nil, nil
	case "flushbatch":
		err = store.FlushBatch()
		return nil, err
	case "set":
		err = store.Set(cmd.Key, cmd.Value)
		return nil, err
	case "get":
		return store.Get(cmd.Key)
	case "mergeset":
		v := utils.BytesToUint64(cmd.Value)
		return nil, store.MergeSet(cmd.Key, v)
	case "delete":
		err = store.Delete(cmd.Key)
		return nil, err
	}

	return nil, errors.New("Failed to execute command.")
}

// GetStore returns a instance of KVStore for the given index name
// If the given index name isn't open, then this method will open
// and cache the index for next use.
func (ng *Engine) GetStore(name string) (store.KVStore, error) {
	var (
		err     error
		storekv store.KVStore
	)

	storekv, err = ng.open(name)

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
