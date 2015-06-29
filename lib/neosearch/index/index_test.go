package index

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/NeowayLabs/neosearch/lib/neosearch/engine"
	"github.com/NeowayLabs/neosearch/lib/neosearch/utils"
)

var DataDirTmp string

func init() {
	var err error
	DataDirTmp, err = ioutil.TempDir("/tmp", "neosearch-index-index")

	if err != nil {
		panic(err)
	}
}

func compareCommands(t *testing.T, commands, expectedCommands []engine.Command) bool {
	if len(commands) != len(expectedCommands) {
		t.Errorf("Length differs: len(cmd) == %d != len(expected) == %d\n", len(commands), len(expectedCommands))

		for _, cmd := range commands {
			cmd.Println()
		}
		return false
	}

	for idx, cmd := range commands {
		expectCmd := expectedCommands[idx]

		if !reflect.DeepEqual(cmd, expectCmd) {
			t.Errorf("Commands differs !!\n")
			cmd.Println()
			expectCmd.Println()
			return false
		}
	}

	return true
}

func TestBuildAddDocument(t *testing.T) {
	var (
		indexName                  = "document-sample"
		indexDir                   = DataDirTmp + "/" + indexName
		commands, expectedCommands []engine.Command
		docJSON                    = []byte(`{"id": 1}`)
		err                        error
		index                      *Index
	)

	cfg := Config{
		Debug:   false,
		DataDir: DataDirTmp,
	}

	err = os.MkdirAll(DataDirTmp, 0755)

	if err != nil {
		goto cleanup
	}

	index, err = New(indexName, cfg, true)

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	if _, err := os.Stat(indexDir); os.IsNotExist(err) {
		t.Errorf("no such file or directory: %s", indexDir)
		goto cleanup
	}

	commands, err = index.BuildAdd(1, docJSON, nil)

	if err != nil {
		t.Error(err.Error())
		goto cleanup
	}

	expectedCommands = []engine.Command{
		{
			Index:     indexName,
			Database:  "document.db",
			Key:       utils.Uint64ToBytes(1),
			KeyType:   engine.TypeUint,
			Value:     docJSON,
			ValueType: engine.TypeString,
			Command:   "set",
		},
		{
			Index:     indexName,
			Database:  "id.idx",
			Key:       utils.Float64ToBytes(1),
			KeyType:   engine.TypeFloat,
			Value:     utils.Uint64ToBytes(1),
			ValueType: engine.TypeUint,
			Command:   "mergeset",
		},
	}

	if !compareCommands(t, commands, expectedCommands) {
		goto cleanup
	}

	docJSON = []byte(`{
            "title": "NeoSearch - Reverse Index",
            "description": "Neoway Full Text Search"
        }`)

	expectedCommands = []engine.Command{
		{
			Index:     indexName,
			Database:  "document.db",
			Command:   "set",
			Key:       utils.Uint64ToBytes(2),
			KeyType:   engine.TypeUint,
			Value:     docJSON,
			ValueType: engine.TypeString,
		},
		{
			Index:     indexName,
			Database:  "description.idx",
			Command:   "mergeset",
			Key:       []byte("neoway"),
			KeyType:   engine.TypeString,
			Value:     utils.Uint64ToBytes(2),
			ValueType: engine.TypeUint,
		},
		{
			Index:     indexName,
			Database:  "description.idx",
			Command:   "mergeset",
			Key:       []byte("full"),
			KeyType:   engine.TypeString,
			Value:     utils.Uint64ToBytes(2),
			ValueType: engine.TypeUint,
		},
		{
			Index:     indexName,
			Database:  "description.idx",
			Command:   "mergeset",
			Key:       []byte("text"),
			KeyType:   engine.TypeString,
			Value:     utils.Uint64ToBytes(2),
			ValueType: engine.TypeUint,
		},
		{
			Index:     indexName,
			Database:  "description.idx",
			Command:   "mergeset",
			Key:       []byte("search"),
			KeyType:   engine.TypeString,
			Value:     utils.Uint64ToBytes(2),
			ValueType: engine.TypeUint,
		},
		{
			Index:     indexName,
			Database:  "description.idx",
			Command:   "mergeset",
			Key:       []byte("neoway full text search"),
			KeyType:   engine.TypeString,
			Value:     utils.Uint64ToBytes(2),
			ValueType: engine.TypeUint,
		},
		{
			Index:     indexName,
			Database:  "title.idx",
			Command:   "mergeset",
			Key:       []byte("neosearch"),
			KeyType:   engine.TypeString,
			Value:     utils.Uint64ToBytes(2),
			ValueType: engine.TypeUint,
		},
		{
			Index:     indexName,
			Database:  "title.idx",
			Command:   "mergeset",
			Key:       []byte("-"),
			KeyType:   engine.TypeString,
			Value:     utils.Uint64ToBytes(2),
			ValueType: engine.TypeUint,
		},
		{
			Index:     indexName,
			Database:  "title.idx",
			Command:   "mergeset",
			Key:       []byte("reverse"),
			KeyType:   engine.TypeString,
			Value:     utils.Uint64ToBytes(2),
			ValueType: engine.TypeUint,
		},
		{
			Index:     indexName,
			Database:  "title.idx",
			Command:   "mergeset",
			Key:       []byte("index"),
			KeyType:   engine.TypeString,
			Value:     utils.Uint64ToBytes(2),
			ValueType: engine.TypeUint,
		},
		{
			Index:     indexName,
			Database:  "title.idx",
			Command:   "mergeset",
			Key:       []byte("neosearch - reverse index"),
			KeyType:   engine.TypeString,
			Value:     utils.Uint64ToBytes(2),
			ValueType: engine.TypeUint,
		},
	}

	commands, err = index.BuildAdd(2, docJSON, nil)

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	if !compareCommands(t, commands, expectedCommands) {
		goto cleanup
	}

cleanup:
	index.Close()
	os.RemoveAll(indexDir)
}

