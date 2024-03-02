package internalgrpc

//go:generate protoc -I ../../../api --go_out=../../../ --go-grpc_out=../../../ --grpc-gateway_out=../../../ --openapiv2_out=../../../ ../../../api/eventservice/eventservice.proto

import (
	"context"
	"net"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	"github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/model"
	"github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/server/grpc/pb"
)

type Server struct {
	addr          string
	logger        Logger
	eventsService EventsService
	server        *grpc.Server
}

type Logger interface {
	Debugf(msg string, args ...interface{})
	Infof(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
	Fatalf(msg string, args ...interface{})
}

type EventsService interface {
	CreateEvent(ctx context.Context, event model.Event) (int, error)
	ModifyEvent(ctx context.Context, event model.Event) error
	RemoveEvent(ctx context.Context, eventId int) error
	GetEvent(ctx context.Context, eventId int) (model.Event, error)
	GetEventsForDay(ctx context.Context, date time.Time) ([]model.Event, error)
	GetEventsForWeek(ctx context.Context, date time.Time) ([]model.Event, error)
	GetEventsForMonth(ctx context.Context, date time.Time) ([]model.Event, error)
}

func NewServer(addr string, logger Logger, eventsService EventsService) *Server {
	return &Server{
		addr:          addr,
		logger:        logger,
		eventsService: eventsService,
	}
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	s.server = grpc.NewServer(grpc.ChainUnaryInterceptor(s.loggerInterceptor()))
	pb.RegisterEventsServiceServer(s.server, NewService(s.eventsService))

	s.logger.Debugf("starting grpc on %s", s.addr)

	return s.server.Serve(lis)
}

func (s *Server) Stop() {
	s.server.GracefulStop()
}

func (s *Server) loggerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		startTime := time.Now()

		meta, ok := metadata.FromIncomingContext(ctx)
		agent := ""
		if ok {
			agent = strings.Join(meta.Get("user-agent"), " ")
		}

		resp, err := handler(ctx, req)
		// Логирование ошибки, если она произошла
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}

		// Получение RemoteAddr из контекста
		p, ok := peer.FromContext(ctx)
		var remoteAddr string
		if ok {
			remoteAddr = p.Addr.String()
		}

		s.logger.Debugf("addr: %s method: %s duration: %v user-agent: [%s] err: '%v'",
			remoteAddr, info.FullMethod, time.Since(startTime), agent, errMsg)

		return resp, err
	}
}
