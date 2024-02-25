package internalhttp

import (
	"context"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/model"
	"github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	ReadHeaderTimeout = time.Second * 5
)

type Server struct {
	addr          string
	grpcEndpoint  string
	logger        Logger
	eventsService EventsService
	server        *http.Server
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

func NewServer(addr, grpcEndpoint string, logger Logger, eventsService EventsService) *Server {
	return &Server{
		addr:          addr,
		logger:        logger,
		eventsService: eventsService,
		grpcEndpoint:  grpcEndpoint,
	}
}

func (s *Server) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// mux := runtime.NewServeMux()

	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				EmitDefaultValues: false,
				EmitUnpopulated:   true,
				// UseProtoNames: true,
			},
		}),
	)

	conn, err := grpc.DialContext(
		context.Background(),
		s.grpcEndpoint,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	err = pb.RegisterEventsServiceHandler(ctx, mux, conn)
	if err != nil {
		return err
	}

	s.server = &http.Server{
		Addr:              s.addr,
		Handler:           loggingMiddleware(mux, s.logger),
		ReadHeaderTimeout: ReadHeaderTimeout,
	}

	s.logger.Debugf("starting http on %s", s.addr)
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
