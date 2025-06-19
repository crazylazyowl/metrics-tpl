package api

import (
	"encoding/json"
	"io"
	"net/http"
)

type Validatable interface {
	Validate() error
}

func readJSON(r io.Reader, j Validatable) error {
	if err := json.NewDecoder(r).Decode(j); err != nil {
		return err
	}
	return j.Validate()
}

func writeJSON(w http.ResponseWriter, status int, j any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(j)
}
