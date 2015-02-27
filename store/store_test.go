package store

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

var DataDirTmp string

func init() {
	var err error
	DataDirTmp, err = ioutil.TempDir("/tmp", "neosearch-")

	if err != nil {
		panic(err)
	}
}

func openDatabase(t *testing.T, name string) KVStore {
	var (
		err   error
		store KVStore
	)

	cfg := KVConfig{
		DataDir: DataDirTmp,
	}

	store, err = New(&cfg)

	if err != nil {
		t.Error(err)
		return nil
	} else if store == nil {
		t.Error("Failed to allocate store")
		return nil
	}

	err = store.Open(name)

	if err != nil {
		t.Error(err)
		return nil
	}

	return store
}

func openDatabaseFail(t *testing.T, name string) {
	var (
		err   error
		store KVStore
	)

	cfg := KVConfig{
		DataDir: DataDirTmp,
	}

	store, err = New(&cfg)

	if err != nil {
		t.Error(err)
		return
	} else if store == nil {
		t.Error("Failed to allocate store")
		return
	}

	err = store.Open(name)

	if err == nil {
		t.Errorf("Should fail... Invalid database name: %s", name)
		return
	}
}

func TestStoreHasBackend(t *testing.T) {
	cfg := KVConfig{
		DataDir: DataDirTmp,
	}

	store, err := New(&cfg)

	if err != nil {
		t.Errorf("You need compile this package with -tags <storage-backend>: %s", err)
		return
	}

	if store == nil {
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
		"asas-a.db",
		"-asasa.db",
		"adasdas-.db",
	}

	for _, dbname := range shouldPass {
		st := openDatabase(t, dbname)
		if st != nil {
			st.Close()
		}

		os.RemoveAll(DataDirTmp + "/" + dbname)
	}

	for _, dbname := range shouldFail {
		openDatabaseFail(t, dbname)
		//os.RemoveAll(DataDirTmp + "/" + dbname)
	}
}

func TestStoreSetGet(t *testing.T) {
	var (
		err    error
		store  KVStore
		data   []byte
		testDb = "test_set.db"
	)

	store = openDatabase(t, testDb)

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
		if err = store.Set(kv.key, kv.value); err != nil {
			t.Error(err)
		}

		if data, err = store.Get(kv.key); err != nil {
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

	data, err = store.Get([]byte("do not exists"))

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
	var testDb = "testbatch.db"
	store := openDatabase(t, testDb)

	store.StartBatch()

	if store.IsBatch() == false {
		t.Error("StartBatch not setting isBatch = true")
		return
	}

	if err := store.Set([]byte{'a'}, []byte{'b'}); err != nil {
		t.Error(err)
		return
	}

	// should returns nil, nil because the key is in the batch cache
	if data, err := store.Get([]byte{'a'}); err != nil || data != nil {
		t.Error("Key set before wasn't in the write batch cache." +
			" Batch mode isnt working")
	}

	if err := store.FlushBatch(); err != nil {
		t.Error(err)
	}

	if store.IsBatch() == true {
		t.Error("FlushBatch doesnt reset the isBatch")
	}

	store.Close()

	os.RemoveAll(DataDirTmp + "/" + testDb)
}
