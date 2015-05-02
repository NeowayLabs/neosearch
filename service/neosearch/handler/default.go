package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type DefaultHandler struct {
	requestVars map[string]string
}

func (h *DefaultHandler) Error(res http.ResponseWriter, errMessage string) {
	errObject := map[string]interface{}{
		"error": errMessage,
	}

	body, err := json.Marshal(errObject)

	if err != nil {
		log.Println("Failed to marshal error object")
		return
	}

	h.WriteJSON(res, body)
}

func (h *DefaultHandler) WriteJSON(res http.ResponseWriter, content []byte) {
	res.Header().Set("Content-Type", "application/json")
	res.Write(content)
}

func (h *DefaultHandler) WriteJSONObject(res http.ResponseWriter, content interface{}) {
	res.Header().Set("Content-Type", "application/json")

	body, err := json.Marshal(content)

	if err != nil {
		log.Printf("Failed to marshal JSON: %s", err.Error())
		return
	}

	h.WriteJSON(res, body)
}

func (h *DefaultHandler) ProcessVars(req *http.Request) map[string]string {
	h.requestVars = mux.Vars(req)

	return h.requestVars
}

func (h *DefaultHandler) GetIndexName() string {
	if h.requestVars == nil {
		return ""
	}

	return h.requestVars["index"]
}

func (h *DefaultHandler) GetDocumentID() string {
	if h.requestVars == nil {
		return ""
	}

	return h.requestVars["id"]
}
