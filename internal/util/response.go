package util

import (
	"encoding/json"
	"net/http"
)

func BuildErrResponse(w http.ResponseWriter, err error) {
	BuildErrResponseWithCode(w, err, http.StatusInternalServerError)
}

func BuildErrResponseWithCode(w http.ResponseWriter, err error, code int) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"err": err.Error(),
	})
	w.WriteHeader(code)
}

func BuildResponseWithBody(w http.ResponseWriter, body any) {
	BuildResponseBodyWithCode(w, body, http.StatusOK)
}

func BuildResponseBodyWithCode(w http.ResponseWriter, body any, code int) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}
