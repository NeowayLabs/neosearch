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

	server.router.Handle("/", &homeHandler)
	server.router.Handle("/{index}", indexHandler)
}

func (server *HTTPServer) Start() error {
	hostPort := server.config.Host + ":" + strconv.Itoa(int(server.config.Port))
	log.Printf("Listening on %s", hostPort)
	err := http.ListenAndServe(hostPort, server.router)

	return err
}
