package home

import (
	"net/http"

	"github.com/NeowayLabs/neosearch/lib/neosearch"
	"github.com/NeowayLabs/neosearch/lib/neosearch/version"
	"github.com/NeowayLabs/neosearch/service/neosearch/handler"
	"github.com/julienschmidt/httprouter"
)

type HomeHandler struct {
	search *neosearch.NeoSearch
	handler.DefaultHandler
}

func NewHomeHandler(ns *neosearch.NeoSearch) *HomeHandler {
	return &HomeHandler{
		search: ns,
	}
}

func (handler *HomeHandler) ServeHTTP(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	response := map[string]string{
		"version": version.Version,
		"status":  "alive",
	}

	handler.WriteJSONObject(res, response)
}
