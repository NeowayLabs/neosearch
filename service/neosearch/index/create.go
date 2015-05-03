package index

import (
	"fmt"
	"net/http"

	"github.com/NeowayLabs/neosearch/lib/neosearch"
	"github.com/NeowayLabs/neosearch/service/neosearch/handler"
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

	body, err := handler.createIndex(indexName)

	if err != nil {
		handler.Error(res, string(body))
		return
	}

	handler.WriteJSON(res, body)
}

func (handler *CreateIndexHandler) createIndex(name string) ([]byte, error) {
	_, err := handler.search.CreateIndex(name)

	if err != nil {
		return nil, err
	}

	response := []byte(fmt.Sprintf("{\"status\": \"Index '%s' created.\"}", name))
	return response, nil
}
