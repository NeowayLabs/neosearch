package neosearch

import (
	"os"
	"testing"
	"time"
)

func TestCreateIndex(t *testing.T) {
	dataDir := "/tmp/neosearch-test"

	os.Mkdir(dataDir, 0755)

	cfg := NewConfig()
	cfg.Option(DataDir(dataDir))
	cfg.Option(Debug(false))

	neo := New(cfg)

	_, err := neo.CreateIndex("test")

	if err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(dataDir + "/test"); os.IsNotExist(err) {
		t.Errorf("no such file or directory: %s", dataDir+"/test")
		return
	}

	_, err = neo.CreateIndex("test")

	if err == nil {
		t.Error("Should FAIL because index already exists")
	}

	_, err = neo.CreateIndex("test2")

	if err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(dataDir + "/test2"); os.IsNotExist(err) {
		t.Errorf("no such file or directory: %s", dataDir+"/test2")
		return
	}

	_, err = neo.CreateIndex("test//")

	if err == nil {
		t.Error("Has invalid name, should fail")
	}

	_, err = neo.CreateIndex("#")

	if err == nil {
		t.Error("Has invalid name, should fail")
	}

	_, err = neo.CreateIndex("-")

	if err == nil {
		t.Error("Has invalid name, should fail")
	}

	_, err = neo.CreateIndex("a")

	if err == nil {
		t.Error("Has invalid name, should fail")
	}

	_, err = neo.CreateIndex("aa")

	if err == nil {
		t.Error("Has invalid name, should fail")
	}

	_, err = neo.CreateIndex("test-$")

	if err == nil {
		t.Error("Has invalid name, should fail")
	}

	neo.Close()
	os.RemoveAll(dataDir)
}

func TestAddDocument(t *testing.T) {
	dataDir := "/tmp/neosearch-test"

	os.Mkdir(dataDir, 0755)

	cfg := NewConfig()

	cfg.Option(DataDir(dataDir))
	cfg.Option(Debug(false))

	neo := New(cfg)

	index, err := neo.CreateIndex("test")

	if err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(dataDir + "/test"); os.IsNotExist(err) {
		t.Errorf("no such file or directory: %s", dataDir+"/test")
		return
	}

	err = index.Add(1, []byte(`{"id": 1, "name": "Neoway Business Solution"}`))

	if err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(dataDir + "/test/document.db"); os.IsNotExist(err) {
		t.Errorf("no such file or directory: %s", dataDir+"/test/document.db")
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

	if string(data) != `{"id": 1, "name": "Neoway Business Solution"}` {
		t.Errorf("Failed to retrieve indexed document")
	}

	filterData, err := index.FilterTerm([]byte("name"), []byte("neoway business solution"))

	if err != nil {
		t.Error(err)
	}

	if len(filterData) != 1 ||
		filterData[0] != `{"id": 1, "name": "Neoway Business Solution"}` {
		t.Errorf("Failed to filter by field name: %v != %s", filterData, `{"id": 1, "name": "Neoway Business Solution"}`)
	}

	filterData, err = index.FilterTerm([]byte("name"), []byte("neoway"))

	if err != nil {
		t.Error(err)
	}

	if len(filterData) != 2 || filterData[0] != `{"id": 1, "name": "Neoway Business Solution"}` ||
		filterData[1] != `{"id": 4, "name": "Neoway Teste"}` {
		t.Errorf("Failed to filter by field name: %s != %s", filterData, `[{"id": 1, "name": "Neoway Business Solution"} {"id": 4, "name": "Neoway Teste"}]`)
	}

	neo.Close()
	os.RemoveAll(dataDir)
}

func TestPrefixMatch(t *testing.T) {
	dataDir := "/tmp/neosearch-test"

	os.Mkdir(dataDir, 0755)

	cfg := NewConfig()
	cfg.Option(DataDir(dataDir))
	cfg.Option(Debug(false))

	neo := New(cfg)

	index, err := neo.CreateIndex("test")

	if err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(dataDir + "/test"); os.IsNotExist(err) {
		t.Errorf("no such file or directory: %s", dataDir+"/test")
		return
	}

	err = index.Add(1, []byte(`{"id": 1, "name": "Neoway Business Solution"}`))

	if err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(dataDir + "/test/document.db"); os.IsNotExist(err) {
		t.Errorf("no such file or directory: %s", dataDir+"/test/document.db")
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

	if string(data) != `{"id": 1, "name": "Neoway Business Solution"}` {
		t.Errorf("Failed to retrieve indexed document")
	}

	values, err := index.MatchPrefix([]byte("name"), []byte("neoway"))

	if err != nil {
		t.Error(err)
	}

	if len(values) != 2 ||
		values[0] != `{"id": 1, "name": "Neoway Business Solution"}` ||
		values[1] != `{"id": 4, "name": "Neoway Teste"}` {
		t.Error("Failed to retrieve documents with 'name' field prefixed with 'neoway'")
	}

	neo.Close()

	if len(neo.Indices) != 0 {
		t.Error("Failed to close all neosearch indices")
	}

	os.RemoveAll(dataDir)
}

func TestBatchAdd(t *testing.T) {
	dataDir := "/tmp/neosearch-test"

	os.Mkdir(dataDir, 0755)

	cfg := NewConfig()
	cfg.Option(DataDir(dataDir))
	cfg.Option(Debug(false))

	neo := New(cfg)

	index, err := neo.CreateIndex("test")

	if err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(dataDir + "/test"); os.IsNotExist(err) {
		t.Errorf("no such file or directory: %s", dataDir+"/test")
		return
	}

	index.Batch()

	err = index.Add(1, []byte(`{"id": 1, "name": "Neoway Business Solution"}`))

	if err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(dataDir + "/test/document.db"); os.IsNotExist(err) {
		t.Errorf("no such file or directory: %s", dataDir+"/test/document.db")
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

	os.RemoveAll(dataDir)
}
