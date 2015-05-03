package index

import (
	"fmt"
	"testing"

	"github.com/NeowayLabs/neosearch/lib/neosearch"
)

func getCreateHandler() *CreateIndexHandler {
	cfg := neosearch.NewConfig()
	cfg.Option(neosearch.DataDir("/tmp/"))
	ns := neosearch.New(cfg)

	handler := NewCreateHandler(ns)

	return handler
}

func TestCreateIndexOK(t *testing.T) {
	handler := getCreateHandler()

	defer func() {
		handler.search.Close()
	}()

	for _, name := range []string{
		"test",
		"about",
		"company",
		"people",
		"apple",
		"sucks",
	} {
		body, err := handler.createIndex(name)

		if err != nil {
			t.Error(err)
			continue
		}

		expected := fmt.Sprintf("{\"status\": \"Index '%s' created.\"}", name)

		if string(body) != expected {
			t.Errorf("REST response differs: Received (%s)\nExpected: (%s)",
				string(body), expected)
		}

		deleteIndex(t, handler.search, name)
	}
}

func TestCreateIndexFail(t *testing.T) {
	handler := getCreateHandler()

	defer func() {
		handler.search.Close()
	}()

	for _, name := range []string{
		"_____",
		"87)*()*)",
		"@#$%*()",
		"a",
		"aa",
	} {
		body, err := handler.createIndex(name)

		if err == nil {
			t.Errorf("Invalid index name '%s' should fail", name)
			deleteIndex(t, handler.search, name)
			continue
		}

		if body != nil {
			t.Errorf("JSON response should be nil for index '%s'", name)
			return
		}

		if err.Error() != "Invalid index name" {
			t.Error("Unexpected error")
		}
	}
}
