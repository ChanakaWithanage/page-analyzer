package analyzer_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/chanaka-withanage/page-analyzer/internal/analyzer"
	"github.com/chanaka-withanage/page-analyzer/internal/fetch"
	"github.com/chanaka-withanage/page-analyzer/pkg/contract"
)

// helper: create an Analyzer Service with a short timeout
func newTestService(t *testing.T) *analyzer.Service {
    t.Helper()
    f := fetch.New(5*time.Second, 3, 1<<20) // 1MB cap
    f.AllowLocal()
    svc := analyzer.New(f)
    svc.SetDefaultTimeout(5 * time.Second)
    return svc
}

func TestAnalyze_BasicHTML(t *testing.T) {
	// fake HTML page
	html := `
		<!DOCTYPE html>
		<html>
		  <head><title>Test Page</title></head>
		  <body>
		    <h1>Heading One</h1>
		    <a href="/internal">Internal</a>
		    <a href="https://external.com">External</a>
		    <form><input type="password" /></form>
		  </body>
		</html>`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(html))
	}))
	defer ts.Close()

	svc := newTestService(t)

	res, err := svc.Analyze(context.Background(), contract.AnalyzeParams{
		URL: ts.URL,
	})
	if err != nil {
		t.Fatalf("Analyze returned error: %v", err)
	}

	// assertions
	if res.Title != "Test Page" {
		t.Errorf("expected title=Test Page, got %q", res.Title)
	}
	if res.Headings["h1"] != 1 {
		t.Errorf("expected 1 h1 heading, got %d", res.Headings["h1"])
	}
	if res.LinksInternal != 1 {
		t.Errorf("expected 1 internal link, got %d", res.LinksInternal)
	}
	if res.LinksExternal != 1 {
		t.Errorf("expected 1 external link, got %d", res.LinksExternal)
	}
	if !res.LoginFormPresent {
		t.Errorf("expected login form present, got false")
	}
	if len(res.Errors) != 0 {
		t.Errorf("expected no errors, got %v", res.Errors)
	}
}

func TestAnalyze_UpstreamError(t *testing.T) {
	// server that always returns 500
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "fail", http.StatusInternalServerError)
	}))
	defer ts.Close()

	svc := newTestService(t)

	res, err := svc.Analyze(context.Background(), contract.AnalyzeParams{
		URL: ts.URL,
	})
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if !strings.Contains(res.Errors[0], "upstream status") {
		t.Errorf("expected upstream status error, got %v", res.Errors)
	}
}

func TestAnalyze_InvalidHTML(t *testing.T) {
	// invalid HTML
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html><head><title>Broken</title></head><body><h1>oops"))
	}))
	defer ts.Close()

	svc := newTestService(t)

	res, err := svc.Analyze(context.Background(), contract.AnalyzeParams{
		URL: ts.URL,
	})
	if err != nil {
		if len(res.Errors) == 0 {
			t.Errorf("expected parser error in Errors, got none")
		}
	}
}
