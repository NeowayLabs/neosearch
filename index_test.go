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

	index, err := neo.CreateIndex("test")

	if err != nil {
		t.Error(err)
	}

	err = index.Add(1, []byte(`{"id": 1, "name": "Neoway Business Solution"}`))

	if err != nil {
		t.Error(err)
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
