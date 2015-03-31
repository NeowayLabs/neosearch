package home

import (
	"net/http"

	"github.com/NeowayLabs/neosearch/neosearch/handler"
	"github.com/NeowayLabs/neosearch/version"
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
