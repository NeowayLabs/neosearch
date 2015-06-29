package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/NeowayLabs/neosearch/lib/neosearch"
	"github.com/NeowayLabs/neosearch/service/neosearch/home"
	"github.com/NeowayLabs/neosearch/service/neosearch/index"
	"github.com/julienschmidt/httprouter"
)

type ServerConfig struct {
	Host string
	Port uint16
}

type HTTPServer struct {
	config *ServerConfig
	router *httprouter.Router
	search *neosearch.NeoSearch
}

func NewConfig() *ServerConfig {
	return &ServerConfig{}
}

func New(search *neosearch.NeoSearch, config *ServerConfig) (*HTTPServer, error) {
	server := HTTPServer{}
	server.config = config
	server.search = search

	server.router = httprouter.New()

	server.createRoutes()
	return &server, nil
}

func (server *HTTPServer) createRoutes() {
	homeHandler := home.HomeHandler{}
	indexHandler := index.New(server.search)
	createIndexHandler := index.NewCreateHandler(server.search)
	deleteIndexHandler := index.NewDeleteHandler(server.search)
	getIndexHandler := index.NewGetHandler(server.search)
	getAnalyzeIndexHandler := index.NewGetAnalyzeHandler(server.search)
	addIndexHandler := index.NewAddHandler(server.search)
	searchIndexHandler := index.NewSearchHandler(server.search)

	server.router.Handle("GET", "/", homeHandler.ServeHTTP)
	server.router.Handle("GET", "/:index", indexHandler.ServeHTTP)
	server.router.Handle("PUT", "/:index", createIndexHandler.ServeHTTP)
	server.router.Handle("DELETE", "/:index", deleteIndexHandler.ServeHTTP)
	server.router.Handle("POST", "/:index", searchIndexHandler.ServeHTTP)
	server.router.Handle("GET", "/:index/:id", getIndexHandler.ServeHTTP)
	server.router.Handle("GET", "/:index/:id/_analyze", getAnalyzeIndexHandler.ServeHTTP)
	server.router.Handle("POST", "/:index/:id", addIndexHandler.ServeHTTP)
}

func (server *HTTPServer) GetRoutes() *httprouter.Router {
	return server.router
}

func (server *HTTPServer) Start() error {
	hostPort := server.config.Host + ":" + strconv.Itoa(int(server.config.Port))
	log.Printf("Listening on %s", hostPort)
	err := http.ListenAndServe(hostPort, server.router)

	return err
}
