package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

var (
	c *http.Client
)

func init() {
	c = &http.Client{
		Timeout: time.Second * 30,
	}
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "9090"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: handleFunc(),
	}

	log.Fatal(server.ListenAndServe())
}
