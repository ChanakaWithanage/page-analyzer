package gateway

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/chanaka-withanage/page-analyzer/pkg/contract"
)

func (s *server) analyze(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	if r.Method != http.MethodPost {
		slog.Warn("invalid method on /analyze", "method", r.Method)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.Warn("invalid JSON payload", "err", err)
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	raw := strings.TrimSpace(body.URL)
	if raw == "" {
		slog.Warn("missing url field in request")
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	u, err := url.Parse(raw)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		slog.Warn("invalid URL submitted", "url", raw, "err", err)
		http.Error(w, "please provide a valid http(s) URL", http.StatusBadRequest)
		return
	}

	slog.Info("starting analysis", "url", u.String())

	res, err := s.svc.Analyze(r.Context(), contract.AnalyzeParams{
		URL: u.String(),
	})

	status := http.StatusOK
	if err != nil {
		slog.Error("analysis failed", "url", u.String(), "err", err, "duration_ms", time.Since(start).Milliseconds())
		status = http.StatusBadGateway
	} else {
		slog.Info("analysis succeeded", "url", u.String(), "duration_ms", time.Since(start).Milliseconds())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		slog.Error("failed to encode response", "url", u.String(), "err", err)
	}
}
