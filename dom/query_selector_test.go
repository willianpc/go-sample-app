package dom_test

import (
	"strings"
	"testing"

	"github.com/willianpc/go-sample-app/dom"
	"golang.org/x/net/html"
)

var rawHtml = `<html>
<body>
	<h1>Hello World</h1>
	<h2>Testing the Query Selector</h2>
	<a href="https://google.com" class="link-class" data-layout="plain">Click here to continue</a>
	<a href="https://exit.cjb.net" class="link-class" data-layout="plain">Or click here to exit</a>
	<span data-layout="plain">This is a simple spam with a custom data</span>
</body>
</html>`

func Test_getA(t *testing.T) {
	doc, err := html.Parse(strings.NewReader(rawHtml))

	if err != nil {
		t.Fatal(err)
	}

	de := dom.DomElement(*doc)

	res := de.QuerySelector("a")

	if l := len(res); l != 2 {
		t.Fatalf("expected 2 but got %d", l)
	}
}

func Test_getClass(t *testing.T) {
	doc, err := html.Parse(strings.NewReader(rawHtml))

	if err != nil {
		t.Fatal(err)
	}

	de := dom.DomElement(*doc)

	res := de.QuerySelector(".link-class")

	if l := len(res); l != 2 {
		t.Fatalf("expected 2 but got %d", l)
	}
}

func Test_getCustomData(t *testing.T) {
	doc, err := html.Parse(strings.NewReader(rawHtml))

	if err != nil {
		t.Fatal(err)
	}

	de := dom.DomElement(*doc)

	res := de.QuerySelector("data-layout")

	if l := len(res); l != 3 {
		t.Fatalf("expected 2 but got %d", l)
	}
}
