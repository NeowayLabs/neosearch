package home

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NeowayLabs/neosearch/lib/neosearch"
	"github.com/gorilla/mux"
)

func getHomeHandler() *HomeHandler {
	cfg := neosearch.NewConfig()
	cfg.Option(neosearch.DataDir("/tmp/"))
	ns := neosearch.New(cfg)
	return NewHomeHandler(ns)
}

func TestHomeInfo(t *testing.T) {
	handler := getHomeHandler()

	router := mux.NewRouter()
	router.HandleFunc("/{index}", func(res http.ResponseWriter, req *http.Request) {
		handler.ServeHTTP(res, req)
	}).Methods("DELETE")

	ts := httptest.NewServer(router)

	defer func() {
		ts.Close()
		handler.search.Close()
	}()
}
