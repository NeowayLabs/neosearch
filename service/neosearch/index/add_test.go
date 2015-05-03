package index

import (
	"fmt"
	"testing"

	"github.com/NeowayLabs/neosearch/lib/neosearch"
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
		`{"id": 0, "bleh": "test"}`,
		`{"id": 1, "title": "ldjfjl"}`,
		`{"id": 2, "title": "hjdfskhfk"}`,
	} {

		err = handler.addDocument("test-ok", uint64(i), []byte(doc))

		if err != nil {
			t.Error(err)
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
