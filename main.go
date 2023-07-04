package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/willianpc/go-sample-app/dom"
	"golang.org/x/net/html"
	"golang.org/x/text/encoding/charmap"
)

func sendError(w http.ResponseWriter, err error) {
	fmt.Fprintf(w, `{"error": %s}`, err.Error())
}

func main() {
	mux := http.NewServeMux()

	c := &http.Client{
		Timeout: time.Second * 30,
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		q := url.QueryEscape(r.URL.Query().Get("q"))

		if q == "" {
			io.WriteString(w, `{"results": 0}`)
			return
		}

		cacheRes := readCache(q)

		var fromCache bool

		if len(cacheRes) == 0 {
			res, err := c.Get("https://www.google.com/search?q=" + q)

			if err != nil {
				sendError(w, err)
				return
			}

			b, err := io.ReadAll(res.Body)

			if err != nil {
				sendError(w, err)
				return
			}

			defer res.Body.Close()

			dec := charmap.Windows1250.NewDecoder()
			b, err = dec.Bytes(b)

			if err != nil {
				sendError(w, err)
				return
			}

			doc, err := html.Parse(bytes.NewReader(b))

			if err != nil {
				sendError(w, err)
				return
			}

			de := dom.DomElement(*doc)
			nodes := de.QuerySelector("h3")

			for _, node := range nodes {
				n := html.Node(node)
				text := dom.InnerText(n)
				cacheRes = append(cacheRes, text)
			}

			if err = writeCache(q, cacheRes); err != nil {
				sendError(w, err)
				return
			}
		} else {
			fromCache = true
		}

		cacheAsArray := `["` + strings.Join(cacheRes, `", "`) + `"]`

		fmt.Fprintf(w, `{
				"total": %d,
				"query": "%s",
				"results": %s,
				"cached": %v
			}
			`, len(cacheRes), q, cacheAsArray, fromCache)
	})

	port := os.Getenv("PORT")

	if port == "" {
		port = "9090"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}
