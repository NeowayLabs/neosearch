package home

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/NeowayLabs/neosearch/lib/neosearch"
	"github.com/NeowayLabs/neosearch/lib/neosearch/config"
	"github.com/julienschmidt/httprouter"
)

var dataDirTmp string

func init() {
	var err error
	dataDirTmp, err = ioutil.TempDir("/tmp", "neosearch-service-home-")
	if err != nil {
		panic(err)
	}
}

func getHomeHandler() *HomeHandler {
	cfg := config.NewConfig()
	cfg.Option(config.DataDir(dataDirTmp))
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
