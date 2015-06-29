package index

import (
	"os"
	"testing"

	"github.com/NeowayLabs/neosearch/lib/neosearch/engine"
	"github.com/NeowayLabs/neosearch/lib/neosearch/utils"
)

func TestSimpleIndexWithMetadata(t *testing.T) {
	var (
		indexName                  = "document-sample-metadata"
		indexDir                   = DataDirTmp + "/" + indexName
		commands, expectedCommands []engine.Command
		docJSON                    = []byte(`{"id": 1}`)
		err                        error
		index                      *Index
		metadata                   = Metadata{
			"id": Metadata{
				"type": "uint",
			},
		}
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

	commands, err = index.BuildAdd(1, docJSON, metadata)

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
			Key:       utils.Uint64ToBytes(1),
			KeyType:   engine.TypeUint,
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

	metadata = Metadata{
		"title": Metadata{
			"type": "string",
		},
		"description": Metadata{
			"type": "string",
		},
	}

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

	commands, err = index.BuildAdd(2, docJSON, metadata)

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
