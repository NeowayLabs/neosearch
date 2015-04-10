package index

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/NeowayLabs/neosearch/engine"
	"github.com/NeowayLabs/neosearch/store"
	"github.com/NeowayLabs/neosearch/utils"
)

const (
	dbName   string = "document.db"
	indexExt string = "idx"

	infoFilename string = "info.json"
)

// Config index
type Config struct {
	// Per-Index configurations

	DataDir     string
	Debug       bool
	CacheSize   int
	EnableCache bool
}

// Index represents an entire index
type Index struct {
	Name string `json:"name"`

	engine *engine.Engine
	config Config

	// Indicates that index abstraction should batch each write command
	enableBatchMode bool

	flushStorages []string

	info      *IndexInfo
	infoMutex *sync.Mutex
}

// ValidateIndexName verifies if name is valid NeoSearch index name
func ValidateIndexName(name string) bool {
	if len(name) < 3 {
		return false
	}

	validName := regexp.MustCompile(`^[a-zA-Z0-9-]+$`)
	return validName.MatchString(name)
}

// New creates new index
func New(name string, cfg Config, create bool) (*Index, error) {
	if !ValidateIndexName(name) {
		return nil, errors.New("Invalid index name")
	}

	index := &Index{
		Name:      name,
		config:    cfg,
		info:      NewIndexInfo(),
		infoMutex: &sync.Mutex{},
	}

	if err := index.setup(create); err != nil {
		return nil, err
	}

	return index, nil
}

func (i *Index) setup(create bool) error {
	dataDir := i.config.DataDir + "/" + i.Name
	if create {
		if err := os.Mkdir(dataDir, 0755); err != nil {
			return err
		}
	}

	// index dataDir = <neosearch datadir> + "/" + <index name>
	i.config.DataDir = dataDir

	i.engine = engine.New(engine.NGConfig{
		KVCfg: &store.KVConfig{
			DataDir:     i.config.DataDir,
			Debug:       i.config.Debug,
			CacheSize:   i.config.CacheSize,
			EnableCache: i.config.EnableCache,
		},
	})

	return i.createInfoFile()
}

func (i *Index) createInfoFile() error {
	statsFile := i.config.DataDir + "/" + infoFilename
	jsonContent, err := json.Marshal(i.info)

	if err != nil {
		return err
	}

	file, err := os.Create(statsFile)

	if err != nil {
		return err
	}

	_, err = file.Write(jsonContent)
	return err
}

func (i *Index) updateInfo(indexCounter map[string]uint64) {
	var (
		field FieldInfo
		ok    bool
	)

	for index, counter := range indexCounter {
		field, ok = i.info.Fields[index]

		if ok {
			field.Size += counter
		} else {
			field = FieldInfo{
				Size: counter,
			}
		}

		i.info.Fields[index] = field
	}

	i.infoMutex.Lock()
	os.Remove(i.config.DataDir + "/" + infoFilename)
	i.createInfoFile()
	i.infoMutex.Unlock()
}

// Batch enables write cache of command before FlushBatch is executed
func (i *Index) Batch() {
	i.enableBatchMode = true
}

// FlushBatch writes the cached commands to disk
// Simple and ugly approach. Only to test the concepts.
func (i *Index) FlushBatch() {
	// flush the WriteBatch
	for _, storeName := range i.flushStorages {
		i.engine.Execute(engine.Command{
			Index:   storeName,
			Command: "flushbatch",
		})

		if i.config.Debug {
			fmt.Printf("Flushing batch storage '%s' of index '%s'.\n",
				storeName,
				i.Name)
		}
	}

	i.flushStorages = make([]string, 0)
}

// Add executes the sequence of commands necessary to index the document `doc`.
func (i *Index) Add(id uint64, doc []byte) error {
	indicesCounter := make(map[string]uint64)

	commands, err := i.BuildAdd(id, doc)

	if err != nil {
		return err
	}

	for _, cmd := range commands {
		_, err := i.engine.Execute(cmd)

		if err != nil {
			return err
		}

		val := indicesCounter[cmd.Index] + 1
		indicesCounter[cmd.Index] = val
	}

	go i.updateInfo(indicesCounter)

	return nil
}

