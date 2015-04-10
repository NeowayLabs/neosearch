package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/NeowayLabs/neosearch"
	"github.com/NeowayLabs/neosearch/neosearch/home"
	"github.com/NeowayLabs/neosearch/neosearch/index"
	"github.com/gorilla/mux"
)

type ServerConfig struct {
	Host string
	Port uint16
}

type HTTPServer struct {
	config *ServerConfig
	router *mux.Router
	search *neosearch.NeoSearch
}

func NewConfig() *ServerConfig {
	return &ServerConfig{}
}

func New(search *neosearch.NeoSearch, config *ServerConfig) (*HTTPServer, error) {
	server := HTTPServer{}
	server.config = config
	server.search = search

	server.router = mux.NewRouter()

	server.createRoutes()
	return &server, nil
}

func (server *HTTPServer) createRoutes() {
	homeHandler := home.HomeHandler{}
	indexHandler := index.New(server.search)
	createIndexHandler := index.NewCreateHandler(server.search)
	deleteIndexHandler := index.NewDeleteHandler(server.search)
	indexGetHandler := index.NewGetHandler(server.search)
	indexAddHandler := index.NewAddHandler(server.search)

	server.router.Handle("/debug/vars", http.DefaultServeMux)

	server.router.Handle("/", &homeHandler).Methods("GET")
	server.router.Handle("/{index}", indexHandler).Methods("GET")
	server.router.Handle("/{index}", createIndexHandler).Methods("PUT")
	server.router.Handle("/{index}", deleteIndexHandler).Methods("DELETE")
	server.router.Handle("/{index}/{id}", indexGetHandler).Methods("GET")
	server.router.Handle("/{index}/{id}", indexAddHandler).Methods("POST")
	//	server.router.Handle("/{index}", indexAddHandler).Methods("POST")
}

func (server *HTTPServer) Start() error {
	hostPort := server.config.Host + ":" + strconv.Itoa(int(server.config.Port))
	log.Printf("Listening on %s", hostPort)
	err := http.ListenAndServe(hostPort, server.router)

	return err
}
