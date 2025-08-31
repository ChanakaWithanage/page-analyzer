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
	for i := 1; i <= 6; i++ {
		tag := "h" + strconv.Itoa(i)
		h[tag] = doc.Find(tag).Length()
	}

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

	login := false
	doc.Find("form").EachWithBreak(func(_ int, f *goquery.Selection) bool {
		if f.Find(`input[type="password"]`).Length() > 0 {
			login = true
			return false
		}
		if f.Find(`input[name*="password" i]`).Length() > 0 {
			login = true
			return false
		}
		return true
	})

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
