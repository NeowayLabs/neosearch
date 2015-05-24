package index

import (
	"errors"
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
	var (
		document []byte
		err      error
		exists   bool
		docID    string
		docIntID int
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

	if req.Method != "POST" {
		err = errors.New("Add document expect a POST request")
		goto error_fatal
	}

	docID = handler.GetDocumentID()
	docIntID, err = strconv.Atoi(docID)

	if err != nil {
		goto error_fatal
	}

	document, err = ioutil.ReadAll(req.Body)

	if err != nil {
		goto error_fatal
	}

	err = handler.addDocument(indexName, uint64(docIntID), document)

	handler.WriteJSON(res, []byte(fmt.Sprintf("{\"status\": \"Document %d indexed.\"}", docIntID)))

	return

error_fatal:
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		handler.Error(res, err.Error())
		return
	}
}

func (handler *AddHandler) addDocument(indexName string, id uint64, document []byte) error {
	index, err := handler.search.OpenIndex(indexName)

	if err != nil {
		return err
	}

	return index.Add(id, document)
}
