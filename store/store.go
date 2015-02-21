package store

import "errors"

// KVStore is the key/value store interface for other backend kv stores.
type KVStore interface {
	// Open the database
	Open(string) error
	Get([]byte) ([]byte, error)
	Set([]byte, []byte) error
	MergeSet([]byte, uint64) error
	Delete([]byte) error
	Close()

	GetIterator() KVIterator
	StartBatch()
	FlushBatch() error

	IsBatch() bool
	IsOpen() bool
}

// KVIterator expose the interface for database iterators.
// This was Based on leveldb interface
type KVIterator interface {
	Valid() bool
	Key() []byte
	Value() []byte
	Next()
	Prev()
	SeekToFirst()
	SeekToLast()
	Seek([]byte)
	GetError() error
	Close()
}

// KVConfig stores the kv configurations
type KVConfig struct {
	Debug       bool
	DataDir     string
	EnableCache bool
	CacheSize   int
}

// KVStoreConstructor is a pointer to constructor of default KVStore
var KVStoreConstructor *func(*KVConfig) (KVStore, error)

// KVStoreName have the name of kv store
var KVStoreName string

// SetDefault set the default kv store
func SetDefault(name string, initPtr *func(*KVConfig) (KVStore, error)) error {
	KVStoreName = name
	KVStoreConstructor = initPtr

	return nil
}

// New initialize the default KV store.
func New(config *KVConfig) (KVStore, error) {
	if KVStoreConstructor != nil {
		return (*KVStoreConstructor)(config)
	}

	return nil, errors.New("No store backend configured...")
}
