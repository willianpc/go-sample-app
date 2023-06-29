package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/willianpc/go-sample-app/dom"
	"golang.org/x/net/html"
)

func getText(n *html.Node) string {
	if n.Type == html.TextNode {
		return strings.ReplaceAll(n.Data, `"`, `\"`)
	}

	if c := n.FirstChild; c != nil {
		return getText(c)
	}

	return ""
}

func main() {
	mux := http.NewServeMux()

	c := http.Client{
		Timeout: time.Second * 30,
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		q := r.URL.Query().Get("q")

		q = url.QueryEscape(q)
		fmt.Println(q)

		if q != "" {
			res, err := c.Get("https://www.google.com/search?q=" + q)

			if err != nil {
				panic(err)
			}

			b, err := io.ReadAll(res.Body)

			// fmt.Println(string(b))

			if err != nil {
				panic(err)
			}

			defer res.Body.Close()

			doc, err := html.Parse(bytes.NewReader(b))

			if err != nil {
				panic(err)
			}

			de := dom.DomElement(*doc)
			nodes := de.QuerySelector("h3")

			buf := []string{}

			for _, node := range nodes {
				n1 := (html.Node)(node)
				text := getText(&n1)

				buf = append(buf, text)
			}

			fmt.Fprintf(w, `{
				"query": "%s",
				"total": %d,
				"results": ["%s"]
			}`, q, len(nodes), strings.Join(buf, `","`))

			return
		}

		io.WriteString(w, `{"results": 0}`)
	})

	server := http.Server{
		Addr:    ":9090",
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}
