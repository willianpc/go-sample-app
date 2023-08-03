package main

import (
	"log"
	"net/http"
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
	mux := http.NewServeMux()
	mux.HandleFunc("/query", handleSearch)

	server := &http.Server{
		Addr:    ":9090",
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}
