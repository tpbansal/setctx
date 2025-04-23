package rpc

import "net/http"

// HandleRequest processes incoming requests
func HandleRequest(w http.ResponseWriter, r *http.Request) {
	// Dummy response
	w.Write([]byte("Request handled successfully"))
}
