package index

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"bitbucket.org/i4k/neosearch/engine"
	"bitbucket.org/i4k/neosearch/store"
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
	Name string

	engine *engine.Engine
	config Config

	// Indicates that index abstraction should batch each write command
	shouldBatch bool

	flushStorages []string
}

func validateIndexName(name string) bool {
	if len(name) < 3 {
		return false
	}

	validName := regexp.MustCompile(`^[a-zA-Z0-9-]+$`)
	return validName.MatchString(name)
}

// New creates new index
func New(name string, cfg Config, create bool) (*Index, error) {
	if !validateIndexName(name) {
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

	return nil
}

// Batch enables write cache of command before FlushBatch is executed
func (i *Index) Batch() {
	i.shouldBatch = true
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

// Add creates new document
func (i *Index) Add(id uint64, doc []byte) error {
	if i.shouldBatch {
		defer func() {
			i.shouldBatch = false
		}()
	}

	err := i.add(id, doc)

	if err != nil {
		return err
	}

	structData := map[string]interface{}{}

	err = json.Unmarshal(doc, &structData)

	if err != nil {
		return err
	}

	return i.indexFields(id, &structData)
}

func (i *Index) add(id uint64, doc []byte) error {
	if i.shouldBatch {
		err := i.enableBatchOn(dbName)

		if err != nil {
			return err
		}

		i.flushStorages = append(i.flushStorages, dbName)

		if i.config.Debug {
			fmt.Printf("Batch mode enabled for storage '%s' of index '%s'.\n",
				dbName,
				i.Name,
			)
		}
	}

	cmd := engine.Command{}

	cmd.Index = "document.db"
	cmd.Command = "set"
	cmd.Key = strconv.AppendUint([]byte(""), id, 10)
	cmd.Value = doc

	_, err := i.engine.Execute(cmd)
	return err
}

// Get retrieves the document by id
func (i *Index) Get(id uint64) ([]byte, error) {
	return i.engine.Execute(engine.Command{
		Index:   "document.db",
		Command: "get",
		Key:     strconv.AppendUint([]byte(""), id, 10),
	})
}

func (i *Index) enableBatchOn(storage string) error {
	if i.config.Debug {
		fmt.Printf("Batch mode enabled for storage '%s' of index '%s'.\n",
			storage,
			i.Name,
		)
	}

	_, err := i.engine.Execute(engine.Command{
		Index:   storage,
		Command: "batch",
	})

	return err
}

func (i *Index) indexFields(id uint64, structData *map[string]interface{}) error {
	var err error

	for key, value := range *structData {
		// optimize here?
		err = i.indexField(id, []byte(key), value)

		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Index) indexField(id uint64, key []byte, value interface{}) error {
	var err error
	v := reflect.TypeOf(value)

	err = nil

	switch v.Kind() {
	case reflect.String:
		err = i.indexString(id, key, value.(string))

		break
	case reflect.Int:
		err = i.indexInt(id, key, value.(int))
		break
	case reflect.Float64:
		err = i.indexFloat64(id, key, value.(float64))
		break
	default:
		fmt.Printf("Unknown type %s: %s\n", v.Kind(), value)
	}

	return err
}

func (i *Index) indexString(id uint64, key []byte, value string) error {
	// default/hardcoded analyser == tokenizer
	value = strings.ToLower(value)
	tokens := strings.Split(value, " ")

	storageName := string(key) + ".idx"

	if i.shouldBatch {
		if err := i.enableBatchOn(storageName); err != nil {
			return err
		}

		i.flushStorages = append(i.flushStorages, storageName)
	}

	// Index each token part
	// TODO: Optimize array of tokens. Need be *unique* tokens
	for _, t := range tokens {
		cmd := engine.Command{}
		cmd.Index = storageName
		cmd.Command = "mergeset"
		cmd.Key = []byte(t)
		cmd.Value = strconv.AppendUint([]byte(""), id, 10)

		_, err := i.engine.Execute(cmd)

		if err != nil {
			return err
		}
	}

	// Index all string
	cmd := engine.Command{}
	cmd.Index = string(key) + ".idx"
	cmd.Command = "mergeset"
	cmd.Key = []byte(value)
	cmd.Value = strconv.AppendUint([]byte(""), id, 10)

	_, err := i.engine.Execute(cmd)

	return err
}

func (i *Index) indexFloat64(id uint64, key []byte, value float64) error {
	storageName := string(key) + ".idx"

	if i.shouldBatch {
		if err := i.enableBatchOn(storageName); err != nil {
			return err
		}

		i.flushStorages = append(i.flushStorages, storageName)
	}

	cmd := engine.Command{}
	cmd.Index = storageName
	cmd.Command = "set"
	cmd.Key = strconv.AppendFloat([]byte(""), value, 'f', -1, 64)
	cmd.Value = strconv.AppendUint([]byte(""), id, 10)

	_, err := i.engine.Execute(cmd)
	return err
}

func (i *Index) indexInt(id uint64, key []byte, value int) error {
	return nil
}

// Close the index
func (i *Index) Close() {
	i.engine.Close()
}