func (i *Index) BuildAdd(id uint64, doc []byte) ([]engine.Command, error) {
	var commands []engine.Command

	if i.enableBatchMode {
		// batchMode says if the BATCH operation is pending on indices
		// If true, then we need run USING <idx> BATCH; on each index.
		// When the Add method returns, this means that every index is
		// propertly in batchMode, then we can disable here.
		defer func() {
			i.enableBatchMode = false
		}()
	}

	docCommands, err := i.buildAddDocument(id, doc)

	if err != nil {
		return nil, err
	}

	structData := map[string]interface{}{}

	err = json.Unmarshal(doc, &structData)

	if err != nil {
		return nil, err
	}

	fieldCommands, err := i.buildIndexFields(id, &structData)

	if err != nil {
		return nil, err
	}

	for _, cmd := range docCommands {
		commands = append(commands, cmd)
	}

	for _, cmd := range fieldCommands {
		commands = append(commands, cmd)
	}

	return commands, nil
}

func (i *Index) buildAddDocument(id uint64, doc []byte) ([]engine.Command, error) {
	var commands []engine.Command

	if i.enableBatchMode {
		cmd, err := i.buildBatchOn(dbName)

		if err == nil {
			commands = append(commands, cmd)
		}
	}

	cmd := engine.Command{}
	cmd.Index = dbName
	cmd.Key = utils.Uint64ToBytes(id)
	cmd.KeyType = engine.TypeUint
	cmd.Value = doc
	cmd.ValueType = engine.TypeString
	cmd.Command = "set"

	commands = append(commands, cmd)
	return commands, nil
}

func (i *Index) buildGet(id uint64) engine.Command {
	return engine.Command{
		Index:   "document.db",
		Command: "get",
		Key:     utils.Uint64ToBytes(id),
		KeyType: engine.TypeUint,
	}
}

// Get retrieves the document by id
func (i *Index) Get(id uint64) ([]byte, error) {
	return i.engine.Execute(i.buildGet(id))
}

func (i *Index) GetAnalyze(id uint64) (engine.Command, error) {
	return i.buildGet(id), nil
}

func (i *Index) buildBatchOn(storage string) (engine.Command, error) {
	if i.config.Debug {
		fmt.Printf("Batch mode enabled for storage '%s' of index '%s'.\n",
			storage,
			i.Name,
		)
	}

	for _, flush := range i.flushStorages {
		if flush == storage {
			return engine.Command{}, errors.New("Index already in batch mode")
		}
	}

	i.flushStorages = append(i.flushStorages, storage)

	command := engine.Command{
		Index:   storage,
		Command: "batch",
	}

	return command, nil
}

func (i *Index) enableBatchOn(storage string) error {
	cmd, err := i.buildBatchOn(storage)

	if err != nil {
		return nil
	}

	_, err = i.engine.Execute(cmd)

	return err
}

// buildIndexFields builds the list of commands to index document fields. Note that
// the order os commands generated by field is sorted lexicografically (sort.Strings)
func (i *Index) buildIndexFields(id uint64, structData *map[string]interface{}) ([]engine.Command, error) {
	var (
		commands []engine.Command
		dataKeys []string
	)

	for key := range *structData {
		dataKeys = append(dataKeys, key)
	}

	sort.Strings(dataKeys)

	for _, key := range dataKeys {
		value := (*structData)[key]
		cmds, err := i.buildIndexField(id, []byte(key), value)

		if err != nil {
			return nil, err
		}

		for _, cmd := range cmds {
			commands = append(commands, cmd)
		}
	}

	return commands, nil
}

