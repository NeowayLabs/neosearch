package index

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NeowayLabs/neosearch"
	"github.com/gorilla/mux"
)

func getServer() (*httptest.Server, *neosearch.NeoSearch) {
	router := mux.NewRouter()
	config := neosearch.NewConfig()
	config.Option(neosearch.DataDir("/tmp/"))
	search := neosearch.New(config)

	indexHandler := New(search)
	createIndexHandler := NewCreateHandler(search)

	router.Handle("/{index}", indexHandler).
		Methods("GET")

	router.Handle("/{index}", createIndexHandler).
		Methods("PUT")

	ts := httptest.NewServer(router)

	return ts, search
}

func deleteIndex(t *testing.T, search *neosearch.NeoSearch, name string) {
	err := search.DeleteIndex(name)

	if err != nil {
		t.Error(err)
	}
}

func TestRESTCreateIndex(t *testing.T) {
	ts, search := getServer()
	defer ts.Close()

	createURL := ts.URL + "/company"

	req, err := http.NewRequest("PUT", createURL, bytes.NewBufferString(""))

	if err != nil {
		t.Error(err)
		return
	}

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		t.Error(err)
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
		return
	}

	if resObj["error"] != nil {
		t.Error(resObj["error"])
		return
	}

	if resObj["status"] == nil {
		t.Error("Failed to create index")
		return
	}

	status := resObj["status"]

	if status != "Index 'company' created." {
		t.Errorf("Failed to create index: %s", status)
		return
	}

	deleteIndex(t, search, "company")
}

func TestRESTIndexInfo(t *testing.T) {
	ts, _ := getServer()
	defer ts.Close()

	infoURL := ts.URL + "/test"

	res, err := http.Get(infoURL)

	if err != nil {
		log.Fatal(err)
	}
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}

	resObj := map[string]interface{}{}

	err = json.Unmarshal(content, &resObj)

	if err != nil {
		t.Errorf("Failed to unmarshal json response: %s", err.Error())
	}

	if resObj["error"] == nil {
		t.Error("Failed: test index should not exists")
	}

	errMsg := resObj["error"]

	if errMsg.(string) != "Index 'test' not found in directory '/tmp'." {
		t.Error("Wrong error message")
	}
}
