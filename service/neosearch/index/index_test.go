package index

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/NeowayLabs/neosearch/lib/neosearch"
	"github.com/NeowayLabs/neosearch/lib/neosearch/config"
	"github.com/NeowayLabs/neosearch/lib/neosearch/index"
)

var dataDirTmp string

func init() {
	var err error
	dataDirTmp, err = ioutil.TempDir("/tmp", "neosearch-service-index-")
	if err != nil {
		panic(err)
	}
}

func getIndexHandler() *IndexHandler {
	cfg := config.NewConfig()
	cfg.Option(config.DataDir(dataDirTmp))
	ns := neosearch.New(cfg)

	handler := New(ns)

	return handler
}

func deleteIndex(t *testing.T, search *neosearch.NeoSearch, name string) {
	err := search.DeleteIndex(name)

	if err != nil {
		t.Error(err)
	}
}

func TestIndexNotExist(t *testing.T) {
	handler := getIndexHandler()

	for _, name := range []string{
		"test",
		"info",
		"lsajldkjal",
		"__",
		"about",
		"hack",
	} {
		_, err := handler.serveIndex(name)

		if err == nil {
			t.Error(fmt.Errorf("Index '%s' shall not exist", name))
			return
		}
	}
}

func addDocs(t *testing.T, index *index.Index) {
	err := index.Add(1, []byte(`{"title": "teste"}`), nil)

	if err != nil {
		t.Error(err)
		return
	}
}

func TestIndexInfo(t *testing.T) {
	handler := getIndexHandler()

	defer func() {
		deleteIndex(t, handler.search, "test-index-info")
		handler.search.Close()
	}()

	index, err := handler.search.CreateIndex("test-index-info")

	addDocs(t, index)

	body, err := handler.serveIndex("test-index-info")

	if err != nil {
		t.Error(err)
		return
	}

	if string(body) != `{"name":"test-index-info"}` {
		t.Errorf("Invalid index info: %s", string(body))
		return
	}
}
