package main

import (
	"log"
	"net/http"
	"os"
	"time"

	redis "github.com/go-redis/redis/v8"
	instana "github.com/instana/go-sensor"
	"github.com/instana/go-sensor/instrumentation/instaredis"
)

var (
	s   instana.TracerLogger
	c   *http.Client
	rdb *redis.Client
)

func init() {
	s = instana.InitCollector(&instana.Options{
		Service: "Go Sample App",
	})

	rdb = redis.NewClient(&redis.Options{Addr: ":6379"})

	instaredis.WrapClient(rdb, s)

	c = &http.Client{
		Timeout:   time.Second * 30,
		Transport: instana.RoundTripper(s, nil),
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
