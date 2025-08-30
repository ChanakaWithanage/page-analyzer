package gateway

import (
	"encoding/json"
	"html/template"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"log"

	"github.com/chanaka-withanage/page-analyzer/pkg/contract"
)

var indexTpl = template.Must(template.ParseFiles(filepath.Join("web", "templates", "index.tmpl.html")))

func (s *server) index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = indexTpl.Execute(w, nil)
}

func (s *server) analyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Expect JSON body like { "url": "https://example.com" }
	var body struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	raw := strings.TrimSpace(body.URL)
	log.Printf("DEBUG: got url=%q", raw)
	if raw == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	u, err := url.Parse(raw)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		http.Error(w, "please provide a valid http(s) URL", http.StatusBadRequest)
		return
	}

	res, err := s.svc.Analyze(r.Context(), contract.AnalyzeParams{
		URL:                 u.String(),
		FetchTimeoutSeconds: 10,
	})

	status := http.StatusOK
	if err != nil {
		status = http.StatusBadGateway
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(res)
}
