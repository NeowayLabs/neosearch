package store

// KVReader is a reader safe for concurrent reads.
type KVReader interface {
	Get([]byte) ([]byte, error)
	GetIterator() KVIterator
	Close() error
}

// KVWriter is a writer safe for concurrent writes.
type KVWriter interface {
	Set([]byte, []byte) error
	Get([]byte) ([]byte, error)
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
	IsOpen() bool

	Reader() KVReader
	Writer() KVWriter

	// Close the database
	Close() error
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
type KVConfig map[string]interface{}
