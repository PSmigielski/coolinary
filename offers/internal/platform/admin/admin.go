// Package admin serves /metrics and /healthz on the admin port (8081).
package admin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Check is one readiness dependency (database, redis, ...).
type Check struct {
	Name string
	Fn   func(ctx context.Context) error
}

type Server struct{ srv *http.Server }

func NewServer(port int, checks []Check, log *slog.Logger) *Server {
	mux := http.NewServeMux()
	mux.Handle("GET /metrics", promhttp.Handler())
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		status, ok := map[string]string{}, true
		for _, c := range checks {
			if err := c.Fn(ctx); err != nil {
				status[c.Name], ok = err.Error(), false
			} else {
				status[c.Name] = "ok"
			}
		}
		w.Header().Set("Content-Type", "application/json")
		if !ok {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": ok, "checks": status})
	})
	return &Server{srv: &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: mux, ReadHeaderTimeout: 5 * time.Second}}
}

func (s *Server) Serve() error {
	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error { return s.srv.Shutdown(ctx) }

// Probe is used by the `healthcheck` subcommand (Docker HEALTHCHECK in distroless images).
func Probe(port int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("http://127.0.0.1:%d/healthz", port), http.NoBody)
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("healthz returned %d", res.StatusCode)
	}
	return nil
}
