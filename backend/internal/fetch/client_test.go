package fetch_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chanaka-withanage/page-analyzer/internal/fetch"
)

func TestFetch_AllowsLocalWhenEnabled(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	}))
	defer ts.Close()

	c := fetch.New(5*time.Second, 3, 1024)
	c.AllowLocal()

	resp, body, err := c.Get(context.Background(), ts.URL)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	defer body.Close()

	data, _ := io.ReadAll(body)
	if string(data) != "hello" {
		t.Errorf("expected 'hello', got %q", string(data))
	}
	resp.Body.Close()
}

func TestFetch_BlocksLocalByDefault(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("nope"))
	}))
	defer ts.Close()

	c := fetch.New(5*time.Second, 3, 1024)

	_, _, err := c.Get(context.Background(), ts.URL)
	if err == nil || err != fetch.ErrPrivateAddr {
		t.Fatalf("expected ErrPrivateAddr, got %v", err)
	}
}

func TestFetch_Timeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.Write([]byte("late"))
	}))
	defer ts.Close()

	c := fetch.New(100*time.Millisecond, 3, 1024)
	c.AllowLocal()

	_, _, err := c.Get(context.Background(), ts.URL)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

func TestFetch_RespectsMaxBytes(t *testing.T) {
	// server sends a large response
	large := make([]byte, 2048)
	for i := range large {
		large[i] = 'x'
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(large)
	}))
	defer ts.Close()

	c := fetch.New(1*time.Second, 3, 1024) // cap at 1KB
	c.AllowLocal()

	_, body, err := c.Get(context.Background(), ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer body.Close()

	data, _ := io.ReadAll(body)
	if len(data) != 1024 {
		t.Errorf("expected max 1024 bytes, got %d", len(data))
	}
}
