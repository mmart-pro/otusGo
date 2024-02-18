package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/model"
)

const (
	ReadHeaderTimeout = time.Second * 5
)

type Server struct {
	addr          Endpointer
	logger        Logger
	eventsService EventsService
	server        *http.Server
}

type Endpointer interface {
	GetEndpoint() string
}

type Logger interface {
	Debugf(msg string, args ...interface{})
	Infof(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
	Fatalf(msg string, args ...interface{})
}

type EventsService interface {
	CreateEvent(ctx context.Context, event model.Event) (int, error)
}

func NewServer(addr Endpointer, logger Logger, eventsService EventsService) *Server {
	return &Server{
		addr:          addr,
		logger:        logger,
		eventsService: eventsService,
	}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", hello)

	addr := s.addr.GetEndpoint()
	s.server = &http.Server{
		Addr:              addr,
		Handler:           loggingMiddleware(mux, s.logger),
		ReadHeaderTimeout: ReadHeaderTimeout,
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
