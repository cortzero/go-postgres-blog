package server

import (
	"log"
	"net/http"
)

// Server contains a server configuration
type Server struct {
	server *http.Server
}

func New(host string, port string) *Server {
	return &Server{
		server: &http.Server{
			Addr: host + ":" + port,
		},
	}
}

func (serv *Server) Start() {
	log.Printf("Server running on http://%s", serv.server.Addr)
	log.Fatal(serv.server.ListenAndServe())
}

func (serv *Server) Close() error {
	// TODO: add resource closure
	return nil
}
