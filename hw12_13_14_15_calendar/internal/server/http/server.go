package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	READ_HEADER_TIMEOUT = time.Second * 5
)

type Server struct {
	addr   Endpointer
	logger Logger
	app    Application
	server *http.Server
}

type Logger interface {
	Debugf(msg string, args ...interface{})
	Infof(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
	Fatalf(msg string, args ...interface{})
}

type Application interface {
	CreateEvent(ctx context.Context, id, title string) error // temp
}

type Endpointer interface {
	GetEndpoint() string
}

func NewServer(addr Endpointer, logger Logger, app Application) *Server {
	return &Server{
		addr:   addr,
		logger: logger,
		app:    app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", hello)

	addr := s.addr.GetEndpoint()
	s.server = &http.Server{
		Addr:              addr,
		Handler:           loggingMiddleware(mux, s.logger),
		ReadHeaderTimeout: READ_HEADER_TIMEOUT,
	}

	s.logger.Debugf("starting http on %s", addr)
	err := s.server.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}
