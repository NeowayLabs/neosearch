package home

import (
	"net/http/httptest"
	"testing"

	"github.com/NeowayLabs/neosearch/lib/neosearch"
	"github.com/julienschmidt/httprouter"
)

func getHomeHandler() *HomeHandler {
	cfg := neosearch.NewConfig()
	cfg.Option(neosearch.DataDir("/tmp/"))
	ns := neosearch.New(cfg)
	return NewHomeHandler(ns)
}

func TestHomeInfo(t *testing.T) {
	handler := getHomeHandler()

	router := httprouter.New()

	router.Handle("GET", "/", handler.ServeHTTP)

	// router.HandleFunc("/{index}", func(res http.ResponseWriter, req *http.Request) {
	// 	handler.ServeHTTP(res, req)
	// }).Methods("DELETE")

	ts := httptest.NewServer(router)

	defer func() {
		ts.Close()
		handler.search.Close()
	}()
}
