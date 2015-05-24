package index

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/NeowayLabs/neosearch/lib/neosearch"
	"github.com/NeowayLabs/neosearch/lib/neosearch/search"
	"github.com/NeowayLabs/neosearch/service/neosearch/handler"
)

type SearchHandler struct {
	handler.DefaultHandler
	search *neosearch.NeoSearch
}

func NewSearchHandler(search *neosearch.NeoSearch) *SearchHandler {
	return &SearchHandler{
		search: search,
	}
}

func (handler *SearchHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var (
		err        error
		exists     bool
		documents  []map[string]interface{}
		outputJSON []byte
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

	dslBytes, err := ioutil.ReadAll(req.Body)

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		handler.Error(res, err.Error())
		return
	}

	dsl := make(map[string]interface{})

	err = json.Unmarshal(dslBytes, &dsl)

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

	if dsl["query"] == nil {
		res.WriteHeader(http.StatusBadRequest)
		handler.Error(res, "No query field specified")
		return
	}

	query, ok := dsl["query"].(map[string]interface{})

	if !ok {
		res.WriteHeader(http.StatusBadRequest)
		handler.Error(res, "Search 'query' field is not a JSON object")
		return
	}

	output := make(map[string]interface{})
	var total uint64

	docs, total, err := search.Search(index, query, 10)

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		handler.Error(res, err.Error())
		return
	}

	documents = make([]map[string]interface{}, len(docs))

	for idx, doc := range docs {
		obj := make(map[string]interface{})
		err = json.Unmarshal([]byte(doc), &obj)

		if err != nil {
			fmt.Println("Failed to unmarshal: ", doc)
			goto error
		}

		documents[idx] = obj
	}

	output["total"] = total
	output["results"] = documents

	outputJSON, err = json.Marshal(output)

	if err != nil {
		goto error
	}

	res.WriteHeader(http.StatusOK)
	handler.WriteJSON(res, outputJSON)

	return

error:
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		handler.Error(res, err.Error())
		return
	}
}
