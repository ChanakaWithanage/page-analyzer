package parser_test

import (
	"io"
	"net/url"
	"strings"
	"testing"

	"github.com/chanaka-withanage/page-analyzer/internal/parser"
)

func mustURL(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}
	return u
}

func TestParse_BasicFields(t *testing.T) {
	html := `
	<!DOCTYPE html>
	<html>
	  <head>
	    <title>My Test Page</title>
	  </head>
	  <body>
	    <h1>Main Heading</h1>
	    <h2>Sub Heading</h2>
	    <a href="/internal">Internal Link</a>
	    <a href="https://example.com">External Link</a>
	    <form><input type="password" /></form>
	  </body>
	</html>`

	r := io.NopCloser(strings.NewReader(html))
	u := mustURL("http://test.local/")

	res, err := parser.Parse(r, u)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Title != "My Test Page" {
		t.Errorf("expected title 'My Test Page', got %q", res.Title)
	}
	if res.Headings["h1"] != 1 {
		t.Errorf("expected 1 h1 heading, got %d", res.Headings["h1"])
	}
	if res.Headings["h2"] != 1 {
		t.Errorf("expected 1 h2 heading, got %d", res.Headings["h2"])
	}
	if len(res.Links) != 2 {
		t.Errorf("expected 2 links, got %d", len(res.Links))
	}
	if !res.LoginFormPresent {
		t.Errorf("expected login form present, got false")
	}
}

func TestParse_HTML5Doctype(t *testing.T) {
	html := "<!DOCTYPE html><html><head><title>X</title></head><body></body></html>"
	r := io.NopCloser(strings.NewReader(html))
	u := mustURL("http://test.local/")

	res, err := parser.Parse(r, u)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.HTMLVersion != "HTML5" {
		t.Errorf("expected HTML5 doctype, got %q", res.HTMLVersion)
	}
}

func TestParse_NoTitle(t *testing.T) {
	html := "<!DOCTYPE html><html><body><h1>No Title</h1></body></html>"
	r := io.NopCloser(strings.NewReader(html))
	u := mustURL("http://test.local/")

	res, err := parser.Parse(r, u)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Title != "" {
		t.Errorf("expected empty title, got %q", res.Title)
	}
}
