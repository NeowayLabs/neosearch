package index

import (
	"encoding/json"
	"net/http"

	"github.com/NeowayLabs/neosearch"
	"github.com/NeowayLabs/neosearch/index"
	"github.com/NeowayLabs/neosearch/neosearch/handler"
)

type IndexHandler struct {
	handler.DefaultHandler

	search *neosearch.NeoSearch
}

func New(search *neosearch.NeoSearch) *IndexHandler {
	handler := IndexHandler{}
	handler.search = search

	return &handler
}

func (handler *IndexHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	handler.ProcessVars(req)
	indexName := handler.GetIndexName()

	if indexName == "" {
		handler.Error(res, "no index supplied")
		return
	} else if !index.ValidateIndexName(indexName) {
		handler.Error(res, "Invalid index name: "+indexName)
		return
	}

	index, err := handler.search.OpenIndex(indexName)

	if err != nil {
		handler.Error(res, err.Error())
		return
	}

	body, err := json.Marshal(&index)

	if err != nil {
		handler.Error(res, err.Error())
		return
	}

	res.Write(body)
}
