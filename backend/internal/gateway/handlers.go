package gateway

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/chanaka-withanage/page-analyzer/pkg/contract"
)

func (s *server) analyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
		URL: u.String(),
	})

	status := http.StatusOK
	if err != nil {
		status = http.StatusBadGateway
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(res)
}
