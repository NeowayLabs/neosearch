package index

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NeowayLabs/neosearch/lib/neosearch"
	"github.com/gorilla/mux"
)

func getAnalyzeGetHandler() *GetAnalyseHandler {
	cfg := neosearch.NewConfig()
	cfg.Option(neosearch.DataDir("/tmp/"))
	ns := neosearch.New(cfg)

	handler := NewGetAnalyzeHandler(ns)
	return handler
}

func TestGetAnalyze(t *testing.T) {
	handler := getAnalyzeGetHandler()

	router := mux.NewRouter()
	router.Handle("/{index}/{id}/analyze", handler).Methods("GET")
	ts := httptest.NewServer(router)

	defer func() {
		ts.Close()
		handler.search.DeleteIndex("test-analyze-ok")
		handler.search.Close()
	}()

	_, err := handler.search.CreateIndex("test-analyze-ok")

	if err != nil {
		t.Error(err)
		return
	}

	for _, testPair := range []struct {
		id  string
		out string
	}{
		{"1", `USING test-analyze-ok.document.db GET uint(1);`},
	} {
		id := testPair.id
		out := testPair.out

		analyzeURL := ts.URL + "/test-analyze-ok/" + id + "/analyze"

		req, err := http.NewRequest("GET", analyzeURL, bytes.NewBufferString(""))

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

		if err == nil {
			t.Error(errors.New("should return a neosearch-cli commands: " + string(content)))
			return
		}

		if string(content) != out {
			t.Error(fmt.Errorf("analyze differs: (%s) != (%s)", string(content), out))
			return
		}
	}
}
