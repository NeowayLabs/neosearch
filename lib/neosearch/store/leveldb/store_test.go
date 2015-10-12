package leveldb

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/NeowayLabs/neosearch/lib/neosearch/store"
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
		err     error
		kvstore store.KVStore
	)

	cfg := store.KVConfig{
		DataDir: DataDirTmp,
	}

	kvstore, err = NewLVDB(&cfg)
	if err != nil {
		t.Error(err)
		return nil
	} else if kvstore == nil {
		t.Error("Failed to allocate store")
		return nil
	}

	err = kvstore.Open(indexName, dbName)

	if err != nil {
		t.Error(err)
		return nil
	}

	return kvstore
}

func openDatabaseFail(t *testing.T, indexName, dbName string) {
	var (
		err     error
		kvstore store.KVStore
	)

	cfg := store.KVConfig{
		DataDir: DataDirTmp,
	}

	kvstore, err = NewLVDB(&cfg)
	if err != nil {
		t.Error(err)
		return
	} else if kvstore == nil {
		t.Error("Failed to allocate store")
		return
	}

	err = kvstore.Open(indexName, dbName)

	if err == nil {
		t.Errorf("Should fail... Invalid database name: %s", dbName)
		return
	}
}

func TestStoreHasBackend(t *testing.T) {
	cfg := store.KVConfig{
		DataDir: DataDirTmp,
	}

	kvstore, err := NewLVDB(&cfg)
	if err != nil {
		t.Errorf("You need compile this package with -tags <storage-backend>: %s", err)
		return
	}

	if kvstore == nil {
		t.Error("Failed to allocate KVStore")
	}
}

func TestOpenDatabase(t *testing.T) {
	shouldPass := []string{
		"123.tt",
		/*		"9999.db",
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
		store  store.KVStore
		data   []byte
		testDb = "test_set.db"
	)

	os.Mkdir(DataDirTmp+string(filepath.Separator)+"sample-store-set-get", 0755)
	store = openDatabase(t, "sample-store-set-get", testDb)

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

	writer := store.Writer()
	reader := store.Reader()

	for _, kv := range shouldPass {
		if err = writer.Set(kv.key, kv.value); err != nil {
			t.Error(err)
		}

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

	store.Close()

	os.RemoveAll(DataDirTmp + "/" + testDb)
}

func TestBatchWrite(t *testing.T) {
	var (
		err    error
		store  store.KVStore
		key    = []byte{'a'}
		value  = []byte{'b'}
		data   []byte
		testDb = "testbatch.db"
	)

	os.Mkdir(DataDirTmp+string(filepath.Separator)+"sample-batch-write", 0755)
	store = openDatabase(t, "sample-batch-write", testDb)

	reader := store.Reader()
	writer := store.Writer()

	writer.StartBatch()

	if writer.IsBatch() == false {
		t.Error("StartBatch not setting isBatch = true")
		return
	}

	if err = writer.Set(key, value); err != nil {
		t.Error(err)
		return
	}

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

	if data, err = reader.Get(key); err != nil {
		t.Error(err)
	} else if data == nil || len(data) != len(value) {
		t.Errorf("Failed to retrieve key '%s'. Retuns: %s", string(key), string(data))
	}

	if !reflect.DeepEqual(data, value) {
		t.Errorf("Data retrieved '%s' != '%s'", string(data), string(value))
	}

	store.Close()

	os.RemoveAll(DataDirTmp + "/" + testDb)
}

func TestBatchMultiWrite(t *testing.T) {
	var (
		err    error
		store  store.KVStore
		data   []byte
		testDb = "test_set-multi.db"
	)

	os.Mkdir(DataDirTmp+string(filepath.Separator)+"sample-batch-multi-write", 0755)
	store = openDatabase(t, "sample-batch-multi-write", testDb)

	reader := store.Reader()
	writer := store.Writer()

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

		if data, err := reader.Get(kv.key); err != nil || data != nil {
			t.Error("Key set before wasn't in the write batch cache." +
				" Batch mode isnt working")
		}
	}

	if err := writer.FlushBatch(); err != nil {
		t.Error(err)
	}

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

	store.Close()

	os.RemoveAll(DataDirTmp + "/" + testDb)
}
