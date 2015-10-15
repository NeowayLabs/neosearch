package goleveldb

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/NeowayLabs/neosearch/lib/neosearch/store"
	"github.com/NeowayLabs/neosearch/lib/neosearch/utils"
)

var DataDirTmp string

func init() {
	var err error
	DataDirTmp, err = ioutil.TempDir("/tmp", "neosearch-leveldb-")

	if err != nil {
		panic(err)
	}
}

func openDatabase(t *testing.T, indexName, dbName string) store.KVStore {
	var (
		err error
		kv  store.KVStore
	)

	cfg := store.KVConfig{
		"dataDir": DataDirTmp,
	}

	kv, err = NewLVDB(cfg)
	if err != nil {
		t.Error(err)
		return nil
	} else if kv == nil {
		t.Error("Failed to allocate store")
		return nil
	}

	err = kv.Open(indexName, dbName)
	if err != nil {
		t.Error(err)
		return nil
	}

	return kv
}

func openDatabaseFail(t *testing.T, indexName, dbName string) {
	var (
		err error
		kv  store.KVStore
	)

	cfg := store.KVConfig{
		"dataDir": DataDirTmp,
	}

	kv, err = NewLVDB(cfg)
	if err != nil {
		t.Error(err)
		return
	} else if kv == nil {
		t.Error("Failed to allocate store")
		return
	}

	err = kv.Open(indexName, dbName)

	if err == nil {
		t.Errorf("Should fail... Invalid database name: %s", dbName)
		return
	}
}

func TestStoreHasBackend(t *testing.T) {
	cfg := store.KVConfig{
		"dataDir": DataDirTmp,
	}

	kv, err := NewLVDB(cfg)
	if err != nil {
		t.Errorf("You need compile this package with -tags <storage-backend>: %s", err)
		return
	}

	if kv == nil {
		t.Error("Failed to allocate KVStore")
	}
}

func TestOpenDatabase(t *testing.T) {
	shouldPass := []string{
		"123.tt",
		/*      "9999.db",
		        "sample.db",
		        "sample.idx",
		        "sample_test.db",
		        "_id.db",
		        "_all.idx",
		        "__.idx",*/
	}

	shouldFail := []string{
		"",
		"1",
		"12",
		"123",
		"1234",
		".db",
		".idx",
		"...db",
		"sample",
		"sample.",
		"sample.a",
		"sample/test.db",
	}

	os.Mkdir(DataDirTmp+string(filepath.Separator)+"sample-ok", 0755)
	os.Mkdir(DataDirTmp+string(filepath.Separator)+"sample-fail", 0755)

	for _, dbname := range shouldPass {
		st := openDatabase(t, "sample-ok", dbname)
		if st != nil {
			st.Close()
		}

		os.RemoveAll(DataDirTmp + "/" + dbname)
	}

	for _, dbname := range shouldFail {
		openDatabaseFail(t, "sample-fail", dbname)
		//os.RemoveAll(DataDirTmp + "/" + dbname)
	}
}

