package index

import (
	"net/http"
	"strconv"

	"github.com/NeowayLabs/neosearch"
	"github.com/NeowayLabs/neosearch/neosearch/handler"
)

type GetHandler struct {
	handler.DefaultHandler
	search *neosearch.NeoSearch
}

func NewGetHandler(search *neosearch.NeoSearch) *GetHandler {
	return &GetHandler{
		search: search,
	}
}

func (handler *GetHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var (
		err      error
		document []byte
		exists   bool
	)

	handler.ProcessVars(req)
	indexName := handler.GetIndexName()

	if exists, err = handler.search.IndexExists(indexName); exists != true && err == nil {
		response := map[string]string{
			"error": "Index '" + indexName + "' doesn't exists.",
		}

		handler.WriteJSONObject(res, response)
		return
	} else if exists == false && err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		handler.Error(res, err.Error())
		return
	}

	docID := handler.GetDocumentID()

	docIntID, err := strconv.Atoi(docID)

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		handler.Error(res, "Invalid document id")
		return
	}

	index, err := handler.search.OpenIndex(indexName)

	if err != nil {
		handler.Error(res, err.Error())
		return
	}

	document, err = index.Get(uint64(docIntID))

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		handler.Error(res, err.Error())
		return
	}

	handler.WriteJSON(res, document)
}
