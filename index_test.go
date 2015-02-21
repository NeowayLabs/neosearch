package neosearch

import (
	"os"
	"testing"
	"time"
)

const DataDirTmp = "/tmp/neosearch-tests"

func init() {
	os.Mkdir(DataDirTmp, 0755)
}

func TestCreateIndex(t *testing.T) {
	cfg := NewConfig()
	cfg.Option(DataDir(DataDirTmp))
	cfg.Option(Debug(false))

	neo := New(cfg)

	shouldPass := []string{
		"test",
		"test2",
	}

	shouldFail := []string{
		"test", // already created index
		"test/",
		"test/kdhakhd",
		"#",
		"a",
		"aa",
		"@",
		"$%¨&*",
	}

	for _, indexName := range shouldPass {
		indexDir := DataDirTmp + "/" + indexName
		_, err := neo.CreateIndex(indexName)

		if err != nil {
			t.Error(err)
			goto cleanup
		}

		if _, err := os.Stat(indexDir); os.IsNotExist(err) {
			t.Errorf("no such file or directory: %s", indexDir)
			goto cleanup
		}
	}

	for _, indexName := range shouldFail {
		_, err := neo.CreateIndex(indexName)

		if err == nil {
			t.Error("Should FAIL because index already exists OR invalid name")
			goto cleanup
		}
	}

cleanup:
	neo.Close()

	for _, indexName := range shouldPass {
		indexDir := DataDirTmp + "/" + indexName
		os.RemoveAll(indexDir)
	}
}

func TestAddDocument(t *testing.T) {
	var (
		data       []byte
		filterData []string
		indexName  = "document-sample"
		indexDir   = DataDirTmp + "/" + indexName
	)

	cfg := NewConfig()

	cfg.Option(DataDir(DataDirTmp))
	cfg.Option(Debug(false))

	neo := New(cfg)

	index, err := neo.CreateIndex(indexName)

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	if _, err := os.Stat(indexDir); os.IsNotExist(err) {
		t.Errorf("no such file or directory: %s", indexDir)
		goto cleanup
	}

	err = index.Add(1, []byte(`{"id": 1, "name": "Neoway Business Solution"}`))

	if err != nil {
		t.Error(err.Error())
		goto cleanup
	}

	if _, err := os.Stat(indexDir + "/document.db"); os.IsNotExist(err) {
		t.Errorf("no such file or directory: %s", indexDir+"/document.db")
		goto cleanup
	}

	err = index.Add(2, []byte(`{"id": 2, "name": "Google Inc."}`))

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	err = index.Add(3, []byte(`{"id": 3, "name": "Facebook Company"}`))

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	err = index.Add(4, []byte(`{"id": 4, "name": "Neoway Teste"}`))

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	data, err = index.Get(1)

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	if string(data) != `{"id": 1, "name": "Neoway Business Solution"}` {
		t.Errorf("Failed to retrieve indexed document")
		goto cleanup
	}

	filterData, err = index.FilterTerm([]byte("name"), []byte("neoway business solution"))

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	if len(filterData) != 1 ||
		filterData[0] != `{"id": 1, "name": "Neoway Business Solution"}` {
		t.Errorf("Failed to filter by field name: %v != %s", filterData, `{"id": 1, "name": "Neoway Business Solution"}`)
		goto cleanup
	}

	filterData, err = index.FilterTerm([]byte("name"), []byte("neoway"))

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	if len(filterData) != 2 || filterData[0] != `{"id": 1, "name": "Neoway Business Solution"}` ||
		filterData[1] != `{"id": 4, "name": "Neoway Teste"}` {
		t.Errorf("Failed to filter by field name: %s != %s", filterData, `[{"id": 1, "name": "Neoway Business Solution"} {"id": 4, "name": "Neoway Teste"}]`)
		goto cleanup
	}

cleanup:
	neo.Close()
	os.RemoveAll(indexDir)
}

