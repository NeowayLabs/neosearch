package server

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NeowayLabs/neosearch"
)

func getServer(t *testing.T) (*httptest.Server, *neosearch.NeoSearch, error) {
	config := neosearch.NewConfig()
	config.Option(neosearch.DataDir("/tmp/"))
	search := neosearch.New(config)
	serverConfig := ServerConfig{
		Host: "0.0.0.0",
		Port: 9500,
	}

	srv, err := New(search, &serverConfig)

	if err != nil {
		t.Error(err.Error())
		return nil, nil, err
	}

	ts := httptest.NewServer(srv.GetRoutes())

	return ts, search, nil
}

func addDocs(t *testing.T, ts *httptest.Server) {
	indexURL := ts.URL + "/company"

	req, err := http.NewRequest("PUT", indexURL, bytes.NewBufferString(""))

	if err != nil {
		t.Error(err.Error())
		return
	}

	client := &http.Client{}
	_, err = client.Do(req)

	if err != nil {
		t.Error(err.Error())
		return
	}

	req, err = http.NewRequest("POST", indexURL+"/1", bytes.NewBufferString(`{"id": 1, "name": "neoway"}`))

	if err != nil {
		t.Error(err.Error())
		return
	}

	client = &http.Client{}
	_, err = client.Do(req)

	if err != nil {
		t.Error(err.Error())
		return
	}

	req, err = http.NewRequest("POST", indexURL+"/2", bytes.NewBufferString(`{"id": 2, "name": "facebook"}`))

	if err != nil {
		t.Error(err.Error())
		return
	}

	client = &http.Client{}
	_, err = client.Do(req)

	if err != nil {
		t.Error(err.Error())
		return
	}

	req, err = http.NewRequest("POST", indexURL+"/3", bytes.NewBufferString(`{"id": 3, "name": "google"}`))

	if err != nil {
		t.Error(err.Error())
		return
	}

	client = &http.Client{}
	_, err = client.Do(req)

	if err != nil {
		t.Error(err.Error())
		return
	}

}

func deleteIndex(t *testing.T, search *neosearch.NeoSearch, name string) {
	err := search.DeleteIndex(name)

	if err != nil {
		t.Error(err)
	}
}

func TestRESTCreateIndex(t *testing.T) {
	ts, search, _ := getServer(t)
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
	ts, _, _ := getServer(t)
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

func TestRESTGetDocuments(t *testing.T) {
	ts, search, _ := getServer(t)
	defer ts.Close()

	addDocs(t, ts)

	indexURL := ts.URL + "/company"

	res, err := http.Get(indexURL + "/1")

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
		return
	}

	if resObj["error"] != nil {
		t.Error(resObj["error"])
		return
	}

	if resObj["name"] != "neoway" {
		t.Error("Invalid document")
	}

	deleteIndex(t, search, "company")

}
