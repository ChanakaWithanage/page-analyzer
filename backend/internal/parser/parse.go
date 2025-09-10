package parser

import (
	"io"
	"log/slog"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type Parsed struct {
	HTMLVersion      string
	Title            string
	Headings         map[string]int
	Links            []string
	LoginFormPresent bool
}

func Parse(r io.Reader, base *url.URL) (*Parsed, error) {
	root, err := html.Parse(r)
	if err != nil {
		slog.Error("failed to parse HTML", "base_url", base.String(), "err", err)
		return nil, err
	}
	doc := goquery.NewDocumentFromNode(root)

	ver := detectDoctype(root)
	title := strings.TrimSpace(doc.Find("title").First().Text())

	h := map[string]int{}
	doc.Find("h1,h2,h3,h4,h5,h6").Each(func(_ int, s *goquery.Selection) {
		if n := s.Nodes; len(n) > 0 {
			tag := strings.ToLower(n[0].Data)
			if tag == "h1" || tag == "h2" || tag == "h3" || tag == "h4" || tag == "h5" || tag == "h6" {
				h[tag]++
			}
		}
	})
	for i := 1; i <= 6; i++ {
		key := "h" + strconv.Itoa(i)
		if _, ok := h[key]; !ok {
			h[key] = 0
		}
	}

	// Links
	var links []string
	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if !ok || href == "" {
			return
		}
		u, err := base.Parse(href)
		if err == nil && u.Scheme != "" && u.Host != "" {
			links = append(links, u.String())
		}
	})

	// Login detection
	login := hasObviousLoginForm(doc) || hasSimpleAuthCTA(doc)

	parsed := &Parsed{
		HTMLVersion:      ver,
		Title:            title,
		Headings:         h,
		Links:            links,
		LoginFormPresent: login,
	}

	slog.Debug("parsed HTML successfully",
		"base_url", base.String(),
		"title", parsed.Title,
		"headings_total", len(parsed.Headings),
		"links_total", len(parsed.Links),
		"login_form_present", parsed.LoginFormPresent,
	)

	return parsed, nil
}

func detectDoctype(n *html.Node) string {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.DoctypeNode {
			name := strings.ToLower(c.Data)
			if name == "html" {
				if c.Attr == nil || len(c.Attr) == 0 {
					return "HTML5"
				}
				data := strings.ToLower(c.Data)
				if strings.Contains(data, "xhtml 1.0") {
					return "XHTML 1.0"
				}
				if strings.Contains(data, "xhtml 1.1") {
					return "XHTML 1.1"
				}
				return "HTML 4.x"
			}
		}
	}
	return "unknown"
}


func hasObviousLoginForm(doc *goquery.Document) bool {
	// Common substrings hinting at auth fields
	hints := []string{
		"login", "signin", "sign-in",
		"username", "user", "email", "e-mail",
		"pwd", "password", "passcode",
	}
	hasHint := func(v string) bool {
		v = strings.ToLower(v)
		for _, h := range hints {
			if strings.Contains(v, h) {
				return true
			}
		}
		return false
	}

	found := false
	doc.Find("form").EachWithBreak(func(_ int, f *goquery.Selection) bool {
		inputs := f.Find("input")

		if inputs.Filter(`[type="password"]`).Length() > 0 {
			found = true
			return false
		}

		authField := false
		inputs.EachWithBreak(func(_ int, inp *goquery.Selection) bool {
			for _, a := range []string{"name", "id", "autocomplete", "placeholder"} {
				if v, ok := inp.Attr(a); ok && hasHint(v) {
					authField = true
					return false
				}
			}
			return true
		})
		if !authField {
			return true
		}

		if f.Find(`button, input[type="submit"], input[type="button"]`).Length() > 0 {
			found = true
			return false
		}

		return true
	})
	return found
}

func hasSimpleAuthCTA(doc *goquery.Document) bool {
	match := func(s string) bool {
		s = strings.ToLower(strings.TrimSpace(strings.Join(strings.Fields(s), " ")))
		return strings.Contains(s, "log in") ||
			strings.Contains(s, "signin") ||
			strings.Contains(s, "sign in") ||
			strings.Contains(s, "account")
	}

	found := false
	doc.Find("button, a").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		if match(s.Text()) {
			found = true
			return false
		}
		return true
	})
	return found
}
