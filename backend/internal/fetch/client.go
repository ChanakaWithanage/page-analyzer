package fetch

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	hc           *http.Client
	maxBytes     int64
	maxRedirects int
	allowLocal   bool
}

func New(timeout time.Duration, maxRedirects int, maxBytes int64) *Client {
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.ResponseHeaderTimeout = 30 * time.Second

	return &Client{
		hc: &http.Client{
			Timeout:   timeout,
			Transport: tr,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= maxRedirects {
					return http.ErrUseLastResponse
				}
				return nil
			},
		},
		maxBytes:     maxBytes,
		maxRedirects: maxRedirects,
		allowLocal:   false,
	}
}

func (c *Client) AllowLocal() {
	c.allowLocal = true
}

var ErrPrivateAddr = errors.New("refusing to fetch private address")

func (c *Client) guard(u *url.URL) error {
	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("unsupported scheme")
	}
	addrs, err := net.DefaultResolver.LookupIPAddr(context.Background(), u.Hostname())
	if err != nil {
		return err
	}
	for _, a := range addrs {
		ip := a.IP
		if !c.allowLocal && (ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsMulticast()) {
			slog.Warn("blocked fetch to private/loopback address", "host", u.Hostname(), "ip", ip.String())
			return ErrPrivateAddr
		}
	}
	return nil
}

func (c *Client) Get(ctx context.Context, raw string) (*http.Response, io.ReadCloser, error) {
	u, err := url.Parse(raw)
	if err != nil {
		slog.Warn("invalid URL parse", "url", raw, "err", err)
		return nil, nil, err
	}
	if err := c.guard(u); err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		slog.Error("failed to create request", "url", raw, "err", err)
		return nil, nil, err
	}
	req.Header.Set("User-Agent", "GoPageAnalyzer/1.0")

	slog.Debug("fetching URL", "url", raw)
	resp, err := c.hc.Do(req)
	if err != nil {
		slog.Error("fetch failed", "url", raw, "err", err)
		return nil, nil, err
	}
	slog.Debug("fetch completed", "url", raw, "status", resp.StatusCode)

	lr := &io.LimitedReader{R: resp.Body, N: c.maxBytes}
	return resp, io.NopCloser(lr), nil
}
