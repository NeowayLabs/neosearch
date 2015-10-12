package store

import "errors"

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
