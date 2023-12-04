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
	rdb = redis.NewClient(&redis.Options{Addr: "redis.default.svc.cluster.local:6379"})

	c = &http.Client{
		Timeout: time.Second * 30,
	}
}

func main() {
	http.HandleFunc("/query", handleSearch)

	log.Fatal(http.ListenAndServe(":9090", nil))
}
