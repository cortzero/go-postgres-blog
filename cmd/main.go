package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/cortzero/go-postgres-blog/internal/data"
	"github.com/cortzero/go-postgres-blog/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	// Loading environment variables
	godotenv.Load()
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	// Instantiating the server
	serv := server.New(host, port)

	// Instantiating a database connection
	database := data.New()
	if err := database.DB.Ping(); err != nil {
		log.Fatal(err)
	}

	// Starting the server
	go serv.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	serv.Close()
	data.Close()
}
