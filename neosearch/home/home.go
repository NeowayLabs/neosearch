package home

import "net/http"

type HomeHandler struct {
}

func (handler *HomeHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("HOME"))
}
