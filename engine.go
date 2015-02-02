package neosearch

import (
	"errors"
	"time"

	"github.com/NeowayLabs/neosearch/store"
)

// StoreCache have the KVStore and LastAccess of the index
type StoreCache struct {
	Store      store.KVStore
	LastAccess time.Time
}

// NGConfig configure the Engine
type NGConfig struct {
	DataDir string
}

// Engine type
type Engine struct {
	stores map[string]StoreCache
	config NGConfig
}

// NewEngine creates a new Engine instance
func NewEngine(config NGConfig) *Engine {
	return &Engine{
		config: config,
	}
}

// Open the index
func (ng *Engine) Open(name string) error {
	ng.stores[name] = StoreCache{
		Store:      store.KVInit(),
		LastAccess: time.Now(),
	}
}

// Execute the given command
func (ng *Engine) Execute(cmd Command) ([]byte, error) {
	var err error

	if ng.stores[cmd.Index] == "" {
		err = ng.Open(cmd.Index)
		if err != nil {
			return nil, err
		}
	}

	storeCache := ng.stores[cmd.Index]
	store := storeCache.Store

	switch cmd.Command {
	case "set":
		err = store.Set(cmd.Key, cmd.Value)
		return nil, error
	case "get":
		return store.Get(cmd.Key)
	}

	return errors.New("Failed to execute command.")
}
