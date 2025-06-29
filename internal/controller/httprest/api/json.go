package api

import (
	"encoding/json"
	"net/http"
)

func readJSON(r *http.Request, j any) error {
	return json.NewDecoder(r.Body).Decode(j)
}

func writeJSON(w http.ResponseWriter, status int, j any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(j)
}
