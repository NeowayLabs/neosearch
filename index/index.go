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

// Index represents an entire index
type Index struct {
	Name string

	engine *engine.Engine
	config Config
}

// Config index
type Config struct {
	// Per-Index configurations

	DataDir     string
	Debug       bool
	CacheSize   int
	EnableCache bool
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

// Add creates new document
func (i *Index) Add(id uint64, doc []byte) error {
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

	for _, t := range tokens {
		cmd := engine.Command{}
		cmd.Index = string(key) + ".idx"
		cmd.Command = "mergeset"
		cmd.Key = []byte(t)
		cmd.Value = strconv.AppendUint([]byte(""), id, 10)

		_, err := i.engine.Execute(cmd)

		if err != nil {
			return err
		}
	}

	cmd := engine.Command{}
	cmd.Index = string(key) + ".idx"
	cmd.Command = "mergeset"
	cmd.Key = []byte(value)
	cmd.Value = strconv.AppendUint([]byte(""), id, 10)

	_, err := i.engine.Execute(cmd)

	return err
}

func (i *Index) indexFloat64(id uint64, key []byte, value float64) error {
	cmd := engine.Command{}
	cmd.Index = string(key) + ".idx"
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
