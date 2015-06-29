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

	"github.com/NeowayLabs/neosearch/lib/neosearch/engine"
	"github.com/NeowayLabs/neosearch/lib/neosearch/store"
	"github.com/NeowayLabs/neosearch/lib/neosearch/utils"
)

const (
	dbName   string = "document.db"
	indexExt string = "idx"
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

	fullDir string
}

// ValidateIndexName verifies if name is valid NeoSearch index name
func ValidateIndexName(name string) bool {
	if len(name) < 3 {
		return false
	}

	validName := regexp.MustCompile(`^[a-zA-Z]+[a-zA-Z0-9_-]+$`)
	return validName.MatchString(name)
}

// New creates new index
func New(name string, cfg Config, create bool) (*Index, error) {
	if !ValidateIndexName(name) {
		return nil, errors.New("Invalid index name")
	}

	index := &Index{
		Name:   name,
		config: cfg,
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

	i.fullDir = dataDir

	i.engine = engine.New(engine.NGConfig{
		KVCfg: &store.KVConfig{
			DataDir:     i.config.DataDir,
			Debug:       i.config.Debug,
			CacheSize:   i.config.CacheSize,
			EnableCache: i.config.EnableCache,
		},
	})

	return nil
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
			Index:    i.Name,
			Database: storeName,
			Command:  "flushbatch",
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
func (i *Index) Add(id uint64, doc []byte, metadata map[string]interface{}) error {
	if metadata == nil {
		metadata = Metadata{}
	}

	commands, err := i.BuildAdd(id, doc, metadata)

	if err != nil {
		return err
	}

	for _, cmd := range commands {
		_, err := i.engine.Execute(cmd)

		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Index) BuildAdd(id uint64, doc []byte, metadata Metadata) ([]engine.Command, error) {
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

	if len(structData) == 0 {
		return nil, errors.New("Empty document")
	}

	fieldCommands, err := i.buildIndexFields(id, "", structData, metadata)

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

	commands = make([]engine.Command, 0, 2)

	if i.enableBatchMode {
		cmd, err := i.buildBatchOn(dbName)

		if err == nil {
			commands = append(commands, cmd)
		}
	}

	cmd := engine.Command{}
	cmd.Database = dbName
	cmd.Index = i.Name
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
		Index:    i.Name,
		Database: dbName,
		Command:  "get",
		Key:      utils.Uint64ToBytes(id),
		KeyType:  engine.TypeUint,
	}
}

// Get retrieves the document by id
func (i *Index) Get(id uint64) ([]byte, error) {
	return i.engine.Execute(i.buildGet(id))
}

func (i *Index) GetAnalyze(id uint64) (engine.Command, error) {
	return i.buildGet(id), nil
}

// GetDocs returns the content of documents specified by docIDs and limited
// by limit.
func (i *Index) GetDocs(docIDs []uint64, limit uint) ([]string, error) {
	var (
		docLen = uint(len(docIDs))
	)

	if docLen > limit {
		docLen = limit
	}

	docs := make([]string, docLen)

	for idx, docID := range docIDs {
		if uint(idx) == docLen {
			break
		}

		if byteDoc, err := i.Get(docID); err == nil {
			docs[idx] = string(byteDoc)
		} else {
			return nil, err
		}
	}

	return docs, nil
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
		Index:     i.Name,
		Database:  storage,
		Command:   "batch",
		Key:       nil,
		KeyType:   engine.TypeNil,
		Value:     nil,
		ValueType: engine.TypeNil,
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
func (i *Index) buildIndexFields(id uint64, baseField string, structData map[string]interface{}, metadata Metadata) ([]engine.Command, error) {
	var (
		commands []engine.Command
		dataKeys []string
	)

	for key := range structData {
		dataKeys = append(dataKeys, key)
	}

	sort.Strings(dataKeys)

	for _, key := range dataKeys {
		metainfo, ok := metadata[key].(Metadata)

		if !ok {
			if i.config.Debug {
				fmt.Printf("[WARN] Metadata not supplied for field '%s'. Metadata = %+v\n", key, metadata)
			}

			metainfo = nil
		}

		value := structData[key]

		fieldKey := key

		if baseField != "" {
			fieldKey = baseField + "." + fieldKey
		}

		cmds, err := i.buildIndexField(id, []byte(fieldKey), value, metainfo)

		if err != nil {
			return nil, err
		}

		for _, cmd := range cmds {
			commands = append(commands, cmd)
		}
	}

	return commands, nil
}

func (i *Index) buildIndexField(id uint64, key []byte, value interface{}, metadata Metadata) ([]engine.Command, error) {
	var (
		commands  []engine.Command
		err       error
		ok        bool
		fieldType string
	)

	vtype := reflect.TypeOf(value)
	fieldType = vtype.String()

	if metadata != nil {
		fieldType, ok = metadata["type"].(string)

		if !ok {
			return nil, fmt.Errorf("Invalid metadata. Field 'type' is required: %+v", metadata)
		}
	}

	switch strings.ToLower(fieldType) {
	case "string":
		vstr, ok := value.(string)

		if !ok {
			return nil, fmt.Errorf("Error indexing field '%s'. Value '%+v' isn't string", string(key), value)
		}

		commands, err = i.buildIndexString(id, key, vstr)
	case "uint", "uint8", "uint16", "uint32", "uint64":
		var vuint uint64

		vuint, ok := value.(uint64)

		if !ok {
			vuint, err = utils.Uint64FromInterface(value, vtype.Kind())

			if err != nil {
				return nil, fmt.Errorf("Error indexing field '%s'. Value '%+v' isn't uint", string(key), value)
			}
		}

		commands, err = i.buildIndexUint64(id, key, vuint)
	case "int", "int8", "int16", "int32", "int64":
		var vint int64

		vint, ok := value.(int64)

		if !ok {
			vint, err = utils.Int64FromInterface(value, vtype.Kind())

			if err != nil {
				return nil, fmt.Errorf("Error indexing field '%s'. Value '%+v' isn't int", string(key), value)
			}
		}

		commands, err = i.buildIndexInt64(id, key, vint)
	case "float", "float32", "float64":
		vfloat, ok := value.(float64)

		if !ok {
			return nil, fmt.Errorf("Error indexing field '%s'. Value '%+v' isn't int", string(key), value)
		}

		commands, err = i.buildIndexFloat64(id, key, vfloat)
	case "slice", "list", "[]interface {}":
		vslice, ok := value.([]interface{})

		if !ok {
			return nil, fmt.Errorf("Error indexing field '%s'. Value '%+v' isn't slice", string(key), value)
		}

		submetadata, ok := metadata["metadata"].(Metadata)

		if !ok {
			submetadata = nil
		}

		commands, err = i.buildIndexSlice(id, key, vslice, submetadata)
	case "object", "map", "map[string]interface {}":
		vobject, ok := value.(map[string]interface{})

		if !ok {
			return nil, fmt.Errorf("Error indexing field '%s'. Value '%+v' isn't object", string(key), value)
		}

		submetadata, ok := metadata["metadata"].(Metadata)

		if !ok {
			submetadata = nil
		}

		commands, err = i.buildIndexFields(id, string(key), vobject, submetadata)
	default:
		errMsg := fmt.Sprintf("Unknown type %s: %s\n", fieldType, value)

		if i.config.Debug {
			fmt.Printf(errMsg)
		}

		return nil, errors.New(errMsg)
	}

	return commands, err
}

// TODO: Index don't take care of item order
func (i *Index) buildIndexSlice(id uint64, key []byte, values []interface{}, metadata Metadata) ([]engine.Command, error) {
	var commands []engine.Command

	storageName := string(key) + ".idx"

	if i.enableBatchMode {
		cmd, err := i.buildBatchOn(storageName)
		if err == nil {
			commands = append(commands, cmd)
		}
	}

	for _, value := range values {
		cmds, err := i.buildIndexField(id, key, value, metadata)

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

	addIndexStringCommand := func(dbase string, key []byte) {
		cmd := engine.Command{}
		cmd.Index = i.Name
		cmd.Database = dbase
		cmd.Command = "mergeset"
		cmd.Key = key
		cmd.KeyType = engine.TypeString
		cmd.Value = utils.Uint64ToBytes(id)
		cmd.ValueType = engine.TypeUint

		commands = append(commands, cmd)
	}

	// Index each token part
	// TODO: Optimize array of tokens. Need be *unique* tokens
	for _, t := range tokens {
		addIndexStringCommand(storageName, []byte(t))
	}

	if len(tokens) == 1 {
		// if there's one token, then no need for index entire string
		return commands, nil
	}

	// Index all string
	addIndexStringCommand(string(key)+".idx", []byte(value))
	return commands, nil
}

func (i *Index) buildIndexCommands(key []byte, cmdKey []byte, cmdVal []byte, keyType uint8) ([]engine.Command, error) {
	var commands []engine.Command

	storageName := string(key) + ".idx"

	if i.enableBatchMode {
		cmd, err := i.buildBatchOn(storageName)
		if err == nil {
			commands = append(commands, cmd)
		}
	}

	cmd := engine.Command{}
	cmd.Index = i.Name
	cmd.Database = storageName
	cmd.Command = "mergeset"
	cmd.Key = cmdKey
	cmd.KeyType = keyType
	cmd.Value = cmdVal
	cmd.ValueType = engine.TypeUint

	commands = append(commands, cmd)
	return commands, nil
}

func (i *Index) buildIndexFloat64(id uint64, key []byte, value float64) ([]engine.Command, error) {
	return i.buildIndexCommands(key, utils.Float64ToBytes(value), utils.Uint64ToBytes(id), engine.TypeFloat)
}

func (i *Index) buildIndexUint64(id uint64, key []byte, value uint64) ([]engine.Command, error) {
	return i.buildIndexCommands(key, utils.Uint64ToBytes(value), utils.Uint64ToBytes(id), engine.TypeUint)
}

func (i *Index) buildIndexInt64(id uint64, key []byte, value int64) ([]engine.Command, error) {
	return i.buildIndexCommands(key, utils.Int64ToBytes(value), utils.Uint64ToBytes(id), engine.TypeInt)
}

// Close the index
func (i *Index) Close() {
	i.engine.Close()
}
