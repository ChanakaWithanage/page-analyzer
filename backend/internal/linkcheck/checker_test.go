package linkcheck_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"
	"time"

	"github.com/chanaka-withanage/page-analyzer/internal/linkcheck"
)

func mustURL(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}
	return u
}

func TestValidate_SuccessAndFailure(t *testing.T) {
	// healthy server
	okSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer okSrv.Close()

	// failing server
	failSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "fail", http.StatusInternalServerError)
	}))
	defer failSrv.Close()

	checker := linkcheck.New(5, 2, 1*time.Second)
	links := []*url.URL{mustURL(okSrv.URL), mustURL(failSrv.URL)}

	results := checker.Validate(context.Background(), links)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	foundOK, foundFail := false, false
	for _, r := range results {
		if r.URL == okSrv.URL && r.Accessible {
			foundOK = true
		}
		if r.URL == failSrv.URL && !r.Accessible {
			foundFail = true
		}
	}
	if !foundOK {
		t.Errorf("expected OK server marked accessible")
	}
	if !foundFail {
		t.Errorf("expected failing server marked inaccessible")
	}
}

func TestValidate_RespectsGlobalConcurrency(t *testing.T) {
	var concurrent int32
	slowest := make(chan struct{})

	// server that blocks until release
	blockingSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cur := atomic.AddInt32(&concurrent, 1)
		if cur > 2 {
			t.Errorf("expected max 2 concurrent requests, got %d", cur)
		}
		<-slowest
		atomic.AddInt32(&concurrent, -1)
	}))
	defer blockingSrv.Close()

	checker := linkcheck.New(2, 2, 1*time.Second)
	var links []*url.URL
	for i := 0; i < 5; i++ {
		links = append(links, mustURL(blockingSrv.URL))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	go func() {
		time.Sleep(100 * time.Millisecond)
		close(slowest)
	}()

	_ = checker.Validate(ctx, links)
}

func TestValidate_Timeout(t *testing.T) {
	// server that sleeps too long
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	checker := linkcheck.New(5, 2, 50*time.Millisecond)
	results := checker.Validate(context.Background(), []*url.URL{mustURL(ts.URL)})

	if results[0].Accessible {
		t.Errorf("expected inaccessible due to timeout, got accessible")
	}
}
