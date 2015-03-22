package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

type DefaultHandler struct {
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
