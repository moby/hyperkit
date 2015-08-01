// Package u12 exposes a CORS http.Handler wrapper.
package u12

import "net/http"

// CORS wraps an http.Handler and sets "Access-Control-Allow-Origin" header to "*".
type CORS struct {
	http.Handler
}

func (h CORS) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	h.Handler.ServeHTTP(w, req)
}
