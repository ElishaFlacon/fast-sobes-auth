package prometheus

import (
	"context"
	"net/http"
	"time"

	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
	def "github.com/ElishaFlacon/fast-sobes-auth/internal/domain"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var _ def.PrometheusServer = (*prometheus)(nil)

type prometheus struct {
	server *http.Server
	log    domain.Logger
}

func NewPrometheus(addr string, log domain.Logger) *prometheus {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	s := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return &prometheus{
		server: s,
		log:    log,
	}
}

func (m *prometheus) Run() {
	go func() {
		if err := m.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			m.log.Fatal("prometheus server failed: %v", err)
		}
	}()
}

func (m *prometheus) Stop(ctx context.Context) {
	if err := m.server.Shutdown(ctx); err != nil {
		m.log.Errorf("Failed to gracefully shutdown prometheus server: %v", err)
	}
}
