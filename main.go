package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	var h handler = handleFunc()
	port := os.Getenv("PORT")

	if port == "" {
		port = "9090"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: h,
	}

	log.Fatal(server.ListenAndServe())
}
