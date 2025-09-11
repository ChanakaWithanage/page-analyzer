package gateway

import (
	"errors"
	"net"
	"net/url"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	errEmpty     = errors.New("url is required")
	errTooLong   = errors.New("url is too long")
	errCtrlChars = errors.New("url contains control characters")
	errScheme    = errors.New("only http/https schemes are allowed")
	errParse     = errors.New("invalid URL")
	errHost      = errors.New("invalid host")
	errUserInfo  = errors.New("userinfo not allowed")
	errPort      = errors.New("invalid port")

	maxURLLength = 2048
)

func normalizeAndValidateURL(raw string) (*url.URL, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, errEmpty
	}
	if len(raw) > maxURLLength {
		return nil, errTooLong
	}
	if hasControlChars(raw) {
		return nil, errCtrlChars
	}

	if !strings.Contains(raw, "://") && looksLikeHost(raw) {
		raw = "https://" + raw
	}

	u, err := url.Parse(raw)
	if err != nil {
		return nil, errParse
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, errScheme
	}
	if u.User != nil {
		return nil, errUserInfo
	}

	host := u.Hostname()
	if host == "" {
		return nil, errHost
	}

	// Validate port if present
	if p := u.Port(); p != "" {
		if !isValidPort(p) {
			return nil, errPort
		}
	}

	return u, nil
}

func hasControlChars(s string) bool {
	for _, r := range s {
		if r < 0x20 || r == 0x7f { // C0/DEL
			return true
		}
	}
	return false
}

func looksLikeHost(s string) bool {
	if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
		return true
	}
	if net.ParseIP(s) != nil {
		return true
	}
	return strings.Contains(s, ".") || strings.EqualFold(s, "localhost")
}

func isValidPort(p string) bool {
	n := 0
	for i := 0; i < len(p); {
		r, sz := utf8.DecodeRuneInString(p[i:])
		if !unicode.IsDigit(r) {
			return false
		}
		n = n*10 + int(r-'0')
		if n > 65535 {
			return false
		}
		i += sz
	}
	return n >= 1
}
