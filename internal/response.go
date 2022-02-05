package common

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondError(w http.ResponseWriter, status int, err error) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	type ErrorResponse struct {
		Errors []string `json:"errors"`
	}
	resp := &ErrorResponse{Errors: make([]string, 0, 1)}
	if err != nil {
		resp.Errors = append(resp.Errors, err.Error())
	}

	enc := json.NewEncoder(w)
	encodingErr := enc.Encode(resp)
	if encodingErr != nil {
		log.Printf("%s", encodingErr.Error())
	}
}

func RespondWithStatus(w http.ResponseWriter, body interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")

	if body == nil {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(status)
		enc := json.NewEncoder(w)
		encodingErr := enc.Encode(body)
		if encodingErr != nil {
			log.Printf("%s", encodingErr.Error())
		}
	}
}
