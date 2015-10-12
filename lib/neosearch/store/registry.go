package store

import "fmt"

// KVStoreConstructor is the register function that every store backend need to implement.
type KVStoreConstructor func(*KVConfig) (KVStore, error)

type KVStoreRegistry map[string]KVStoreConstructor

var stores = make(KVStoreRegistry, 0)

func RegisterKVStore(name string, constructor KVStoreConstructor) {
	_, exists := stores[name]
	if exists {
		panic(fmt.Errorf("attempted to register duplicate store named '%s'", name))
	}
	stores[name] = constructor
}

func KVStoreConstructorByName(name string) KVStoreConstructor {
	return stores[name]
}

/*
// New initialize the default KV store.
func New(config *KVConfig) (KVStore, error) {
	if KVStoreConstructor != nil {
		return KVStoreConstructor(config)
	}

	return nil, errors.New("No store backend configured...")
}
*/
