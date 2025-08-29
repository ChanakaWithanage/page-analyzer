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

    if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
        if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB max memory
            http.Error(w, "invalid multipart form", http.StatusBadRequest)
            return
        }
    } else {
        if err := r.ParseForm(); err != nil {
            http.Error(w, "invalid form", http.StatusBadRequest)
            return
        }
    }

    raw := strings.TrimSpace(r.FormValue("url"))
    log.Printf("DEBUG: got url=%q", raw)
    if raw == "" {
        http.Error(w, "url is required", http.StatusBadRequest)
        return
    }

	u, err := url.Parse(raw)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		http.Error(w, "please provide a valid http(s) URL", http.StatusBadRequest); return
	}

	res, err := s.svc.Analyze(r.Context(), contract.AnalyzeParams{
		URL:                 u.String(),
		FetchTimeoutSeconds: 10,
	})
	status := http.StatusOK
	if err != nil {
		status = http.StatusBadGateway // upstream issues
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(res)
}
