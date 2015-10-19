package index

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/NeowayLabs/neosearch/lib/neosearch"
	"github.com/NeowayLabs/neosearch/lib/neosearch/config"
	"github.com/julienschmidt/httprouter"
)

func getGetHandler() *GetHandler {
	cfg := config.NewConfig()
	cfg.Option(config.DataDir(dataDirTmp))
	ns := neosearch.New(cfg)

	handler := NewGetHandler(ns)

	return handler
}

func TestGetDocumentsOK(t *testing.T) {
	handler := getGetHandler()

	defer func() {
		handler.search.DeleteIndex("test-get-ok")
		handler.search.Close()
	}()

	ind, err := handler.search.CreateIndex("test-get-ok")

	if err != nil {
		t.Error(err)
		return
	}

	router := httprouter.New()

	router.Handle("GET", "/:index/:id", handler.ServeHTTP)

	ts := httptest.NewServer(router)

	for i, doc := range []string{
		`{"id": 0, "bleh": "test"}`,
		`{"id": 1, "title": "ldjfjl"}`,
		`{"id": 2, "title": "hjdfskhfk"}`,
	} {

		err = ind.Add(uint64(i), []byte(doc), nil)

		if err != nil {
			t.Error(err)
			return
		}

		getURL := ts.URL + "/test-get-ok/" + strconv.Itoa(i)

		req, err := http.NewRequest("GET", getURL, bytes.NewBufferString(""))

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

		if string(content) != doc {
			t.Errorf("Differs: %s != %s", string(content), doc)
			t.Errorf("Failed to get document: %s", string(content))
			return
		}
	}

}
