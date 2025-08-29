package gateway

import (
	"encoding/json"
	"html/template"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"log"
)

var indexTpl = template.Must(template.ParseFiles(filepath.Join("web", "templates", "index.tmpl.html")))

func index(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = indexTpl.Execute(w, nil)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func analyze(w http.ResponseWriter, r *http.Request) {
    log.Printf("analyze11: raw url form value = %q", r.FormValue("url"))
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	raw := strings.TrimSpace(r.FormValue("url"))
	if raw == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

    log.Printf("analyze22: raw url form value = %q", r.FormValue("url"))
	u, err := url.Parse(raw)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		http.Error(w, "please provide a valid http(s) URL", http.StatusBadRequest)
		return
	}

	resp := map[string]any{
		"url":          u.String(),
		"html_version": "unknown",
		"title":        "",
		"headings":     map[string]int{},
		"links": map[string]int{
			"internal":     0,
			"external":     0,
			"inaccessible": 0,
		},
		"login_form_present": false,
		"warnings":           []string{},
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
