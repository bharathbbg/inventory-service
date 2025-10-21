// Initialize REST server
package rest

import (
	"encoding/json"
	"net/http"
)

// NewRouter returns an http.Handler for the REST API.
// It accepts the service as an interface{} to avoid coupling to concrete service types here.
// Add more routes and handlers that use the actual service methods as you implement them.
func NewRouter(_ interface{}) http.Handler {
	mux := http.NewServeMux()

	// simple health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// TODO: add inventory endpoints that use the provided service, e.g.:
	// mux.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) { ... })

	return mux
}
