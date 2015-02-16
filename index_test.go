package neosearch

import (
	"fmt"
	"os"
	"testing"
)

func TestCreateIndex(t *testing.T) {
	dataDir := "/tmp/neosearch-test"

	os.Mkdir(dataDir, 0755)

	neo := New(Config{
		DataDir: dataDir,
		Debug:   false,
	})

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

	neo := New(Config{
		DataDir: dataDir,
		Debug:   false,
	})

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

	if string(filterData) != "[1]" {
		t.Errorf("Failed to filter by field name: %v != %s", filterData, "[1]")
	}

	filterData, err = index.FilterTerm([]byte("name"), []byte("neoway"))

	if err != nil {
		t.Error(err)
	}

	fmt.Println("term out: ", string(filterData))

	if string(filterData) != "[1,4]" {
		t.Errorf("Failed to filter by field name: %s != %s", filterData, "[1,4]")
	}

	neo.Close()
	os.RemoveAll(dataDir)
}

func TestPrefixMatch(t *testing.T) {
	dataDir := "/tmp/neosearch-test"

	os.Mkdir(dataDir, 0755)

	neo := New(Config{
		DataDir: dataDir,
		Debug:   false,
	})

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

	fmt.Println("Found at: ", values)

	neo.Close()
	os.RemoveAll(dataDir)
}
