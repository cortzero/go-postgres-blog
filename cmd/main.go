package main

import (
	"os"
	"os/signal"

	"github.com/cortzero/go-postgres-blog/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	serv := server.New(host, port)

	go serv.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	serv.Close()
}
