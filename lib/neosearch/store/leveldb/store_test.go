// +build leveldb

package leveldb

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/NeowayLabs/neosearch/lib/neosearch/store"
	"github.com/NeowayLabs/neosearch/lib/neosearch/store/test"
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

func TestStoreHasBackend2(t *testing.T) {
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

func TestOpenDatabase2(t *testing.T) {
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

func TestStoreSetGet2(t *testing.T) {
	var (
		kv     store.KVStore
		testDb = "test_set.db"
	)

	os.Mkdir(DataDirTmp+string(filepath.Separator)+"sample-store-set-get", 0755)
	if kv = openDatabase(t, "sample-store-set-get", testDb); kv == nil {
		return
	}

	test.CommonTestStoreSetGet(t, kv)

	kv.Close()
	os.RemoveAll(DataDirTmp + "/" + testDb)
}

func TestBatchWrite2(t *testing.T) {
	var (
		kv     store.KVStore
		testDb = "testbatch.db"
	)

	os.Mkdir(DataDirTmp+string(filepath.Separator)+"sample-batch-write", 0755)
	if kv = openDatabase(t, "sample-batch-write", testDb); kv == nil {
		return
	}

	test.CommonTestBatchWrite(t, kv)

	kv.Close()
	os.RemoveAll(DataDirTmp + "/" + testDb)
}

func TestBatchMultiWrite2(t *testing.T) {
	var (
		kv     store.KVStore
		testDb = "test_set-multi.db"
	)

	os.Mkdir(DataDirTmp+string(filepath.Separator)+"sample-batch-multi-write", 0755)
	if kv = openDatabase(t, "sample-batch-multi-write", testDb); kv == nil {
		return
	}

	test.CommonTestBatchMultiWrite(t, kv)

	kv.Close()
	os.RemoveAll(DataDirTmp + "/" + testDb)
}

func TestStoreMergeSet2(t *testing.T) {
	var (
		kv     store.KVStore
		testDb = "test_mergeset.db"
	)

	os.Mkdir(DataDirTmp+string(filepath.Separator)+"sample-store-merge-set", 0755)
	if kv = openDatabase(t, "sample-store-merge-set", testDb); kv == nil {
		return
	}

	test.CommonTestStoreMergeSet(t, kv)

	kv.Close()
	os.RemoveAll(DataDirTmp + "/" + testDb)
}

func TestStoreIterator2(t *testing.T) {
	var (
		kv     store.KVStore
		testDb = "test_iterator.db"
	)

	os.Mkdir(DataDirTmp+string(filepath.Separator)+"sample-store-iterator", 0755)
	if kv = openDatabase(t, "sample-store-iterator", testDb); kv == nil {
		return
	}

	test.CommonTestStoreIterator(t, kv)

	kv.Close()
	os.RemoveAll(DataDirTmp + "/" + testDb)
}
