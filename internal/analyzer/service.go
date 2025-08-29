package analyzer

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/chanaka-withanage/page-analyzer/internal/fetch"
	"github.com/chanaka-withanage/page-analyzer/internal/parser"
	"github.com/chanaka-withanage/page-analyzer/pkg/contract"
	"github.com/chanaka-withanage/page-analyzer/internal/linkcheck"
)

type Service struct {
	fetch *fetch.Client
}

func New(fetchClient *fetch.Client) *Service {
	return &Service{fetch: fetchClient}
}

func (s *Service) Analyze(ctx context.Context, p contract.AnalyzeParams) (*contract.AnalyzeResult, error) {

	ctx, cancel := context.WithTimeout(ctx, time.Duration(p.FetchTimeoutSeconds)*time.Second)
	defer cancel()

	res := &contract.AnalyzeResult{
		URL:       p.URL,
		Headings:  map[string]int{},
		Warnings:  []string{},
		Errors:    []string{},
	}

	resp, body, err := s.fetch.Get(ctx, p.URL)
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
		return res, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		res.Errors = append(res.Errors, fmt.Sprintf("upstream status: %d", resp.StatusCode))
		return res, fmt.Errorf("upstream returned %d", resp.StatusCode)
	}

	// parse HTML
	u, _ := url.Parse(p.URL)
	parsed, err := parser.Parse(body, u)
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
		return res, err
	}

	// fill results
	res.HTMLVersion = parsed.HTMLVersion
	res.Title = parsed.Title
	res.Headings = parsed.Headings
	res.LoginFormPresent = parsed.LoginFormPresent

	// classify links
	host := u.Host
    var urlObjs []*url.URL
    for _, l := range parsed.Links {
    	lu, err := url.Parse(l)
    	if err != nil {
    		continue
    	}
    	urlObjs = append(urlObjs, lu)
    	if sameHost(host, lu.Host) {
    		res.LinksInternal++
    	} else {
    		res.LinksExternal++
    	}
    }

    if len(urlObjs) > 0 {
    	checker := linkcheck.New(10, 2, 3*time.Second) // tune: 10 global, 2 per-host, 3s timeout
    	results := checker.Validate(ctx, urlObjs)

    	bad := 0
    	for _, r := range results {
    		if !r.Accessible {
    			bad++
    		}
    	}
    	res.LinksInaccessible = bad
    }

	return res, nil
}

func sameHost(a, b string) bool {
	return strings.EqualFold(a, b)
}
