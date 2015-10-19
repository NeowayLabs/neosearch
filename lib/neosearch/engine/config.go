package engine

import (
	"github.com/NeowayLabs/neosearch/lib/neosearch/store"
	"github.com/NeowayLabs/neosearch/lib/neosearch/store/goleveldb"
)

const (
	// OpenCacheSize is the default value for the maximum number of
	// open database files. This value can be override by
	// Config.OpenCacheSize.
	DefaultOpenCacheSize int = 100

	// BatchSize is the default size of cached operations before
	// a write batch occurs. You can override this value with
	// Config.BatchSize.
	DefaultBatchSize int = 5000

	// DefaultKVStore set the default KVStore
	DefaultKVStore string = goleveldb.KVName
)

// Config configure the Engine
type Config struct {
	// OpenCacheSize adjust the length of maximum number of
	// open indices. This is a LRU cache, the least used
	// database open will be closed when needed.
	OpenCacheSize int `yaml:"openCacheSize"`

	// BatchSize batch write size
	BatchSize int `yaml:"batchSize"`

	// KVStore configure the kvstore to be used
	KVStore string `yaml:"kvstore"`

	// KVStore specific options to kvstore
	KVConfig store.KVConfig `yaml:"kvconfig"`
}

// New creates new Config
func NewConfig() *Config {
	return &Config{
		DefaultOpenCacheSize,
		DefaultBatchSize,
		DefaultKVStore,
		nil,
	}
}
