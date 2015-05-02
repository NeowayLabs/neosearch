package home

import (
	"net/http"

	"github.com/NeowayLabs/neosearch/lib/neosearch/version"
	"github.com/NeowayLabs/neosearch/service/neosearch/handler"
)

type HomeHandler struct {
	handler.DefaultHandler
}

func (handler *HomeHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	response := map[string]string{
		"version": version.Version,
		"status":  "alive",
	}

	handler.WriteJSONObject(res, response)
}