func TestBuildAddDocumentWithBatchMode(t *testing.T) {
	var (
		indexName                  = "document-sample2"
		indexDir                   = DataDirTmp + "/" + indexName
		commands, expectedCommands []engine.Command
		docJSON                    = []byte(`{"id": 1}`)
		err                        error
		index                      *Index
	)

	cfg := Config{
		Debug:   false,
		DataDir: DataDirTmp,
	}

	err = os.MkdirAll(DataDirTmp, 0755)

	if err != nil {
		goto cleanup
	}

	index, err = New(indexName, cfg, true)

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	if _, err := os.Stat(indexDir); os.IsNotExist(err) {
		t.Errorf("no such file or directory: %s", indexDir)
		goto cleanup
	}

	index.Batch()
	commands, err = index.BuildAdd(1, docJSON, nil)

	if err != nil {
		t.Error(err.Error())
		goto cleanup
	}

	expectedCommands = []engine.Command{
		{
			Index:     indexName,
			Database:  "document.db",
			Command:   "batch",
			Key:       nil,
			KeyType:   engine.TypeNil,
			Value:     nil,
			ValueType: engine.TypeNil,
		},
		{
			Index:     indexName,
			Database:  "document.db",
			Key:       utils.Uint64ToBytes(1),
			KeyType:   engine.TypeUint,
			Value:     docJSON,
			ValueType: engine.TypeString,
			Command:   "set",
		},
		{
			Index:     indexName,
			Database:  "id.idx",
			Command:   "batch",
			Key:       nil,
			KeyType:   engine.TypeNil,
			Value:     nil,
			ValueType: engine.TypeNil,
		},
		{
			Index:     indexName,
			Database:  "id.idx",
			Key:       utils.Float64ToBytes(1),
			KeyType:   engine.TypeFloat,
			Value:     utils.Uint64ToBytes(1),
			ValueType: engine.TypeUint,
			Command:   "mergeset",
		},
	}

	if !compareCommands(t, commands, expectedCommands) {
		goto cleanup
	}

	docJSON = []byte(`{
            "title": "NeoSearch - Reverse Index",
            "description": "Neoway Full Text Search"
        }`)

	expectedCommands = []engine.Command{
		// The "document.db" doesn't need a batch command,
		// because this was already done by document '1'
		//{
		//	Index:   "document.db",
		//	Command: "batch",
		//},
		{
			Index:     indexName,
			Database:  "document.db",
			Command:   "set",
			Key:       utils.Uint64ToBytes(2),
			KeyType:   engine.TypeUint,
			Value:     docJSON,
			ValueType: engine.TypeString,
		},
		{
			Index:     indexName,
			Database:  "description.idx",
			Command:   "batch",
			Key:       nil,
			KeyType:   engine.TypeNil,
			Value:     nil,
			ValueType: engine.TypeNil,
		},
		{
			Index:     indexName,
			Database:  "description.idx",
			Command:   "mergeset",
			Key:       []byte("neoway"),
			KeyType:   engine.TypeString,
			Value:     utils.Uint64ToBytes(2),
			ValueType: engine.TypeUint,
		},
		{
			Index:     indexName,
			Database:  "description.idx",
			Command:   "mergeset",
			Key:       []byte("full"),
			KeyType:   engine.TypeString,
			Value:     utils.Uint64ToBytes(2),
			ValueType: engine.TypeUint,
		},
		{
			Index:     indexName,
			Database:  "description.idx",
			Command:   "mergeset",
			Key:       []byte("text"),
			KeyType:   engine.TypeString,
			Value:     utils.Uint64ToBytes(2),
			ValueType: engine.TypeUint,
		},
		{
			Index:     indexName,
			Database:  "description.idx",
			Command:   "mergeset",
			Key:       []byte("search"),
			KeyType:   engine.TypeString,
			Value:     utils.Uint64ToBytes(2),
			ValueType: engine.TypeUint,
		},
		{
			Index:     indexName,
			Database:  "description.idx",
			Command:   "mergeset",
			Key:       []byte("neoway full text search"),
			KeyType:   engine.TypeString,
			Value:     utils.Uint64ToBytes(2),
			ValueType: engine.TypeUint,
		},
		{
			Index:     indexName,
			Database:  "title.idx",
			Command:   "batch",
			Key:       nil,
			KeyType:   engine.TypeNil,
			Value:     nil,
			ValueType: engine.TypeNil,
		},
		{
			Index:     indexName,
			Database:  "title.idx",
			Command:   "mergeset",
			Key:       []byte("neosearch"),
			KeyType:   engine.TypeString,
			Value:     utils.Uint64ToBytes(2),
			ValueType: engine.TypeUint,
		},
		{
			Index:     indexName,
			Database:  "title.idx",
			Command:   "mergeset",
			Key:       []byte("-"),
			KeyType:   engine.TypeString,
			Value:     utils.Uint64ToBytes(2),
			ValueType: engine.TypeUint,
		},
		{
			Index:     indexName,
			Database:  "title.idx",
			Command:   "mergeset",
			Key:       []byte("reverse"),
			KeyType:   engine.TypeString,
			Value:     utils.Uint64ToBytes(2),
			ValueType: engine.TypeUint,
		},
		{
			Index:     indexName,
			Database:  "title.idx",
			Command:   "mergeset",
			Key:       []byte("index"),
			KeyType:   engine.TypeString,
			Value:     utils.Uint64ToBytes(2),
			ValueType: engine.TypeUint,
		},
		{
			Index:     indexName,
			Database:  "title.idx",
			Command:   "mergeset",
			Key:       []byte("neosearch - reverse index"),
			KeyType:   engine.TypeString,
			Value:     utils.Uint64ToBytes(2),
			ValueType: engine.TypeUint,
		},
	}

	index.Batch()
	commands, err = index.BuildAdd(2, docJSON, nil)

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	if !compareCommands(t, commands, expectedCommands) {
		goto cleanup
	}

cleanup:
	index.Close()
	os.RemoveAll(indexDir)
}
