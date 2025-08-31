package analyzer

import (
	"context"
	"fmt"
	"log/slog"
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
	return &Service{fetch: fetchClient, defaultTimeout: 30 * time.Second}
}

func (s *Service) SetDefaultTimeout(d time.Duration) {
	s.defaultTimeout = d
}

func (s *Service) Analyze(ctx context.Context, p contract.AnalyzeParams) (*contract.AnalyzeResult, error) {
	timeout := s.defaultTimeout
	if p.FetchTimeoutSeconds > 0 {
		timeout = time.Duration(p.FetchTimeoutSeconds) * time.Second
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	start := time.Now()
	slog.Info("analysis started", "url", p.URL, "timeout", timeout)

	res := &contract.AnalyzeResult{
		URL:      p.URL,
		Headings: map[string]int{},
		Warnings: []string{},
		Errors:   []string{},
	}

	resp, body, err := s.fetch.Get(ctx, p.URL)
	if err != nil {
		slog.Error("fetch failed", "url", p.URL, "err", err)
		res.Errors = append(res.Errors, err.Error())
		return res, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		slog.Warn("upstream returned non-2xx", "url", p.URL, "status", resp.StatusCode)
		res.Errors = append(res.Errors, fmt.Sprintf("upstream status: %d", resp.StatusCode))
		return res, fmt.Errorf("upstream returned %d", resp.StatusCode)
	}

	u, _ := url.Parse(p.URL)
	parsed, err := parser.Parse(body, u)
	if err != nil {
		slog.Error("parse failed", "url", p.URL, "err", err)
		res.Errors = append(res.Errors, err.Error())
		return res, err
	}

	res.HTMLVersion = parsed.HTMLVersion
	res.Title = parsed.Title
	res.Headings = parsed.Headings
	res.LoginFormPresent = parsed.LoginFormPresent

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
		slog.Debug("validating links", "url", p.URL, "count", len(urlObjs))
		checker := linkcheck.New(10, 2, s.defaultTimeout/2)
		results := checker.Validate(ctx, urlObjs)

		bad := 0
		for _, r := range results {
			if !r.Accessible {
				bad++
			}
		}
		res.LinksInaccessible = bad
		slog.Info("link validation complete", "url", p.URL, "bad_links", bad)
	}

	slog.Info("analysis finished",
		"url", p.URL,
		"duration_ms", time.Since(start).Milliseconds(),
		"title", res.Title,
		"headings", len(res.Headings),
	)

	return res, nil
}

func sameHost(a, b string) bool {
	return strings.EqualFold(a, b)
}
