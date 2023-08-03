package main

import (
	"log"
	"net/http"
	"os"
	"time"

	redis "github.com/go-redis/redis/v8"
)

var (
	c   *http.Client
	rdb *redis.Client
)

func init() {
	rdb = redis.NewClient(&redis.Options{Addr: ":6379"})

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
