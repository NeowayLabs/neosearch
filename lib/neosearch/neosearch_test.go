package neosearch

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/NeowayLabs/neosearch/lib/neosearch/config"
	"github.com/NeowayLabs/neosearch/lib/neosearch/index"
)

var DataDirTmp string

func init() {
	var err error
	DataDirTmp, err = ioutil.TempDir("/tmp", "neosearch-index-")

	if err != nil {
		panic(err)
	}
}

func TestCreateIndex(t *testing.T) {
	cfg := config.NewConfig()
	cfg.Option(config.DataDir(DataDirTmp))

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

func TestOpenIndexCache(t *testing.T) {
	var (
		err  error
		indx *index.Index
	)

	cfg := config.NewConfig()
	cfg.Option(config.DataDir(DataDirTmp))

	neo := New(cfg)

	_, err = neo.CreateIndex("test-cache")

	if err != nil {
		t.Error(err)
	}

	if neo.GetIndices().Len() != 1 {
		t.Error("Cache problem")
		goto cleanup
	}

	indx, err = neo.OpenIndex("test-cache")

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	if indx == nil {
		t.Error("Failed to open index")
		goto cleanup
	}

	if neo.GetIndices().Len() != 1 {
		t.Error("Cache problem")
		goto cleanup
	}

	indx, err = neo.OpenIndex("test-cache")

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	if indx == nil {
		t.Error("Failed to open index")
		goto cleanup
	}

	if neo.GetIndices().Len() != 1 {
		t.Error("Cache problem")
		goto cleanup
	}

cleanup:
	neo.Close()
	neo.DeleteIndex("test-cache")

}

func TestDeleteIndex(t *testing.T) {
	cfg := config.NewConfig()
	cfg.Option(config.DataDir(DataDirTmp))

	neo := New(cfg)

	err := neo.DeleteIndex("lsdlas")

	if err == nil {
		t.Error("Failed: Index does not exists yet")
	}

	_, err = neo.CreateIndex("test")

	if err != nil {
		t.Error(err)
	}

	err = neo.DeleteIndex("test")

	if err != nil {
		t.Error(err)
	}
}

func TestAddDocument(t *testing.T) {
	var (
		data       []byte
		filterData []string
		indexName  = "document-sample"
		indexDir   = DataDirTmp + "/" + indexName
		total      uint64
	)

	cfg := config.NewConfig()
	cfg.Option(config.DataDir(DataDirTmp))

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

	err = index.Add(1, []byte(`{"id": 1, "name": "Neoway Business Solution"}`), nil)

	if err != nil {
		t.Error(err.Error())
		goto cleanup
	}

	if _, err := os.Stat(indexDir + "/document.db"); os.IsNotExist(err) {
		t.Errorf("no such file or directory: %s", indexDir+"/document.db")
		goto cleanup
	}

	err = index.Add(2, []byte(`{"id": 2, "name": "Google Inc."}`), nil)

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	err = index.Add(3, []byte(`{"id": 3, "name": "Facebook Company"}`), nil)

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	err = index.Add(4, []byte(`{"id": 4, "name": "Neoway Teste"}`), nil)

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

	filterData, total, err = index.FilterTerm([]byte("name"), []byte("neoway business solution"), 0)

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	if total != 1 || len(filterData) != 1 ||
		filterData[0] != `{"id": 1, "name": "Neoway Business Solution"}` {
		t.Errorf("Failed to filter by field name: %v != %s", filterData, `{"id": 1, "name": "Neoway Business Solution"}`)
		goto cleanup
	}

	filterData, total, err = index.FilterTerm([]byte("name"), []byte("neoway"), 0)

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	if total != 2 || len(filterData) != 2 || !reflect.DeepEqual(filterData, []string{
		`{"id": 1, "name": "Neoway Business Solution"}`,
		`{"id": 4, "name": "Neoway Teste"}`,
	}) {
		t.Errorf("Failed to filter by field name: %s != %s", filterData, `[{"id": 1, "name": "Neoway Business Solution"} {"id": 4, "name": "Neoway Teste"}]`)
		goto cleanup
	}

cleanup:
	neo.Close()
	os.RemoveAll(indexDir)
}

func BenchmarkAddDocuments(b *testing.B) {
	var (
		indexName = "document-bench-sample"
	)

	dataDirTmp, _ := ioutil.TempDir("/tmp", "neosearch-index-bench-")

	indexDir := dataDirTmp + "/" + indexName

	cfg := config.NewConfig()
	cfg.Option(config.DataDir(dataDirTmp))

	neo := New(cfg)

	index, err := neo.CreateIndex(indexName)

	if err != nil {
		b.Fatal(err)
		goto cleanup
	}

	if _, err := os.Stat(indexDir); os.IsNotExist(err) {
		b.Fatal(fmt.Sprintf("no such file or directory: %s", indexDir))
		goto cleanup
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err = index.Add(1, []byte(`{"id": 1, "name": "Neoway Business Solution"}`), nil)

		if err != nil {
			b.Fatal(err)
			goto cleanup
		}
	}

cleanup:
	neo.Close()
	os.RemoveAll(indexDir)
}

func TestAddDocumentWithObject(t *testing.T) {
	var (
		data       []byte
		filterData []string
		indexName  = "document-object-sample"
		indexDir   = DataDirTmp + "/" + indexName
		total      uint64
	)

	docNeoway := `{
    "id": 1,
    "name": "Neoway Business Solution",
    "address": {
        "city": "Florianópolis",
        "district": "Itacorubi",
        "street": "Patricio Farias",
        "latlon": [
            -27.545198,
            -48.504827
        ]
    }
}`

	docGoogle := `{
    "id": 2,
    "name": "Google Inc.",
    "address": {
        "city": "Mountain View",
        "street": "Amphitheatre Parkway",
        "latlon": [
            37.422541,
            -122.084221
        ]
    }
}`

	docFacebook := `{
    "id": 3,
    "name": "Facebook Company",
    "address": {
        "city": "Menlo Park",
        "street": "Hacker Way",
        "latlon": [
            37.484770,
            -122.147914
        ]
    }
}`

	cfg := config.NewConfig()
	cfg.Option(config.DataDir(DataDirTmp))

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

	err = index.Add(1, []byte(docNeoway), nil)

	if err != nil {
		t.Error(err.Error())
		goto cleanup
	}

	if _, err := os.Stat(indexDir + "/document.db"); os.IsNotExist(err) {
		t.Errorf("no such file or directory: %s", indexDir+"/document.db")
		goto cleanup
	}

	err = index.Add(2, []byte(docGoogle), nil)

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	err = index.Add(3, []byte(docFacebook), nil)

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	data, err = index.Get(1)

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	if string(data) != docNeoway {
		t.Errorf("Failed to retrieve indexed document")
		goto cleanup
	}

	filterData, total, err = index.FilterTerm([]byte("name"), []byte("neoway business solution"), 0)

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	if total != 1 || len(filterData) != 1 ||
		filterData[0] != docNeoway {
		t.Errorf("Failed to filter by field name: %v != %s", filterData, docNeoway)
		goto cleanup
	}

	filterData, total, err = index.FilterTerm([]byte("name"), []byte("neoway"), 0)

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	if total != 1 || len(filterData) != 1 || !reflect.DeepEqual(filterData, []string{
		docNeoway,
	}) {
		t.Errorf("Failed to filter by field name: %s != %s", filterData, `[`+docNeoway+`]`)
		goto cleanup
	}

	filterData, total, err = index.FilterTerm([]byte("address.city"), []byte("menlo"), 0)

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	if total != 1 || len(filterData) != 1 || !reflect.DeepEqual(filterData, []string{
		docFacebook,
	}) {
		t.Errorf("Failed to filter by field name: %s != %s", filterData, `[`+docFacebook+`]`)
		goto cleanup
	}

	filterData, total, err = index.FilterTerm([]byte("address.street"), []byte("hacker"), 0)

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	if total != 1 || len(filterData) != 1 || !reflect.DeepEqual(filterData, []string{
		docFacebook,
	}) {
		t.Errorf("Failed to filter by field name: %s != %s", filterData, `[`+docFacebook+`]`)
		goto cleanup
	}

cleanup:
	neo.Close()
	os.RemoveAll(indexDir)
}

func TestComplexDocument(t *testing.T) {
	var (
		data       []byte
		filterData []string
		indexName  = "document-complex-object"
		indexDir   = DataDirTmp + "/" + indexName
		total      uint64
	)

	metadata := index.Metadata{
		"creationDate": index.Metadata{
			"type":   "date",
			"format": "Jan 2, 2006 at 3:04pm (MST)",
		},
	}

	document1 := `{
    "cnpj": "01458782000170",
    "companyName": "Some company name",
    "creationDate": "Nov 10, 2009 at 3:00pm (UTC)",
    "status": {
        "date": "Jan 01, 2015 at 9:00pm (UTC)",
        "description": "ACTIVE",
        "especial": {}
    },
    "address": {
        "district": "some place",
        "city": "some city",
        "state": "AA",
        "zipcode": "52211",
        "street": "some street",
        "number": 23132
    },
    "info": {
        "hasEmail": false,
        "hasTel": false,
        "hasDomain": false
    },
    "employee": [
        {
            "name": "John doe",
            "age": 25
        },
        {
            "name": "Mary Doe",
            "age": 22
        }
    ]
}`

	document2 := `{
    "cnpj": "01458782000170",
    "companyName": "Another company",
    "creationDate": "Nov 10, 2011 at 3:00pm (UTC)",
    "status": {
        "date": "Jan 10, 2015 at 9:00pm (UTC)",
        "description": "ACTIVE",
        "especial": {}
    },
    "address": {
        "district": "",
        "city": "other city",
        "state": "BB",
        "zipcode": "999",
        "street": "other street",
        "number": 111
    },
    "info": {
        "hasEmail": true,
        "hasTel": true,
        "hasDomain": false
    },
    "employee": [
        {
            "name": "John Snow",
            "age": 33
        }
    ]
}`

	cfg := config.NewConfig()
	cfg.Option(config.DataDir(DataDirTmp))

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

	err = index.Add(1, []byte(document1), metadata)

	if err != nil {
		t.Error(err.Error())
		goto cleanup
	}

	if _, err := os.Stat(indexDir + "/document.db"); os.IsNotExist(err) {
		t.Errorf("no such file or directory: %s", indexDir+"/document.db")
		goto cleanup
	}

	err = index.Add(2, []byte(document2), metadata)

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	data, err = index.Get(1)

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	if string(data) != document1 {
		t.Errorf("Failed to retrieve indexed document")
		goto cleanup
	}

	filterData, total, err = index.FilterTerm([]byte("companyName"), []byte("another"), 0)

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	if total != 1 || len(filterData) != 1 ||
		filterData[0] != document2 {
		t.Errorf("Failed to filter by field name: %v != %s", filterData, document2)
		goto cleanup
	}

	filterData, total, err = index.FilterTerm([]byte("companyName"), []byte("some"), 0)

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	if total != 1 || len(filterData) != 1 || !reflect.DeepEqual(filterData, []string{
		document1,
	}) {
		t.Errorf("Failed to filter by field name: %s != %s", filterData, `[`+document1+`]`)
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

	cfg := config.NewConfig()
	cfg.Option(config.DataDir(DataDirTmp))

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

	err = index.Add(1, []byte(`{"id": 1, "name": "Neoway Business Solution"}`), nil)

	if err != nil {
		t.Error(err.Error())
		goto cleanup
	}

	if _, err := os.Stat(indexDir + "/document.db"); os.IsNotExist(err) {
		t.Errorf("no such file or directory: %s", indexDir+"/document.db")
		goto cleanup
	}

	err = index.Add(2, []byte(`{"id": 2, "name": "Google Inc."}`), nil)

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	err = index.Add(3, []byte(`{"id": 3, "name": "Facebook Company"}`), nil)

	if err != nil {
		t.Error(err)
		goto cleanup
	}

	err = index.Add(4, []byte(`{"id": 4, "name": "Neoway Teste"}`), nil)

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

	if neo.GetIndices().Len() != 0 {
		t.Error("Failed to close all neosearch indices")
	}

	os.RemoveAll(indexDir)
}

func TestBatchAdd(t *testing.T) {
	var (
		indexName = "test-batch"
		indexDir  = DataDirTmp + "/" + indexName
	)

	cfg := config.NewConfig()
	cfg.Option(config.DataDir(DataDirTmp))

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

	err = index.Add(1, []byte(`{"id": 1, "name": "Neoway Business Solution"}`), nil)

	if err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(indexDir + "/document.db"); os.IsNotExist(err) {
		t.Errorf("no such file or directory: %s", indexDir+"/document.db")
		return
	}

	err = index.Add(2, []byte(`{"id": 2, "name": "Google Inc."}`), nil)

	if err != nil {
		t.Error(err)
	}

	err = index.Add(3, []byte(`{"id": 3, "name": "Facebook Company"}`), nil)

	if err != nil {
		t.Error(err)
	}

	err = index.Add(4, []byte(`{"id": 4, "name": "Neoway Teste"}`), nil)

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

	if neo.GetIndices().Len() != 0 {
		t.Error("Failed to close all neosearch indices")
	}

	os.RemoveAll(indexDir)
}
