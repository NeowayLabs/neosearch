package index

import (
	"os"
	"testing"

	"github.com/NeowayLabs/neosearch/lib/neosearch/engine"
	"github.com/NeowayLabs/neosearch/lib/neosearch/utils"
)

func TestBuildAddObjectDocument(t *testing.T) {
	var (
		indexName                  = "document-with-object-sample"
		indexDir                   = DataDirTmp + "/" + indexName
		commands, expectedCommands []engine.Command
		docJSON                    = []byte(`
                {
                    "id": 1,
                    "address": {
                        "city": "florianópolis",
                        "district": "Itacorubi",
                        "street": "Patricio Farias",
                        "latlon": [
                            -27.545198,
                            -48.504827
                        ]
                    }
                }`)
		err   error
		index *Index
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

	commands, err = index.BuildAdd(1, docJSON)

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
			Database:  "address.city.idx",
			Key:       []byte("florianópolis"),
			KeyType:   engine.TypeString,
			Value:     utils.Uint64ToBytes(1),
			ValueType: engine.TypeUint,
			Command:   "mergeset",
		},
		{
			Index:     indexName,
			Database:  "address.district.idx",
			Key:       []byte("itacorubi"),
			KeyType:   engine.TypeString,
			Value:     utils.Uint64ToBytes(1),
			ValueType: engine.TypeUint,
			Command:   "mergeset",
		},
		{
			Index:     indexName,
			Database:  "address.latlon.idx",
			Key:       utils.Float64ToBytes(-27.545198),
			KeyType:   engine.TypeFloat,
			Value:     utils.Uint64ToBytes(1),
			ValueType: engine.TypeUint,
			Command:   "mergeset",
		},
		{
			Index:     indexName,
			Database:  "address.latlon.idx",
			Key:       utils.Float64ToBytes(-48.504827),
			KeyType:   engine.TypeFloat,
			Value:     utils.Uint64ToBytes(1),
			ValueType: engine.TypeUint,
			Command:   "mergeset",
		},
		{
			Index:     indexName,
			Database:  "address.street.idx",
			Key:       []byte("patricio"),
			KeyType:   engine.TypeString,
			Value:     utils.Uint64ToBytes(1),
			ValueType: engine.TypeUint,
			Command:   "mergeset",
		},
		{
			Index:     indexName,
			Database:  "address.street.idx",
			Key:       []byte("farias"),
			KeyType:   engine.TypeString,
			Value:     utils.Uint64ToBytes(1),
			ValueType: engine.TypeUint,
			Command:   "mergeset",
		},
		{
			Index:     indexName,
			Database:  "address.street.idx",
			Key:       []byte("patricio farias"),
			KeyType:   engine.TypeString,
			Value:     utils.Uint64ToBytes(1),
			ValueType: engine.TypeUint,
			Command:   "mergeset",
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

cleanup:
	index.Close()
	os.RemoveAll(indexDir)
}
