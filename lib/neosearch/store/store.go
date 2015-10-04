package store

import "errors"

// KVReader is a reader safe for concurrent reads.
type KVReader interface {
	Get([]byte) ([]byte, error)
	GetIterator() KVIterator
}

// KVWriter is a writer safe for concurrent writes.
type KVWriter interface {
	Set([]byte, []byte) error
	MergeSet([]byte, uint64) error
	Delete([]byte) error

	StartBatch()
	FlushBatch() error
	IsBatch() bool
}

// KVStore is the key/value store interface for backend kv stores.
type KVStore interface {
	// Open the database
	Open(string, string) error

	// Close the database
	Close() error
	IsOpen() bool

	Reader() KVReader
	NewReader() KVReader
	Writer() KVWriter
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
	Close() error
}

// KVConfig stores the kv configurations
type KVConfig struct {
	Debug       bool
	DataDir     string
	EnableCache bool
	CacheSize   int
}

// KVFuncConstructor is the register function that every store backend need to implement.
type KVFuncConstructor func(*KVConfig) (KVStore, error)

// KVStoreConstructor is a pointer to constructor of default KVStore
var KVStoreConstructor KVFuncConstructor

// KVStoreName have the name of kv store
var KVStoreName string

// SetDefault set the default kv store
func SetDefault(name string, initPtr KVFuncConstructor) error {
	KVStoreName = name
	KVStoreConstructor = initPtr

	return nil
}

// New initialize the default KV store.
func New(config *KVConfig) (KVStore, error) {
	if KVStoreConstructor != nil {
		return KVStoreConstructor(config)
	}

	return nil, errors.New("No store backend configured...")
}
