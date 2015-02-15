package engine

import (
	"errors"
	"strconv"
	"time"

	"bitbucket.org/i4k/neosearch/store"
)

// StoreCache have the KVStore and LastAccess of the index
type StoreCache struct {
	Store      *store.KVStore
	LastAccess time.Time
}

// NGConfig configure the Engine
type NGConfig struct {
	KVCfg *store.KVConfig
}

// Engine type
type Engine struct {
	stores map[string]StoreCache
	config NGConfig
}

// Command defines a NeoSearch command
type Command struct {
	Index   string
	Command string
	Key     []byte
	Value   []byte
}

// New creates a new Engine instance
func New(config NGConfig) *Engine {
	return &Engine{
		config: config,
		stores: make(map[string]StoreCache),
	}
}

// Open the index
func (ng *Engine) Open(name string) error {
	storekv, err := store.KVInit(ng.config.KVCfg)

	if err != nil {
		return err
	}

	ng.stores[name] = StoreCache{
		Store:      storekv,
		LastAccess: time.Now(),
	}

	err = (*storekv).Open(name)
	return err
}

// Execute the given command
func (ng *Engine) Execute(cmd Command) ([]byte, error) {
	var err error

	if ng.stores[cmd.Index].Store == nil {
		err = ng.Open(cmd.Index)
		if err != nil {
			return nil, err
		}
	}

	storeCache := ng.stores[cmd.Index]
	store := storeCache.Store

	switch cmd.Command {
	case "set":
		err = (*store).Set(cmd.Key, cmd.Value)
		return nil, err
	case "get":
		return (*store).Get(cmd.Key)
	case "mergeset":
		v, _ := strconv.ParseInt(string(cmd.Value), 10, 64)
		return nil, (*store).MergeSet(cmd.Key, uint64(v))
	case "delete":
		err = (*store).Delete(cmd.Key)
		return nil, err
	}

	return nil, errors.New("Failed to execute command.")
}

// Close all of the opened databases
func (ng *Engine) Close() {
	for _, stc := range ng.stores {
		st := stc.Store
		(*st).Close()
	}
}
