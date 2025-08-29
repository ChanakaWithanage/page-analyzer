package gateway

import (
	"log"
	"net/http"
	"time"

	"github.com/chanaka-withanage/page-analyzer/internal/analyzer"
)

type server struct {
	svc *analyzer.Service
}

func NewMuxWithService(svc *analyzer.Service) http.Handler {
	s := &server{svc: svc}
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.index)
	mux.HandleFunc("/analyze", s.analyze)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	return withLogging(mux)
}

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
