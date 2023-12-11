package main

import (
	"log"
	"net/http"
	"time"

	redis "github.com/go-redis/redis/v8"
	instana "github.com/instana/go-sensor"
	"github.com/instana/go-sensor/instrumentation/instaredis"
)

var (
	col instana.TracerLogger
	c   *http.Client
	rdb *redis.Client
)

func init() {
	col = instana.InitCollector(&instana.Options{
		Service:   "Go Sample App",
		AgentHost: "instana-agent.instana-agent.svc.cluster.local",
	})

	rdb = redis.NewClient(&redis.Options{Addr: "redis.lero.svc.cluster.local:6379"})

	instaredis.WrapClient(rdb, col)

	c = &http.Client{
		Timeout:   time.Second * 30,
		Transport: instana.RoundTripper(col, nil),
	}
}

func main() {
	http.HandleFunc("/query", instana.TracingHandlerFunc(col, "/query", handleSearch))

	log.Fatal(http.ListenAndServe(":9090", nil))
}
