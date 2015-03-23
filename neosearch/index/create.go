package index

import (
	"fmt"
	"net/http"

	"github.com/NeowayLabs/neosearch"
	"github.com/NeowayLabs/neosearch/neosearch/handler"
)

type CreateIndexHandler struct {
	handler.DefaultHandler
	search *neosearch.NeoSearch
}

func NewCreateHandler(search *neosearch.NeoSearch) *CreateIndexHandler {
	return &CreateIndexHandler{
		search: search,
	}
}

func (handler *CreateIndexHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	handler.ProcessVars(req)
	indexName := handler.GetIndexName()

	if exists, err := handler.search.IndexExists(indexName); exists == true && err == nil {
		response := map[string]string{
			"error": "Index '" + indexName + "' already exists.",
		}

		handler.WriteJSONObject(res, response)
		return
	} else if exists == false && err != nil {
		handler.Error(res, err.Error())
		return
	}

	_, err := handler.search.CreateIndex(indexName)

	if err != nil {
		handler.Error(res, err.Error())
		return
	}

	handler.WriteJSON(res, []byte(fmt.Sprintf("{\"status\": \"Index '%s' created.\"}", indexName)))
}
