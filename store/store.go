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
}

// KVConfig stores the kv configurations
type KVConfig struct {
}

// KVStoreConstructor is a pointer to constructor of default KVStore
var KVStoreConstructor func(*KVConfig) (*KVStore, error)

// KVStoreName have the name of kv store
var KVStoreName string

// SetDefault set the default kv store
func SetDefault(name string, initPtr func(*KVConfig) (*KVStore, error)) {
	KVStoreName = name
	KVStoreConstructor = initPtr
}

// KVInit initialize the default KV store.
func KVInit() (*KVStore, error) {
	if KVStoreConstructor != nil {
		return KVStoreConstructor(&KVConfig{})
	}

	return nil, errors.New("No store backend configured...")
}
