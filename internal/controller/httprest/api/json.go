package api

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type Validatable interface {
	Validate() error
}

func readCompressedJSON(r *http.Request, j Validatable) error {
	if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		return readJSON(r.Body, j)
	}
	reader, err := gzip.NewReader(r.Body)
	if err != nil {
		return err
	}
	defer reader.Close()
	return readJSON(reader, j)
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
