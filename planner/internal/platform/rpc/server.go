// Package rpc serves connect-go handlers over h2c. One port answers gRPC
// (for grpc-js / grpcio / grpc-go clients) and the Connect protocol over
// HTTP/1.1 (curl-friendly), plus grpc.health.v1 and, in dev, reflection.
package rpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
)

type Server struct {
	mux     *http.ServeMux
	names   []string
	Health  *grpchealth.StaticChecker
	srv     *http.Server
	log     *slog.Logger
	reflect bool
}

func NewServer(port int, reflection bool, log *slog.Logger) *Server {
	s := &Server{mux: http.NewServeMux(), log: log, reflect: reflection}
	var protocols http.Protocols
	protocols.SetHTTP1(true)
	protocols.SetUnencryptedHTTP2(true)
	s.srv = &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           s.mux,
		Protocols:         &protocols,
		ReadHeaderTimeout: 5 * time.Second,
	}
	return s
}

// Handle mounts a generated connect handler, e.g. kitchenv1connect.NewKitchenServiceHandler.
func (s *Server) Handle(serviceName, path string, h http.Handler) {
	s.names = append(s.names, serviceName)
	s.mux.Handle(path, h)
}

// Options returns the interceptors every handler should use.
func Options(log *slog.Logger) []connect.HandlerOption {
	return []connect.HandlerOption{
		connect.WithInterceptors(NewServerInterceptor(log)),
	}
}

func (s *Server) Serve() error {
	s.Health = grpchealth.NewStaticChecker(s.names...)
	s.mux.Handle(grpchealth.NewHandler(s.Health))
	if s.reflect {
		reflector := grpcreflect.NewStaticReflector(append(s.names, grpchealth.HealthV1ServiceName)...)
		s.mux.Handle(grpcreflect.NewHandlerV1(reflector))
		s.mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
	}
	ln, err := net.Listen("tcp", s.srv.Addr)
	if err != nil {
		return err
	}
	s.log.Info("gRPC server listening", "addr", s.srv.Addr, "services", s.names, "reflection", s.reflect)
	if err := s.srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error { return s.srv.Shutdown(ctx) }
