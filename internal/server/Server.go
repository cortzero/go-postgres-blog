package server

import (
	"log"
	"net/http"

	"github.com/cortzero/go-postgres-blog/internal/data"
	"github.com/cortzero/go-postgres-blog/internal/server/handlers"
)

// Server contains a server configuration
type Server struct {
	server *http.Server
}

func New(host string, port string) *Server {
	userHandler := handlers.NewUserHandler(data.NewUserRepository())
	mux := http.NewServeMux()
	mux.Handle("/api/v1/users", userHandler)
	mux.Handle("/api/v1/users/", userHandler)
	mux.Handle("/api/v1/users/{id}", userHandler)
	mux.Handle("/api/v1/users/{id}/", userHandler)
	return &Server{
		server: &http.Server{
			Addr:    host + ":" + port,
			Handler: mux,
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
