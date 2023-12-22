package server

import (
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
