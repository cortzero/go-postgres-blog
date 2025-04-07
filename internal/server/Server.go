package server

import (
	"log"
	"net/http"

	"github.com/cortzero/go-postgres-blog/internal/data"
	"github.com/cortzero/go-postgres-blog/internal/server/handlers"
	"github.com/cortzero/go-postgres-blog/internal/service/services"
)

// Server contains a server configuration
type Server struct {
	server *http.Server
}

func New(host string, port string) *Server {
	// Database Connection
	conn := data.New()

	// User Service
	userService := services.NewUserService(data.NewUserRepository(conn))

	// User Handler
	userHandler := handlers.NewUserHandler(userService)

	// Post Handler
	postHandler := handlers.NewPostHandler(data.NewPostRepository(conn))

	// Creating the Server Mux
	mux := http.NewServeMux()

	// Mapping User endpoints to the user handler
	mux.Handle("/api/v1/users", userHandler)
	mux.Handle("/api/v1/users/", userHandler)
	mux.Handle("/api/v1/users/{id}", userHandler)
	mux.Handle("/api/v1/users/{id}/", userHandler)

	// Mapping Post endpoints to the post handler
	mux.Handle("/api/v1/posts", postHandler)
	mux.Handle("/api/v1/posts/", postHandler)
	mux.Handle("/api/v1/posts/{id}", postHandler)
	mux.Handle("/api/v1/posts/{id}/", postHandler)

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
