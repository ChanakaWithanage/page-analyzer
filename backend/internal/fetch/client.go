package fetch

import (
	"context"
	"errors"
	"net"
	"io"          // 👈 add this
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	hc           *http.Client
	maxBytes     int64
	maxRedirects int
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
	}
}

var ErrPrivateAddr = errors.New("refusing to fetch private address")

//SSRF guard
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
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsMulticast() {
			return ErrPrivateAddr
		}
	}
	return nil
}

func (c *Client) Get(ctx context.Context, raw string) (*http.Response, io.ReadCloser, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, nil, err
	}
	if err := c.guard(u); err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("User-Agent", "GoPageAnalyzer/1.0")

    resp, err := c.hc.Do(req)
    if err != nil {
        return nil, nil, err
    }
    // Wrap body to enforce a hard cap (defense vs. huge pages)
    lr := &io.LimitedReader{R: resp.Body, N: c.maxBytes}
    return resp, io.NopCloser(lr), nil

}
