package internalhttp

import (
	"net/http"
)

//nolint:deadcode,unused
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO
	})
}
