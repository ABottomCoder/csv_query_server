package server

import (
	"encoding/json"
	"net/http"
)

func NewServer(addr string, registerHandler func(mux *http.ServeMux)) *http.Server {
	mux := http.NewServeMux()
	registerHandler(mux)

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return server
}

func JsonResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
