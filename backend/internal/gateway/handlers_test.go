package gateway_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chanaka-withanage/page-analyzer/internal/analyzer"
	"github.com/chanaka-withanage/page-analyzer/internal/fetch"
	"github.com/chanaka-withanage/page-analyzer/internal/gateway"
	"github.com/chanaka-withanage/page-analyzer/pkg/contract"
)

// helper: create a minimal backend with fetch.AllowLocal enabled
func newTestHandler() http.Handler {
	f := fetch.New(5*time.Second, 3, 1<<20)
	f.AllowLocal()
	svc := analyzer.New(f)
	svc.SetDefaultTimeout(5 * time.Second)
	return gateway.NewMuxWithService(svc)
}

// fake upstream page
func startFakePage(html string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(html))
	}))
}

func TestAnalyzeHandler_ReturnsJSON(t *testing.T) {
	// given: fake page with title + heading
	page := startFakePage(`
		<!DOCTYPE html>
		<html>
		  <head><title>Gateway Test</title></head>
		  <body><h1>Hello</h1></body>
		</html>`)
	defer page.Close()

	handler := newTestHandler()
	srv := httptest.NewServer(handler)
	defer srv.Close()

	// send request to /api/analyze
	reqBody, _ := json.Marshal(map[string]string{"url": page.URL})
	resp, err := http.Post(srv.URL+"/api/analyze", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		t.Fatalf("POST /api/analyze failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}

	var result contract.AnalyzeResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}

	// assertions
	if result.Title != "Gateway Test" {
		t.Errorf("expected title 'Gateway Test', got %q", result.Title)
	}
	if result.Headings["h1"] != 1 {
		t.Errorf("expected 1 h1 heading, got %d", result.Headings["h1"])
	}
	if len(result.Errors) != 0 {
		t.Errorf("expected no errors, got %v", result.Errors)
	}
}

func TestAnalyzeHandler_InvalidJSON(t *testing.T) {
	handler := newTestHandler()
	srv := httptest.NewServer(handler)
	defer srv.Close()

	resp, err := http.Post(srv.URL+"/api/analyze", "application/json", bytes.NewBuffer([]byte(`{bad-json}`)))
	if err != nil {
		t.Fatalf("POST /api/analyze failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d", resp.StatusCode)
	}
}
