package index

import (
	"net/http"
	"strconv"

	"github.com/NeowayLabs/neosearch/lib/neosearch"
	"github.com/NeowayLabs/neosearch/lib/neosearch/engine"
	"github.com/NeowayLabs/neosearch/service/neosearch/handler"
)

type GetAnalyseHandler struct {
	handler.DefaultHandler

	search *neosearch.NeoSearch
}

func NewGetAnalyzeHandler(search *neosearch.NeoSearch) *GetAnalyseHandler {
	handler := GetAnalyseHandler{
		search: search,
	}

	return &handler
}

func (handler *GetAnalyseHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var (
		err    error
		cmd    engine.Command
		exists bool
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

	cmd, err = index.GetAnalyze(uint64(docIntID))

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		handler.Error(res, err.Error())
		return
	}

	res.Write([]byte(cmd.Reverse()))
}
