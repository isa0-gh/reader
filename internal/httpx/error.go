// Package httpx holds small HTTP response helpers shared across handlers and
// middleware, which can't import each other's error-writing logic directly
// without an import cycle.
package httpx

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

// WriteError mirrors http.Error's signature (message, then status code) so
// call sites can switch over with a plain string replace, but returns a JSON
// body consistent with the rest of the API instead of plain text.
func WriteError(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: error})
}
