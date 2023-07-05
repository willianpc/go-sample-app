package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	instana "github.com/instana/go-sensor"
	ot "github.com/opentracing/opentracing-go"
	"github.com/willianpc/go-sample-app/dom"
	"golang.org/x/net/html"
	"golang.org/x/text/encoding/charmap"
)

func sendError(w http.ResponseWriter, err error) {
	fmt.Fprintf(w, `{"error": %s}`, err.Error())
}

func handleSearch() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		q := url.QueryEscape(r.URL.Query().Get("q"))

		if q == "" {
			_, _ = io.WriteString(w, `{"results": 0}`)
			return
		}

		cacheRes := readCache(q)

		var fromCache bool

		if len(cacheRes) == 0 {
			var ps ot.Span
			var ok bool

			if ps, ok = instana.SpanFromContext(r.Context()); !ok {
				// no parent span found, create an artificial parent span
				ps = s.StartSpan("entry-span")
			}

			// Important! To finish the span
			defer ps.Finish()

			// New context containing the parent span, mandatory
			rctx := instana.ContextWithSpan(r.Context(), ps)

			req, err := http.NewRequest(http.MethodGet, "https://www.google.com/search?q="+q, nil)

			if err != nil {
				sendError(w, err)
				return
			}

			res, err := c.Do(req.WithContext(rctx))

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

			if err = writeCache(r.Context(), q, cacheRes); err != nil {
				sendError(w, err)
				return
			}
		} else {
			fromCache = true
		}

		cacheAsArray := `["` + strings.Join(cacheRes, `", "`) + `"]`

		fmt.Fprintf(w, `{"total": %d,"query": "%s","results": %s, "cached": %v}`, len(cacheRes), q, cacheAsArray, fromCache)
	}
}

func handleFunc() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/query", instana.TracingHandlerFunc(s, "/query", handleSearch()))

	return mux
}
