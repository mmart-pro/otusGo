package events

import (
	"context"
	"time"

	"github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/errors"
	"github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/model"
)

type EventsService struct {
	logger        Logger
	eventsStorage EventsStorage
}

type Logger interface {
	Debugf(msg string, args ...interface{})
	Infof(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
	Fatalf(msg string, args ...interface{})
}

type EventsStorage interface {
	InsertEvent(ctx context.Context, event model.Event) (int, error)
	UpdateEvent(ctx context.Context, event model.Event) error
	DeleteEvent(ctx context.Context, eventId int) error
	GetEvent(ctx context.Context, eventId int) (model.Event, error)

	GetEvents(ctx context.Context, dateFrom, dateTo time.Time) ([]model.Event, error)
	IsTimeFree(ctx context.Context, excludeId int, dateFrom, dateTo time.Time) (bool, error)
}

func NewEventsService(logger Logger, eventsStorage EventsStorage) *EventsService {
	return &EventsService{
		logger:        logger,
		eventsStorage: eventsStorage,
	}
}

func (s *EventsService) CreateEvent(ctx context.Context, event model.Event) (int, error) {
	if !event.ValidPeriod() {
		return -1, errors.ErrIncorrectDates
	}

	// check that destination time is free
	if free, err := s.eventsStorage.IsTimeFree(ctx, 0, event.StartDatetime, event.EndDatetime); err != nil {
		return -1, err
	} else if !free {
		return -1, errors.ErrDateTimeBusy
	}

	return s.eventsStorage.InsertEvent(ctx, event)
}

func (s *EventsService) ModifyEvent(ctx context.Context, event model.Event) error {
	if !event.ValidPeriod() {
		return errors.ErrIncorrectDates
	}

	// check that destination time is free
	if free, err := s.eventsStorage.IsTimeFree(ctx, event.Id, event.StartDatetime, event.EndDatetime); err != nil {
		return err
	} else if !free {
		return errors.ErrDateTimeBusy
	}

	return s.eventsStorage.UpdateEvent(ctx, event)
}

func (s *EventsService) RemoveEvent(ctx context.Context, eventId int) error {
	return s.eventsStorage.DeleteEvent(ctx, eventId)
}

func (s *EventsService) GetEvent(ctx context.Context, eventId int) (model.Event, error) {
	return s.eventsStorage.GetEvent(ctx, eventId)
}

func (s *EventsService) GetEventsForDay(ctx context.Context, date time.Time) ([]model.Event, error) {
	dateTo := date.AddDate(0, 0, 1)
	return s.eventsStorage.GetEvents(ctx, date, dateTo)
}

func (s *EventsService) GetEventsForWeek(ctx context.Context, date time.Time) ([]model.Event, error) {
	dateTo := date.AddDate(0, 0, 7)
	return s.eventsStorage.GetEvents(ctx, date, dateTo)
}

func (s *EventsService) GetEventsForMonth(ctx context.Context, date time.Time) ([]model.Event, error) {
	dateTo := date.AddDate(0, 1, 0)
	return s.eventsStorage.GetEvents(ctx, date, dateTo)
}
