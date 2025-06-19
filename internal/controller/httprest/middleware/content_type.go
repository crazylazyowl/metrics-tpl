package middleware

import (
	"fmt"
	"net/http"
)

const (
	jsonContentType     = "application/json"
	textHtmlContentType = "text/html"
)

func JSONContentType(next http.Handler) http.Handler {
	return ContentType(jsonContentType, next)
}

func ContentType(contentType string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != contentType {
			http.Error(w, fmt.Sprintf("Content-Type must be '%s'", contentType), http.StatusUnsupportedMediaType)
			return
		}
		next.ServeHTTP(w, r)
	})
}