func TestPrefixMatch(t *testing.T) {
	var (
		data      []byte
		values    []string
		indexName = "test-prefix"
		indexDir  = DataDirTmp + "/" + indexName
	)

	cfg := NewConfig()
	cfg.Option(DataDir(DataDirTmp))
	cfg.Option(Debug(false))

	neo := New(cfg)

	index, err := neo.CreateIndex(indexName)

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	if _, err := os.Stat(indexDir); os.IsNotExist(err) {
		t.Errorf("no such file or directory: %s", indexDir)
		goto cleanup
	}

	err = index.Add(1, []byte(`{"id": 1, "name": "Neoway Business Solution"}`))

	if err != nil {
		t.Error(err.Error())
		goto cleanup
	}

	if _, err := os.Stat(indexDir + "/document.db"); os.IsNotExist(err) {
		t.Errorf("no such file or directory: %s", indexDir+"/document.db")
		goto cleanup
	}

	err = index.Add(2, []byte(`{"id": 2, "name": "Google Inc."}`))

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	err = index.Add(3, []byte(`{"id": 3, "name": "Facebook Company"}`))

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	err = index.Add(4, []byte(`{"id": 4, "name": "Neoway Teste"}`))

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	data, err = index.Get(1)

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	if string(data) != `{"id": 1, "name": "Neoway Business Solution"}` {
		t.Errorf("Failed to retrieve indexed document")
		goto cleanup
	}

	values, err = index.MatchPrefix([]byte("name"), []byte("neoway"))

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	if len(values) != 2 ||
		values[0] != `{"id": 1, "name": "Neoway Business Solution"}` ||
		values[1] != `{"id": 4, "name": "Neoway Teste"}` {
		t.Error("Failed to retrieve documents with 'name' field prefixed with 'neoway'")
		goto cleanup
	}

cleanup:
	neo.Close()

	if len(neo.Indices) != 0 {
		t.Error("Failed to close all neosearch indices")
	}

	os.RemoveAll(indexDir)
}

func TestBatchAdd(t *testing.T) {
	var (
		indexName = "test-batch"
		indexDir  = DataDirTmp + "/" + indexName
	)

	cfg := NewConfig()
	cfg.Option(DataDir(DataDirTmp))
	cfg.Option(Debug(false))

	neo := New(cfg)

	index, err := neo.CreateIndex(indexName)

	if err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(indexDir); os.IsNotExist(err) {
		t.Errorf("no such file or directory: %s", indexDir)
		return
	}

	index.Batch()

	err = index.Add(1, []byte(`{"id": 1, "name": "Neoway Business Solution"}`))

	if err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(indexDir + "/document.db"); os.IsNotExist(err) {
		t.Errorf("no such file or directory: %s", indexDir+"/document.db")
		return
	}

	err = index.Add(2, []byte(`{"id": 2, "name": "Google Inc."}`))

	if err != nil {
		t.Error(err)
	}

	err = index.Add(3, []byte(`{"id": 3, "name": "Facebook Company"}`))

	if err != nil {
		t.Error(err)
	}

	err = index.Add(4, []byte(`{"id": 4, "name": "Neoway Teste"}`))

	if err != nil {
		t.Error(err)
	}

	data, err := index.Get(1)

	if err != nil {
		t.Error(err)
	}

	if string(data) == `{"id": 1, "name": "Neoway Business Solution"}` {
		t.Errorf("Failed!!! Batch mode doesnt working")
	}

	index.FlushBatch()

	batchWork := false

	for i := 0; i < 3; i++ {
		data, err := index.Get(1)

		if err != nil {
			t.Error(err)
		}

		if string(data) == `{"id": 1, "name": "Neoway Business Solution"}` {
			batchWork = true
			break
		}

		time.Sleep(time.Second * 3)
	}

	if !batchWork {
		t.Error("Failed to execute batch commands")
	}

	neo.Close()

	if len(neo.Indices) != 0 {
		t.Error("Failed to close all neosearch indices")
	}

	os.RemoveAll(indexDir)
}
