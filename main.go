package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	redis "github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var (
	c   *http.Client
	rdb *redis.Client
)

func init() {
	rdb = redis.NewClient(&redis.Options{Addr: ":6379"})

	if err := redisotel.InstrumentTracing(rdb); err != nil {
		panic(err)
	}

	c = &http.Client{
		Timeout:   time.Second * 30,
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
}

func main() {
	ctx := context.Background()

	fn, err := initProvider()

	if err != nil {
		log.Fatal("Trace provider not initialized", err)
	}

	defer func() {
		log.Println("Shutting down the trace provider")
		fn(ctx)
	}()

	http.Handle("/query", otelhttp.NewHandler(http.HandlerFunc(handleSearch), "/query"))

	log.Fatal(http.ListenAndServe("localhost:9090", nil))
}
