package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/willianpc/go-sample-app/dom"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/net/html"
	"golang.org/x/text/encoding/charmap"
)

func handleSearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	q := url.QueryEscape(r.URL.Query().Get("q"))

	// Throws HTTP 500 Error
	if q == "" {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, `{"error": "A value for 'q' must be provided"}`)
		return
	}

	// Read from cache
	cacheRes := readCache(r.Context(), q)

	var fromCache bool

	if len(cacheRes) > 0 {
		fromCache = true
	} else {
		cacheRes = dataFromGoogle(r, q)
	}

	cacheAsArray := `["` + strings.Join(cacheRes, `", "`) + `"]`

	fmt.Fprintf(w, `{"total": %d,"query": "%s","results": %s, "cached": %v}`, len(cacheRes), q, cacheAsArray, fromCache)
}

func dataFromGoogle(incomingRequest *http.Request, q string) []string {
	var cacheRes []string

	tr := otel.Tracer("otel-tracer")
	ctx := incomingRequest.Context()

	var parentSpan trace.Span
	parentSpan = trace.SpanFromContext(ctx)

	if !parentSpan.SpanContext().HasSpanID() {
		opts := []trace.SpanStartOption{
			trace.WithSpanKind(trace.SpanKindServer),
		}

		ctx, parentSpan = tr.Start(ctx, "client-call-parent-span", opts...)
	}

	defer parentSpan.End()
	clientReq, _ := http.NewRequestWithContext(ctx, "GET", "https://www.google.com/search?q="+q, nil)

	clientResp, _ := c.Do(clientReq)

	body, _ := io.ReadAll(clientResp.Body)
	defer clientResp.Body.Close()

	dec := charmap.Windows1250.NewDecoder()
	body, _ = dec.Bytes(body)

	doc, _ := html.Parse(bytes.NewReader(body))

	de := dom.DomElement(*doc)
	nodes := de.QuerySelector("h3")

	for _, node := range nodes {
		n := html.Node(node)
		text := dom.InnerText(n)
		cacheRes = append(cacheRes, text)
	}

	// Update cache
	_ = writeCache(incomingRequest.Context(), q, cacheRes)

	return cacheRes
}
