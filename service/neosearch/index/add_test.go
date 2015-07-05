package index

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/NeowayLabs/neosearch/lib/neosearch"
	"github.com/julienschmidt/httprouter"
)

func getAddDocHandler() *AddHandler {
	cfg := neosearch.NewConfig()
	cfg.Option(neosearch.DataDir("/tmp/"))
	ns := neosearch.New(cfg)

	handler := NewAddHandler(ns)

	return handler
}

func TestAddDocumentsOK(t *testing.T) {
	handler := getAddDocHandler()

	defer func() {
		handler.search.DeleteIndex("test-ok")
		handler.search.Close()
	}()

	_, err := handler.search.CreateIndex("test-ok")

	if err != nil {
		t.Error(err)
		return
	}

	for i, doc := range []string{
		`{"doc": {"id": 0, "bleh": "test"}}`,
		`{"doc": {"id": 1, "title": "ldjfjl"}}`,
		`{"doc": {"id": 2, "title": "hjdfskhfk"}}`,
	} {

		err = handler.addDocument("test-ok", uint64(i), []byte(doc))

		if err != nil {
			t.Error(err)
			return
		}
	}
}

func TestAddDocumentsREST_OK(t *testing.T) {
	handler := getAddDocHandler()
	router := httprouter.New()
	router.Handle("POST", "/:index/:id", handler.ServeHTTP)
	ts := httptest.NewServer(router)

	defer func() {
		handler.search.DeleteIndex("test-rest-add-ok")
		ts.Close()
		handler.search.Close()
	}()

	_, err := handler.search.CreateIndex("test-rest-add-ok")

	if err != nil {
		t.Error(err)
		return
	}

	for i, doc := range []string{
		`{"doc": {"id": 0, "bleh": "test"}}`,
		`{"doc": {"id": 1, "title": "ldjfjl"}}`,
		`{"doc": {"id": 2, "title": "hjdfskhfk"}}`,
	} {
		addURL := ts.URL + "/test-rest-add-ok/" + strconv.Itoa(i)

		req, err := http.NewRequest("POST", addURL, bytes.NewBufferString(doc))

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
			t.Errorf("Returned value: %s", string(content))
			return
		}

		if resObj["error"] != nil {
			t.Error(resObj["error"])
			return
		}

		if resObj["status"] == nil {
			t.Error("Failed to add document")
			return
		}

		status := resObj["status"]
		expected := "Document " + strconv.Itoa(i) + " indexed."

		if status != expected {
			t.Errorf("Differs: %s != %s", status, expected)
			t.Errorf("Failed to add document: %s", status)
			return
		}
	}
}

func TestAddDocumentsFail(t *testing.T) {
	handler := getAddDocHandler()

	defer func() {
		handler.search.DeleteIndex("test-fail")
		handler.search.Close()
	}()

	_, err := handler.search.CreateIndex("test-fail")

	if err != nil {
		t.Error(err)
		return
	}

	for i, doc := range []string{
		`{}`,
		`{"metadata": {}}`,
		`{"t": "sçdçs"}`,
		``,
		`test`,
		`   `,
		`[]`,
		`[{}]`,
	} {

		err = handler.addDocument("test-fail", uint64(i), []byte(doc))

		if err == nil {
			t.Error(fmt.Errorf("Invalid document: %s", doc))
			return
		}
	}
}
