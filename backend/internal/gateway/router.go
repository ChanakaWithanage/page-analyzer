package gateway

import (
	"log/slog"
	"net/http"
	"time"
	"fmt"

	"github.com/chanaka-withanage/page-analyzer/internal/analyzer"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type server struct {
	svc *analyzer.Service
}

func NewMuxWithService(svc *analyzer.Service) http.Handler {
	s := &server{svc: svc}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/analyze", s.analyze)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.Handle("/metrics", promhttp.Handler())

	return withCORS(withLogging(mux))
}

type responseRecorder struct {
	http.ResponseWriter
	status int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.status = code
	rr.ResponseWriter.WriteHeader(code)
}

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        rec := &responseRecorder{ResponseWriter: w, status: 200}

        next.ServeHTTP(rec, r)

        duration := time.Since(start).Seconds()
        slog.Info("request completed",
        			"method", r.Method,
        			"path", r.URL.Path,
        			"status", rec.status,
        			"duration_ms", time.Since(start).Milliseconds(),
        			"remote_addr", r.RemoteAddr,
        		)

        // Prometheus metrics
        requestsTotal.WithLabelValues(r.URL.Path, r.Method, fmt.Sprint(rec.status)).Inc()
        requestDuration.WithLabelValues(r.URL.Path).Observe(duration)
    })
}
