package parser

import (
	"io"
	"net/url"
	"strings"
	"strconv"

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

	// parse into an html.Node tree
	root, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	doc := goquery.NewDocumentFromNode(root)

	// --- HTML version
	ver := detectDoctype(root)

	// --- Title
	title := strings.TrimSpace(doc.Find("title").First().Text())

	// --- Headings h1..h6
	h := map[string]int{}
	for i := 1; i <= 6; i++ {
        tag := "h" + strconv.Itoa(i) // ✅ fixed
        h[tag] = doc.Find(tag).Length()
    }

	// --- Links
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

	// --- Login form detection
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

	return &Parsed{
		HTMLVersion:      ver,
		Title:            title,
		Headings:         h,
		Links:            links,
		LoginFormPresent: login,
	}, nil
}

func detectDoctype(n *html.Node) string {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.DoctypeNode {
			name := strings.ToLower(c.Data)
			if name == "html" {
				if c.Attr == nil || len(c.Attr) == 0 {
					return "HTML5"
				}
				// crude heuristic: legacy doctypes
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
