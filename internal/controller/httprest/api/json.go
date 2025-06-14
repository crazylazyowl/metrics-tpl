package api

import (
	"encoding/json"
	"net/http"
)

type Validatable interface {
	Validate() error
}

func readJSON(r *http.Request, j Validatable) error {
	if err := json.NewDecoder(r.Body).Decode(j); err != nil {
		return err
	}
	return j.Validate()
}

func writeJSON(w http.ResponseWriter, status int, j any) error {
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(j)
}
