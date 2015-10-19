package index

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NeowayLabs/neosearch/lib/neosearch"
	"github.com/NeowayLabs/neosearch/lib/neosearch/config"
	"github.com/julienschmidt/httprouter"
)

func getSearchHandler() *SearchHandler {
	cfg := config.NewConfig()
	cfg.Option(config.DataDir(dataDirTmp))
	ns := neosearch.New(cfg)

	handler := NewSearchHandler(ns)
	return handler
}

func addDocumentsForSearch(indexName string) (*SearchHandler, error) {
	handler := getSearchHandler()

	ind, err := handler.search.CreateIndex(indexName)

	if err != nil {
		return nil, err
	}

	for i, doc := range []string{
		`{"id": 0, "name": "Neoway Business Solution"}`,
		`{"id": 1, "name": "Facebook Inc"}`,
		`{"id": 2, "name": "Google Inc"}`,
	} {

		err = ind.Add(uint64(i), []byte(doc), nil)

		if err != nil {
			return nil, err
		}
	}

	return handler, nil
}

func TestSimpleSearch(t *testing.T) {
	handler, err := addDocumentsForSearch("search-simple")

	if err != nil {
		t.Error(err)
		return
	}

	router := httprouter.New()

	router.Handle("POST", "/:index", handler.ServeHTTP)

	ts := httptest.NewServer(router)

	defer func() {
		handler.search.DeleteIndex("search-simple")
		ts.Close()
		handler.search.Close()
	}()

	searchURL := ts.URL + "/search-simple"

	dsl := `
        {
            "from": 0,
            "size": 10,
            "query": {
                "$and": [
                    {"name": "neoway"}
                ]
            }
        }`

	req, err := http.NewRequest("POST", searchURL, bytes.NewBufferString(dsl))

	if err != nil {
		t.Error(err)
		return
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		t.Error(err)
		return
	}

	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
		return
	}

	resObj := map[string]interface{}{}

	err = json.Unmarshal(content, &resObj)

	if err != nil {
		t.Error(err)
		t.Errorf("Returned value: %s", string(content))
		return
	}

	if resObj["error"] != nil {
		t.Error(resObj["error"])
		return
	}

	if resObj["total"] == nil {
		t.Error("Invalid results response. Field 'total' is required.")
		return
	}

	total, ok := resObj["total"].(float64)

	if !ok {
		t.Errorf("Total must be an float: %+v", total)
		return
	}

	if int(total) != 1 {
		t.Errorf("Search problem. Returns %d but the correct is %d", total, 1)
		return
	}
}

func TestSimpleANDSearch(t *testing.T) {
	handler, err := addDocumentsForSearch("simple-and-search")

	if err != nil {
		t.Error(err)
		return
	}

	router := httprouter.New()

	router.Handle("POST", "/:index", handler.ServeHTTP)

	ts := httptest.NewServer(router)

	defer func() {
		handler.search.DeleteIndex("simple-and-search")
		ts.Close()
		handler.search.Close()
	}()

	searchURL := ts.URL + "/simple-and-search"

	dsl := `
        {
            "from": 0,
            "size": 10,
            "query": {
                "$and": [
                    {"name": "inc"},
                    {"name": "facebook"}
                ]
            }
        }`

	req, err := http.NewRequest("POST", searchURL, bytes.NewBufferString(dsl))

	if err != nil {
		t.Error(err)
		return
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		t.Error(err)
		return
	}

	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
		return
	}

	resObj := map[string]interface{}{}

	err = json.Unmarshal(content, &resObj)

	if err != nil {
		t.Error(err)
		t.Errorf("Returned value: %s", string(content))
		return
	}

	if resObj["error"] != nil {
		t.Error(resObj["error"])
		return
	}

	if resObj["total"] == nil {
		t.Error("Invalid results response. Field 'total' is required.")
		return
	}

	total, ok := resObj["total"].(float64)

	if !ok {
		t.Errorf("Total must be an float: %+v", total)
		return
	}

	if int(total) != 1 {
		t.Errorf("Search problem. Returns %d but the correct is %d", total, 1)
		return
	}

	results, ok := resObj["results"].([]interface{})

	if !ok || results == nil {
		t.Errorf("No results: %+v", resObj["results"])
		return
	}

	if len(results) != 1 {
		t.Errorf("Results length is invalid: %d", len(results))
		return
	}

	r := results[0].(map[string]interface{})

	if r["name"].(string) != "Facebook Inc" {
		t.Errorf("Invalid result: %+v", r)
	}
}
