package server

import (
	"context"
	"net/http"

	"github.com/puppetlabs/insights-instrumentation/metrics"
)

// Options are the Server configuration options
type Options struct {
	// BindAddr is the address:port to listen on
	BindAddr string

	// Path is the URI path to handle requests for.
	// Default is /
	Path string
}

// Server delegates http requests for metrics on a configured path to the Metrics
// collector.
type Server struct {
	bindAddr string
	m        *metrics.Metrics
	path     string
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != s.path {
		http.NotFound(w, r)

		return
	}

	handler := NewHandler(s.m)
	handler.ServeHTTP(w, r)
}

func (s Server) Run(ctx context.Context) error {
	hs := &http.Server{Addr: s.bindAddr, Handler: s}

	go func() {
		<-ctx.Done()
		hs.Shutdown(ctx)
	}()

	if err := hs.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

// New returns a new Server
func New(m *metrics.Metrics, opts Options) *Server {
	path := opts.Path

	if path == "" {
		path = "/"
	}

	return &Server{
		bindAddr: opts.BindAddr,
		m:        m,
		path:     path,
	}
}
