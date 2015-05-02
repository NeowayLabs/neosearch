package index

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/NeowayLabs/neosearch/lib/neosearch"
	"github.com/NeowayLabs/neosearch/service/neosearch/handler"
)

type AddHandler struct {
	handler.DefaultHandler
	search *neosearch.NeoSearch
}

func NewAddHandler(search *neosearch.NeoSearch) *AddHandler {
	return &AddHandler{
		search: search,
	}
}

func (handler *AddHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	handler.ProcessVars(req)
	indexName := handler.GetIndexName()

	if exists, err := handler.search.IndexExists(indexName); exists != true && err == nil {
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

	if req.Method != "POST" {
		res.WriteHeader(http.StatusBadRequest)
		handler.Error(res, "Add document expect a POST request")
		return
	}

	docID := handler.GetDocumentID()

	docIntID, err := strconv.Atoi(docID)

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		handler.Error(res, "Invalid document id")
		return
	}

	document, err := ioutil.ReadAll(req.Body)

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		handler.Error(res, err.Error())
		return
	}

	index, err := handler.search.OpenIndex(indexName)

	if err != nil {
		handler.Error(res, err.Error())
		return
	}

	err = index.Add(uint64(docIntID), document)

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		handler.Error(res, err.Error())
		return
	}

	handler.WriteJSON(res, []byte(fmt.Sprintf("{\"status\": \"Document %d indexed.\"}", docIntID)))
}
