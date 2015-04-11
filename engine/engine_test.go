package engine

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/NeowayLabs/neosearch/store"
	"github.com/NeowayLabs/neosearch/utils"
)

var DataDirTmp string
var sampleIndex = "sample"

func init() {
	var err error
	DataDirTmp, err = ioutil.TempDir("/tmp", "neosearch-index-")

	if err != nil {
		panic(err)
	}

	os.Mkdir(DataDirTmp+"/"+sampleIndex, 0755)
}

func execSequence(t *testing.T, ng *Engine, cmds []Command) {
	for _, cmd := range cmds {
		_, err := ng.Execute(cmd)

		if err != nil {
			t.Error(err)
		}
	}
}

func cmpIterator(t *testing.T, itReturns []map[int64]string, ng *Engine, seek []byte, index, database string) {
	storekv, err := ng.GetStore(index, database)

	if err != nil {
		t.Error(err)
	}

	it := storekv.GetIterator()

	it.Seek(seek)

	for i := 0; i < len(itReturns); i++ {
		for key, value := range itReturns[i] {
			if !it.Valid() {
				t.Errorf("Failed to seek to '%d' key", key)
			}

			val := it.Value()

			if len(val) == 0 || string(val) != value {
				t.Errorf("Failed to get '%d' key", key)
			} else if err := it.GetError(); err != nil {
				t.Error(err)
			}
		}

		it.Next()
	}
}

// TestEngineIntegerKeyOrder verifies if the chosen storage engine is really
// a LSM database ordered by key with ByteWise comparator.
func TestEngineIntegerKeyOrder(t *testing.T) {
	ng := New(NGConfig{
		KVCfg: &store.KVConfig{
			DataDir: DataDirTmp,
		},
		OpenCacheSize: 1,
	})

	cmds := []Command{
		{
			Index:    "sample",
			Database: "test.idx",
			Key:      []byte("AAA"),
			Value:    []byte("value AAA"),
			Command:  "set",
		},
		{
			Index:    "sample",
			Database: "test.idx",
			Key:      utils.Uint64ToBytes(1),
			Value:    []byte("value 1"),
			Command:  "set",
		},
		{
			Index:    "sample",
			Database: "test.idx",
			Key:      []byte("BBB"),
			Value:    []byte("value BBB"),
			Command:  "set",
		},
		{
			Index:    "sample",
			Database: "test.idx",
			Key:      utils.Uint64ToBytes(2),
			Value:    []byte("value 2"),
			Command:  "set",
		},
		{
			Index:    "sample",
			Database: "test.idx",
			Key:      utils.Uint64ToBytes(2000),
			Value:    []byte("value 2000"),
			Command:  "set",
		},
		{
			Index:    "sample",
			Database: "test.idx",
			Key:      []byte("2000"),
			Value:    []byte("value 2000"),
			Command:  "set",
		},
		{
			Index:    "sample",
			Database: "test.idx",
			Key:      utils.Uint64ToBytes(100000),
			Value:    []byte("value 100000"),
			Command:  "set",
		},
		{
			Index:    "sample",
			Database: "test.idx",
			Key:      utils.Uint64ToBytes(1000000),
			Value:    []byte("value 1000000"),
			Command:  "set",
		},
		{
			Index:    "sample",
			Database: "test.idx",
			Key:      utils.Uint64ToBytes(10000000000000),
			Value:    []byte("value 10000000000000"),
			Command:  "set",
		},
	}

	execSequence(t, ng, cmds)

	itReturns := []map[int64]string{
		{
			1: "value 1",
		},
		{
			2: "value 2",
		},
		{
			2000: "value 2000",
		},
		{
			100000: "value 100000",
		},
		{
			1000000: "value 1000000",
		},
		{
			10000000000000: "value 10000000000000",
		},
	}

	cmpIterator(t, itReturns, ng, utils.Uint64ToBytes(1), "sample", "test.idx")

	defer ng.Close()
	os.RemoveAll(DataDirTmp)
}
