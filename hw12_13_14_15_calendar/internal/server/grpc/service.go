package internalgrpc

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/model"
	"github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/server/grpc/pb"
)

type Service struct {
	pb.UnimplementedEventsServiceServer
	eventsService EventsService
}

func NewService(eventsService EventsService) *Service {
	return &Service{
		eventsService: eventsService,
	}
}

func (s *Service) Health(context.Context, *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func (s *Service) CreateEvent(ctx context.Context, e *pb.Event) (*pb.CreateEventResponse, error) {
	id, err := s.eventsService.CreateEvent(ctx, model.Event{
		Title:           e.Title,
		Description:     e.Description,
		UserId:          int(e.UserId),
		NotifyBeforeMin: int(e.NotifyBeforeMin),
		StartDatetime:   e.StartDatetime.AsTime(),
		EndDatetime:     e.EndDatetime.AsTime(),
	})
	return &pb.CreateEventResponse{
		Id: int32(id),
	}, err
}

func (s *Service) ModifyEvent(ctx context.Context, e *pb.Event) (*empty.Empty, error) {
	err := s.eventsService.ModifyEvent(ctx, model.Event{
		Id:              int(e.Id),
		Title:           e.Title,
		Description:     e.Description,
		UserId:          int(e.UserId),
		NotifyBeforeMin: int(e.NotifyBeforeMin),
		StartDatetime:   e.StartDatetime.AsTime(),
		EndDatetime:     e.EndDatetime.AsTime(),
	})
	return &emptypb.Empty{}, err
}

func (s *Service) RemoveEvent(ctx context.Context, e *pb.EventIdRequest) (*empty.Empty, error) {
	err := s.eventsService.RemoveEvent(ctx, int(e.Id))
	return &emptypb.Empty{}, err
}

func (s *Service) GetEvent(ctx context.Context, e *pb.EventIdRequest) (*pb.Event, error) {
	ev, err := s.eventsService.GetEvent(ctx, int(e.Id))
	return &pb.Event{
		Id:              int32(ev.Id),
		Title:           ev.Title,
		Description:     ev.Description,
		StartDatetime:   timestamppb.New(ev.StartDatetime),
		EndDatetime:     timestamppb.New(ev.EndDatetime),
		UserId:          int32(ev.UserId),
		NotifyBeforeMin: int32(ev.NotifyBeforeMin),
	}, err
}

func (s *Service) GetEventsForDay(ctx context.Context, r *pb.DateRequest) (*pb.Events, error) {
	t, err := time.Parse("2006-01-02", r.Date)
	if err != nil {
		return nil, err
	}
	evs, err := s.eventsService.GetEventsForDay(ctx, t)
	if err != nil {
		return nil, err
	}
	return evsToPbEvents(evs), err
}

func (s *Service) GetEventsForWeek(ctx context.Context, r *pb.DateRequest) (*pb.Events, error) {
	t, err := time.Parse("2006-01-02", r.Date)
	if err != nil {
		return nil, err
	}
	evs, err := s.eventsService.GetEventsForWeek(ctx, t)
	if err != nil {
		return nil, err
	}
	return evsToPbEvents(evs), err
}

func (s *Service) GetEventsForMonth(ctx context.Context, r *pb.DateRequest) (*pb.Events, error) {
	t, err := time.Parse("2006-01-02", r.Date)
	if err != nil {
		return nil, err
	}
	evs, err := s.eventsService.GetEventsForMonth(ctx, t)
	if err != nil {
		return nil, err
	}
	return evsToPbEvents(evs), err
}

func evsToPbEvents(evs []model.Event) *pb.Events {
	result := make([]*pb.Event, 0, len(evs))
	for _, e := range evs {
		result = append(result, &pb.Event{
			Id:              int32(e.Id),
			Title:           e.Title,
			Description:     e.Description,
			StartDatetime:   timestamppb.New(e.StartDatetime),
			EndDatetime:     timestamppb.New(e.EndDatetime),
			UserId:          int32(e.UserId),
			NotifyBeforeMin: int32(e.NotifyBeforeMin),
		})
	}
	return &pb.Events{
		Events: result,
	}
}
