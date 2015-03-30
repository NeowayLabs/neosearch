package engine

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/NeowayLabs/neosearch/store"
	"github.com/NeowayLabs/neosearch/utils"
)

const (
	// OpenCacheSize is the default value for the maximum number of
	// open database files. This value can be override by
	// NGConfig.OpenCacheSize.
	OpenCacheSize int = 100

	// BatchSize is the default size of cached operations before
	// a write batch occurs. You can override this value with
	// NGConfig.BatchSize.
	BatchSize int = 5000
)

// NGConfig configure the Engine
type NGConfig struct {
	KVCfg *store.KVConfig
	// OpenCacheSize adjust the length of maximum number of
	// open indices. This is a LRU cache, the least used
	// database open will be closed when needed.
	OpenCacheSize int

	// Default batch write size
	BatchSize int
}

// Engine type
type Engine struct {
	stores       StoreCache
	storeEntries []StoreEntry
	config       NGConfig
}

const (
	TypeUint   = iota + 1
	TypeInt    = iota
	TypeFloat  = iota
	TypeString = iota
	TypeBinary = iota // TODO: TBD
)

// Command defines a NeoSearch command
type Command struct {
	Index     string
	Command   string
	Key       []byte
	KeyType   uint8
	Value     []byte
	ValueType uint8

	Batch bool
}

func (c Command) Println() {
	line := c.Reverse()
	fmt.Println(line)
}

func (c Command) Reverse() string {
	var (
		keyStr string
		valStr string
		line   string
	)

	if c.Key != nil {
		if c.KeyType == TypeString {
			keyStr = `'` + string(c.Key) + `'`
		} else if c.KeyType == TypeUint {
			keyStr = `uint(` + strconv.Itoa(int(utils.BytesToUint64(c.Key))) + `)`
		} else if c.KeyType == TypeInt {
			keyStr = `int(` + strconv.Itoa(int(utils.BytesToInt64(c.Key))) + `)`
		} else {
			panic(fmt.Errorf("Invalid command value type: %d", c.ValueType))
		}
	}

	if c.Value != nil {
		if c.ValueType == TypeString {
			valStr = `'` + string(c.Value) + `'`
		} else if c.ValueType == TypeUint {
			valStr = `uint(` + strconv.Itoa(int(utils.BytesToUint64(c.Value))) + `)`
		} else if c.ValueType == TypeInt {
			valStr = `int(` + strconv.Itoa(int(utils.BytesToInt64(c.Value))) + `)`
		} else {
			panic(fmt.Errorf("Invalid command key type: %d", c.KeyType))
		}
	}

	switch c.Command {
	case "set", "mergeset":
		line = fmt.Sprintf("USING %s %s %s %s;", c.Index, c.Command, keyStr, valStr)
	case "batch", "flushbatch":
		line = fmt.Sprintf("USING %s %s;", c.Index, c.Command)
	case "get", "delete":
		line = fmt.Sprintf("USING %s %s %s;", c.Index, c.Command, keyStr)
	default:
		panic(fmt.Errorf("Invalid command: %s: %v", c.Command, c))
	}

	return line
}

// New creates a new Engine instance
// Engine is the generic interface to access database/index files.
// You can execute commands directly to database using Execute method
// acquire direct iterators using the Store interface.
func New(config NGConfig) *Engine {
	ng := &Engine{
		config: config,
		stores: make(StoreCache),
	}

	if ng.config.OpenCacheSize == 0 {
		ng.config.OpenCacheSize = OpenCacheSize
	}

	if ng.config.BatchSize == 0 {
		ng.config.BatchSize = BatchSize
	}

	return ng
}

// cacheClean ensures that only OpenCacheSize indexes are open.
// Closing each of the least accessed of them, until the engine has the
// correct max number of database open (OpenCacheSize config).
func (ng *Engine) cacheClean() {
	var entries = len(ng.storeEntries)
	if entries < ng.config.OpenCacheSize {
		return
	}

	delEntries := ng.storeEntries[ng.config.OpenCacheSize:entries]

	for i := range delEntries {
		entry := delEntries[i]
		store := entry.Store
		store.Close()

		delete(ng.stores, entry.Name)
	}

	ng.storeEntries = ng.storeEntries[0:ng.config.OpenCacheSize]
}

// Open the index and cache then for future uses
func (ng *Engine) open(name string) error {
	storekv, err := store.New(ng.config.KVCfg)

	if err != nil {
		return err
	}

	entry := StoreEntry{
		Store:      storekv,
		LastAccess: time.Now().Nanosecond(),
		Name:       name,
	}

	ng.storeEntries = append(ng.storeEntries, entry)
	sort.Sort(ByLastAccess(ng.storeEntries))

	ng.stores[name] = &entry

	ng.cacheClean()

	err = storekv.Open(name)
	return err
}

// Execute the given command
func (ng *Engine) Execute(cmd Command) ([]byte, error) {
	var err error

	store, err := ng.GetStore(cmd.Index)

	if ng.config.KVCfg.Debug {
		fmt.Println(cmd)
	}

	if err != nil {
		return nil, err
	}

	switch cmd.Command {
	case "batch":
		store.StartBatch()
		return nil, nil
	case "flushbatch":
		err = store.FlushBatch()
		return nil, err
	case "set":
		err = store.Set(cmd.Key, cmd.Value)
		return nil, err
	case "get":
		return store.Get(cmd.Key)
	case "mergeset":
		v := utils.BytesToUint64(cmd.Value)
		return nil, store.MergeSet(cmd.Key, v)
	case "delete":
		err = store.Delete(cmd.Key)
		return nil, err
	}

	return nil, errors.New("Failed to execute command.")
}

// GetStore returns a instance of KVStore for the given index name
// If the given index name isn't open, then this method will open
// and cache the index for next use.
func (ng *Engine) GetStore(name string) (store.KVStore, error) {
	var (
		err   error
		isNew bool
	)

	if ng.stores[name] == nil {
		isNew = true
		err = ng.open(name)
		if err != nil {
			return nil, err
		}
	}

	storeCache := ng.stores[name]

	if !isNew {
		storeCache.LastAccess = time.Now().Nanosecond()
		sort.Sort(ByLastAccess(ng.storeEntries))
	}

	return storeCache.Store, nil
}

// Close all of the open databases
func (ng *Engine) Close() {
	for idx, stc := range ng.stores {
		st := stc.Store
		if st.IsOpen() {
			st.Close()
		}

		ng.stores[idx] = nil
	}
}
