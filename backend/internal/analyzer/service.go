package analyzer

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/chanaka-withanage/page-analyzer/internal/fetch"
	"github.com/chanaka-withanage/page-analyzer/internal/linkcheck"
	"github.com/chanaka-withanage/page-analyzer/internal/parser"
	"github.com/chanaka-withanage/page-analyzer/pkg/contract"
)

type Service struct {
	fetch          *fetch.Client
	defaultTimeout time.Duration
}

func New(fetchClient *fetch.Client) *Service {
	// fallback to 30s if not overridden
	return &Service{fetch: fetchClient, defaultTimeout: 30 * time.Second}
}

func (s *Service) SetDefaultTimeout(d time.Duration) {
	s.defaultTimeout = d
}

func (s *Service) Analyze(ctx context.Context, p contract.AnalyzeParams) (*contract.AnalyzeResult, error) {
	// pick timeout: explicit param > default
	timeout := s.defaultTimeout
	if p.FetchTimeoutSeconds > 0 {
		timeout = time.Duration(p.FetchTimeoutSeconds) * time.Second
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	res := &contract.AnalyzeResult{
		URL:      p.URL,
		Headings: map[string]int{},
		Warnings: []string{},
		Errors:   []string{},
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

	// validate links concurrently
	if len(urlObjs) > 0 {
		checker := linkcheck.New(10, 2, s.defaultTimeout/2) // use half of default for per-link validation
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