func (i *Index) buildIndexField(id uint64, key []byte, value interface{}) ([]engine.Command, error) {
	var (
		err      error
		commands []engine.Command
	)

	v := reflect.TypeOf(value)

	err = nil

	switch v.Kind() {
	case reflect.String:
		commands, err = i.buildIndexString(id, key, value.(string))
		break
	case reflect.Int:
		commands, err = i.buildIndexInt64(id, key, value.(int64))
		break
	case reflect.Float64:
		commands, err = i.buildIndexFloat64(id, key, value.(float64))
		break
	case reflect.Slice:
		commands, err = i.buildIndexSlice(id, key, value.([]interface{}))
		break
	default:
		errMsg := fmt.Sprintf("Unknown type %s: %s\n", v.Kind(), value)

		if i.config.Debug {
			fmt.Printf(errMsg)
		}

		return nil, errors.New(errMsg)
	}

	return commands, err
}

func (i *Index) buildIndexSlice(id uint64, key []byte, values []interface{}) ([]engine.Command, error) {
	var commands []engine.Command

	storageName := string(key) + ".idx"

	if i.enableBatchMode {
		cmd, err := i.buildBatchOn(storageName)
		if err == nil {
			commands = append(commands, cmd)
		}
	}

	for _, value := range values {
		cmds, err := i.buildIndexField(id, key, value)

		if err != nil {
			return nil, err
		}

		for _, cmd := range cmds {
			commands = append(commands, cmd)
		}
	}

	return commands, nil
}

func (i *Index) buildIndexString(id uint64, key []byte, value string) ([]engine.Command, error) {
	var commands []engine.Command

	// default/hardcoded analyser == tokenizer
	value = strings.Trim(value, " ")
	value = strings.ToLower(value)
	tokens := strings.Split(value, " ")

	storageName := string(key) + ".idx"

	if i.enableBatchMode {
		cmd, err := i.buildBatchOn(storageName)
		if err == nil {
			commands = append(commands, cmd)
		}
	}

	// Index each token part
	// TODO: Optimize array of tokens. Need be *unique* tokens
	for _, t := range tokens {
		cmd := engine.Command{}
		cmd.Index = storageName
		cmd.Command = "mergeset"
		cmd.Key = []byte(t)
		cmd.KeyType = engine.TypeString
		cmd.Value = utils.Uint64ToBytes(id)
		cmd.ValueType = engine.TypeUint

		commands = append(commands, cmd)
	}

	// Index all string
	cmd := engine.Command{}
	cmd.Index = string(key) + ".idx"
	cmd.Command = "mergeset"
	cmd.Key = []byte(value)
	cmd.KeyType = engine.TypeString
	cmd.Value = utils.Uint64ToBytes(id)
	cmd.ValueType = engine.TypeUint

	commands = append(commands, cmd)
	return commands, nil
}

func (i *Index) buildIndexFloat64(id uint64, key []byte, value float64) ([]engine.Command, error) {
	var commands []engine.Command

	storageName := string(key) + ".idx"

	if i.enableBatchMode {
		cmd, err := i.buildBatchOn(storageName)

		if err == nil {
			commands = append(commands, cmd)
		}
	}

	cmd := engine.Command{}
	cmd.Index = storageName
	cmd.Command = "mergeset"
	cmd.Key = utils.Float64ToBytes(value)
	cmd.KeyType = engine.TypeFloat
	cmd.Value = utils.Uint64ToBytes(id)
	cmd.ValueType = engine.TypeUint

	commands = append(commands, cmd)
	return commands, nil
}

func (i *Index) buildIndexInt64(id uint64, key []byte, value int64) ([]engine.Command, error) {
	var commands []engine.Command

	storageName := string(key) + ".idx"

	if i.enableBatchMode {
		cmd, err := i.buildBatchOn(storageName)
		if err == nil {
			commands = append(commands, cmd)
		}
	}

	cmd := engine.Command{}
	cmd.Index = storageName
	cmd.Command = "mergeset"
	cmd.Key = utils.Int64ToBytes(value)
	cmd.KeyType = engine.TypeInt
	cmd.Value = utils.Uint64ToBytes(id)
	cmd.ValueType = engine.TypeUint

	commands = append(commands, cmd)

	return commands, nil
}

// Close the index
func (i *Index) Close() {
	i.engine.Close()
}
