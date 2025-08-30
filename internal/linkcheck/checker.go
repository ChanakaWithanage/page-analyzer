package linkcheck

import (
	"context"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Result struct {
	URL        string
	Accessible bool
	StatusCode int
	Err        string
}

type Checker struct {
	client            *http.Client
	globalConcurrency int
	perHostLimit      int
	timeout           time.Duration
}

func New(globalConcurrency, perHostLimit int, timeout time.Duration) *Checker {
	return &Checker{
		client: &http.Client{
			Timeout: timeout,
		},
		globalConcurrency: globalConcurrency,
		perHostLimit:      perHostLimit,
		timeout:           timeout,
	}
}

func (c *Checker) Validate(ctx context.Context, links []*url.URL) []Result {
	results := make([]Result, len(links))

	globalSem := make(chan struct{}, c.globalConcurrency)
	hostSems := sync.Map{}

	var wg sync.WaitGroup
	for i, link := range links {
		wg.Add(1)
		go func(i int, u *url.URL) {
			defer wg.Done()

			// global limit
			select {
			case globalSem <- struct{}{}:
				defer func() { <-globalSem }()
			case <-ctx.Done():
				results[i] = Result{URL: u.String(), Err: "context cancelled"}
				return
			}

			// per-host limit
			h := u.Hostname()
			val, _ := hostSems.LoadOrStore(h, make(chan struct{}, c.perHostLimit))
			hostSem := val.(chan struct{})
			select {
			case hostSem <- struct{}{}:
				defer func() { <-hostSem }()
			case <-ctx.Done():
				results[i] = Result{URL: u.String(), Err: "context cancelled"}
				return
			}

			reqCtx, cancel := context.WithTimeout(ctx, c.timeout)
			defer cancel()

			req, _ := http.NewRequestWithContext(reqCtx, http.MethodHead, u.String(), nil)
			resp, err := c.client.Do(req)
			if err != nil {
				results[i] = Result{URL: u.String(), Accessible: false, Err: err.Error()}
				return
			}
			defer resp.Body.Close()

			ok := resp.StatusCode >= 200 && resp.StatusCode < 400
			results[i] = Result{URL: u.String(), Accessible: ok, StatusCode: resp.StatusCode}
		}(i, link)
	}
	wg.Wait()
	return results
}
