package index

import (
	"net/http"

	"github.com/NeowayLabs/neosearch"
)

type IndexHandle struct {
	search *neosearch.NeoSearch
}

func (handler *IndexHandle) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("test server handler"))
}
