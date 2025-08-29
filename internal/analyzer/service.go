package analyzer

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/chanaka-withanage/page-analyzer/pkg/contract"
	"github.com/chanaka-withanage/page-analyzer/internal/fetch"
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

	// Read/Discard now; next step we'll parse.
	n, _ := io.Copy(io.Discard, body)
	res.Warnings = append(res.Warnings, fmt.Sprintf("fetched %d bytes, parsing to be added", n))
	return res, nil
}
