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
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var body struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.Warn("invalid JSON payload", "err", err)
		writeError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}

	raw := strings.TrimSpace(body.URL)
	if raw == "" {
		slog.Warn("missing url field in request")
		writeError(w, http.StatusBadRequest, "url is required")
		return
	}

	// regex validation
	if !isValidURL(raw) {
		slog.Warn("url failed regex validation", "url", raw)
		writeError(w, http.StatusBadRequest, "please provide a valid http(s) URL")
		return
	}

	// strict parsing
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" {
		slog.Warn("invalid URL after parse", "url", raw, "err", err)
		writeError(w, http.StatusBadRequest, "please provide a valid http(s) URL")
		return
	}

	slog.Info("starting analysis", "url", u.String())

	res, err := s.svc.Analyze(r.Context(), contract.AnalyzeParams{
		URL: u.String(),
	})

	status := http.StatusOK
	if err != nil {
		slog.Error("analysis failed",
			"url", u.String(),
			"err", err,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		status = http.StatusBadGateway
	} else {
		slog.Info("analysis succeeded",
			"url", u.String(),
			"duration_ms", time.Since(start).Milliseconds(),
		)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		slog.Error("failed to encode response", "url", u.String(), "err", err)
	}
}
