package dom

import (
	"strings"

	"golang.org/x/net/html"
)

type DomElement html.Node

func genericSelector(e *DomElement, s string) bool {
	return e.Type == html.ElementNode && e.Data == s
}

func classSelector(e *DomElement, s string) bool {
	if e.Type != html.ElementNode {
		return false
	}

	for _, attr := range e.Attr {
		if attr.Key == "class" && attr.Val == s[1:] {
			return true
		}
	}

	return false
}

func customDataSelector(e *DomElement, s string) bool {
	if e.Type != html.ElementNode {
		return false
	}

	for _, attr := range e.Attr {
		if strings.HasPrefix(attr.Key, "data-") {
			return true
		}
	}

	return false
}

func parseSelector(e *DomElement, s string) bool {
	var fn func(e *DomElement, s string) bool
	fn = genericSelector

	if strings.HasPrefix(s, ".") {
		fn = classSelector
	}

	if strings.HasPrefix(s, "data-") {
		fn = customDataSelector
	}

	return fn(e, s)
}

// QuerySelector receives a selector s that will be retrived from the DOM element, if found.
// Returns a slice of DomElement or an empty slice
func (e *DomElement) QuerySelector(s string) []DomElement {
	var nodes []DomElement

	if parseSelector(e, s) {
		nodes = append(nodes, DomElement(*e))
	}

	for c := e.FirstChild; c != nil; c = c.NextSibling {
		child := DomElement(*c)
		next := child.QuerySelector(s)

		if len(next) > 0 {
			nodes = append(nodes, next...)
		}
	}

	return nodes
}