func TestStoreSetGet(t *testing.T) {
	var (
		err    error
		kv     store.KVStore
		data   []byte
		testDb = "test_set.db"
	)

	os.Mkdir(DataDirTmp+string(filepath.Separator)+"sample-store-set-get", 0755)
	if kv = openDatabase(t, "sample-store-set-get", testDb); kv == nil {
		return
	}

	type kvTest struct {
		key   []byte
		value []byte
	}

	shouldPass := []kvTest{
		{
			key:   []byte{'t', 'e', 's', 't', 'e'},
			value: []byte{'i', '4', 'k'},
		},
		{
			key: []byte{'p', 'l', 'a', 'n', '9'},
			value: []byte{'f', 'r', 'o', 'm',
				'o', 'u', 't', 'e', 'r', 's',
				's', 'p', 'a', 'c', 'e', '!'},
		},
		{
			key:   []byte{'t', 'h', 'e', 'm', 'a', 't', 'r', 'i', 'x'},
			value: []byte{'h', 'a', 's', 'y', 'o', 'u'},
		},
	}

	writer := kv.Writer()
	if writer == nil {
		t.Error("Writer not created!")
		return
	}

	for _, kv := range shouldPass {
		if err = writer.Set(kv.key, kv.value); err != nil {
			t.Error(err)
		}
	}

	reader := kv.Reader()
	if reader == nil {
		t.Error("Reader not created!")
		return
	}
	defer func() {
		err := reader.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	for _, kv := range shouldPass {
		if data, err = reader.Get(kv.key); err != nil {
			t.Error(err)
			continue
		} else if data == nil || len(data) != len(kv.value) {
			t.Errorf("Failed to retrieve key '%s'. Retuns: %s", string(kv.key), string(data))
			continue
		}

		if !reflect.DeepEqual(data, kv.value) {
			t.Errorf("Data retrieved '%s' != '%s'", string(data), string(kv.value))
		}
	}

	data, err = reader.Get([]byte("do not exists"))
	if err != nil {
		t.Error(err)
	}

	// key does not exists, data should be nil
	if data != nil {
		t.Error("key 'does not exists' returning wrong value")
	}

	kv.Close()

	os.RemoveAll(DataDirTmp + "/" + testDb)
}

func TestBatchWrite(t *testing.T) {
	var (
		err    error
		kv     store.KVStore
		key    = []byte{'a'}
		value  = []byte{'b'}
		data   []byte
		testDb = "testbatch.db"
	)

	os.Mkdir(DataDirTmp+string(filepath.Separator)+"sample-batch-write", 0755)
	if kv = openDatabase(t, "sample-batch-write", testDb); kv == nil {
		return
	}

	writer := kv.Writer()
	if writer == nil {
		t.Error("Writer not created!")
		return
	}

	writer.StartBatch()

	if writer.IsBatch() == false {
		t.Error("StartBatch not setting isBatch = true")
		return
	}

	if err = writer.Set(key, value); err != nil {
		t.Error(err)
		return
	}

	reader := kv.Reader()
	if reader == nil {
		t.Error("Reader not created!")
		return
	}
	defer func() {
		err := reader.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	// should returns nil, nil because the key is in the batch cache
	if data, err = reader.Get(key); err != nil || data != nil {
		t.Error("Key set before wasn't in the write batch cache." +
			" Batch mode isnt working")
	}

	if err = writer.FlushBatch(); err != nil {
		t.Error(err)
	}

	if writer.IsBatch() == true {
		t.Error("FlushBatch doesnt reset the isBatch")
	}

	newReader := kv.Reader()
	if newReader == nil {
		t.Error("newReader not created!")
		return
	}
	defer func() {
		err := newReader.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	if data, err = newReader.Get(key); err != nil {
		t.Error(err)
	} else if data == nil || len(data) != len(value) {
		t.Errorf("Failed to retrieve key '%s'. Retuns: %s", string(key), string(data))
	}

	if !reflect.DeepEqual(data, value) {
		t.Errorf("Data retrieved '%s' != '%s'", string(data), string(value))
	}

	kv.Close()

	os.RemoveAll(DataDirTmp + "/" + testDb)
}

func TestBatchMultiWrite(t *testing.T) {
	var (
		err    error
		kv     store.KVStore
		data   []byte
		testDb = "test_set-multi.db"
	)

	os.Mkdir(DataDirTmp+string(filepath.Separator)+"sample-batch-multi-write", 0755)
	if kv = openDatabase(t, "sample-batch-multi-write", testDb); kv == nil {
		return
	}

	writer := kv.Writer()
	if writer == nil {
		t.Error("Writer not created!")
		return
	}

	writer.StartBatch()

	type kvTest struct {
		key   []byte
		value []byte
	}

	shouldPass := []kvTest{
		{
			key:   []byte{'t', 'e', 's', 't', 'e'},
			value: []byte{'i', '4', 'k'},
		},
		{
			key: []byte{'p', 'l', 'a', 'n', '9'},
			value: []byte{'f', 'r', 'o', 'm',
				'o', 'u', 't', 'e', 'r', 's',
				's', 'p', 'a', 'c', 'e', '!'},
		},
		{
			key:   []byte{'t', 'h', 'e', 'm', 'a', 't', 'r', 'i', 'x'},
			value: []byte{'h', 'a', 's', 'y', 'o', 'u'},
		},
	}

	for _, kv := range shouldPass {
		if err = writer.Set(kv.key, kv.value); err != nil {
			t.Error(err)
		}
	}

	reader := kv.Reader()
	if reader == nil {
		t.Error("Reader not created!")
		return
	}
	defer func() {
		err := reader.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	for _, kv := range shouldPass {
		if data, err := reader.Get(kv.key); err != nil || data != nil {
			t.Error("Key set before wasn't in the write batch cache." +
				" Batch mode isnt working")
		}
	}

	if err := writer.FlushBatch(); err != nil {
		t.Error(err)
	}

	newReader := kv.Reader()
	if newReader == nil {
		t.Error("Reader not created!")
		return
	}
	defer func() {
		err := newReader.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	for _, kv := range shouldPass {
		if data, err = newReader.Get(kv.key); err != nil {
			t.Error(err)
			continue
		} else if data == nil || len(data) != len(kv.value) {
			t.Errorf("Failed to retrieve key '%s'. Retuns: %s", string(kv.key), string(data))
			continue
		}

		if !reflect.DeepEqual(data, kv.value) {
			t.Errorf("Data retrieved '%s' != '%s'", string(data), string(kv.value))
		}
	}

	kv.Close()

	os.RemoveAll(DataDirTmp + "/" + testDb)
}

func TestStoreMergeSet(t *testing.T) {
	var (
		err    error
		kv     store.KVStore
		data   []byte
		testDb = "test_mergeset.db"
	)

	os.Mkdir(DataDirTmp+string(filepath.Separator)+"sample-store-merge-set", 0755)
	if kv = openDatabase(t, "sample-store-merge-set", testDb); kv == nil {
		return
	}

	writer := kv.Writer()
	if writer == nil {
		t.Error("Writer not created!")
		return
	}

	key := []byte{'t', 'e', 's', 't', 'e'}
	values := []uint64{0, 2, 1}

	result := append(utils.Uint64ToBytes(values[0]), utils.Uint64ToBytes(values[2])...)
	result = append(result, utils.Uint64ToBytes(values[1])...)

	for _, value := range values {
		if err = writer.MergeSet(key, value); err != nil {
			t.Error(err)
			return
		}
	}

	reader := kv.Reader()
	if reader == nil {
		t.Error("Reader not created!")
		return
	}
	defer func() {
		err := reader.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	if data, err = reader.Get(key); err != nil {
		t.Error(err)
	} else if data == nil || len(data) != len(result) {
		t.Errorf("Failed to retrieve key '%s'. Retuns: %+v", string(key), data)
	}

	if !reflect.DeepEqual(data, result) {
		t.Errorf("Data retrieved '%+v' != '%+v'", data, result)
	}

	kv.Close()
	os.RemoveAll(DataDirTmp + "/" + testDb)
}

func TestStoreIterator(t *testing.T) {
	var (
		err    error
		kv     store.KVStore
		testDb = "test_iterator.db"
	)

	os.Mkdir(DataDirTmp+string(filepath.Separator)+"sample-store-iterator", 0755)
	if kv = openDatabase(t, "sample-store-iterator", testDb); kv == nil {
		return
	}

	writer := kv.Writer()

	type kvTest struct {
		key   []byte
		value []byte
	}

	data := []kvTest{
		{
			key: []byte{'p', 'l', 'a', 'n', '9'},
			value: []byte{'f', 'r', 'o', 'm',
				'o', 'u', 't', 'e', 'r', 's',
				's', 'p', 'a', 'c', 'e', '!'},
		},
		{
			key:   []byte{'t', 'e', 's', 't', 'e'},
			value: []byte{'i', '4', 'k'},
		},
		{
			key:   []byte{'t', 'h', 'e', 'm', 'a', 't', 'r', 'i', 'x'},
			value: []byte{'h', 'a', 's', 'y', 'o', 'u'},
		},
	}

	writer.StartBatch()
	for _, kv := range data {
		if err = writer.Set(kv.key, kv.value); err != nil {
			t.Error(err)
		}
	}

	err = writer.FlushBatch()
	if err != nil {
		t.Fatal(err)
	}

	reader := kv.Reader()
	defer func() {
		err := reader.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()
	it := reader.GetIterator()
	it.Seek([]byte{0})
	keys := make([][]byte, 0, len(data))
	for it.Valid() {
		k := it.Key()
		key := make([]byte, len(k))
		copy(key, k)
		keys = append(keys, key)
		it.Next()
	}

	if len(keys) != len(data) {
		t.Errorf("expected same number of keys, got %d != %d", len(keys), len(data))
	}
	for i, dk := range data {
		if !reflect.DeepEqual(dk.key, keys[i]) {
			t.Errorf("expected key %s got %s", dk.key, keys[i])
		}
	}

	err = it.Close()
	if err != nil {
		t.Fatal(err)
	}
}
