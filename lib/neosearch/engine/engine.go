package engine

import (
	"errors"

	"github.com/NeowayLabs/neosearch/lib/neosearch/cache"
	"github.com/NeowayLabs/neosearch/lib/neosearch/store"
	"github.com/NeowayLabs/neosearch/lib/neosearch/utils"
)

// Engine type
type Engine struct {
	stores cache.Cache
	config *Config
	debug  bool
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
func New(cfg *Config) *Engine {
	if cfg.OpenCacheSize == 0 {
		cfg.OpenCacheSize = DefaultOpenCacheSize
	}

	if cfg.BatchSize == 0 {
		cfg.BatchSize = DefaultBatchSize
	}

	if cfg.KVStore == "" {
		cfg.KVStore = DefaultKVStore
	}

	if cfg.KVConfig == nil {
		cfg.KVConfig = store.KVConfig{}
	}

	ng := &Engine{
		config: cfg,
		stores: cache.NewLRUCache(cfg.OpenCacheSize),
	}

	if debug, ok := cfg.KVConfig["debug"].(bool); ok {
		ng.debug = debug
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
		storeConstructor := store.KVStoreConstructorByName(ng.config.KVStore)
		if storeConstructor == nil {
			return nil, errors.New("Unknown storage type")
		}

		storekv, err = storeConstructor(ng.config.KVConfig)
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

	if ng.debug {
		cmd.Println()
	}

	if err != nil {
		return nil, err
	}

	writer := store.Writer()

	reader := store.Reader()
	defer func() {
		reader.Close()
	}()

	switch cmd.Command {
	case "batch":
		writer.StartBatch()
		return nil, nil
	case "flushbatch":
		err = writer.FlushBatch()
		return nil, err
	case "set":
		err = writer.Set(cmd.Key, cmd.Value)
		return nil, err
	case "get":
		return reader.Get(cmd.Key)
	case "mergeset":
		v := utils.BytesToUint64(cmd.Value)
		return nil, writer.MergeSet(cmd.Key, v)
	case "delete":
		err = writer.Delete(cmd.Key)
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
