package index

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NeowayLabs/neosearch/lib/neosearch"
	"github.com/NeowayLabs/neosearch/lib/neosearch/config"
	"github.com/julienschmidt/httprouter"
)

func getDeleteHandler() *DeleteIndexHandler {
	cfg := config.NewConfig()
	cfg.Option(config.DataDir(dataDirTmp))
	ns := neosearch.New(cfg)

	handler := NewDeleteHandler(ns)

	return handler

}

func TestDeleteServeHTTP_OK(t *testing.T) {
	handler := getDeleteHandler()

	router := httprouter.New()

	router.Handle("DELETE", "/:index", handler.ServeHTTP)

	ts := httptest.NewServer(router)

	defer func() {
		ts.Close()
		handler.search.Close()
	}()

	for _, name := range []string{
		"test-delete-serve-http",
		"delete-this-index",
		"lsdfjlsjflsdjfl",
		"LOL",
	} {

		_, err := handler.search.CreateIndex(name)

		if err != nil {
			t.Error(err)
			return
		}

		deleteURL := ts.URL + "/" + name

		req, err := http.NewRequest("DELETE", deleteURL, bytes.NewBufferString(""))

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
			t.Error("Failed to delete index")
			return
		}

		status := resObj["status"]

		if status != "Index '"+name+"' deleted." {
			t.Errorf("Failed to delete index: %s", status)
			return
		}

		_, err = handler.search.OpenIndex(name)

		if err == nil {
			t.Error(errors.New("Index '" + name + "' should not exist."))
			return
		}
	}
}

func TestDeleteIndex(t *testing.T) {
	handler := getDeleteHandler()
	_, err := handler.search.CreateIndex("test-delete")

	if err != nil {
		t.Error(err)
		return
	}

	defer func() {
		handler.search.Close()
	}()

	err = handler.deleteIndex("test-delete")

	if err != nil {
		t.Error(err)
		return
	}

	err = handler.deleteIndex("ldjfklsjfl")

	if err == nil {
		t.Error(errors.New("should fail: index doesn't exist"))
	}
}
