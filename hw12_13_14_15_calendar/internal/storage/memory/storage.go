package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/errors"
	"github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/model"
)

type Storage struct {
	mu     sync.RWMutex
	lastId int
	events []model.Event
}

func NewStorage() *Storage {
	return &Storage{
		events: make([]model.Event, 0, 100),
	}
}

func (*Storage) Connect(_ context.Context) error {
	return nil
}

func (*Storage) Close() error {
	return nil
}

func (s *Storage) InsertEvent(_ context.Context, event model.Event) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.lastId += 1
	event.Id = s.lastId
	s.events = append(s.events, event)

	return event.Id, nil
}

func (s *Storage) UpdateEvent(_ context.Context, event model.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	i := indexById(s.events, event.Id)
	if i < 0 {
		return errors.ErrEventNotFound
	}

	s.events[i] = event

	return nil
}

func (s *Storage) DeleteEvent(_ context.Context, eventId int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	i := indexById(s.events, eventId)
	if i < 0 {
		return errors.ErrEventNotFound
	}

	s.events[i] = s.events[len(s.events)-1]
	s.events = s.events[:len(s.events)-1]

	return nil
}

func (s *Storage) GetEvent(_ context.Context, eventId int) (model.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	i := indexById(s.events, eventId)
	if i < 0 {
		return model.Event{}, errors.ErrEventNotFound
	}

	return s.events[i], nil
}

func (s *Storage) GetEvents(_ context.Context, dateFrom, dateTo time.Time) ([]model.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]model.Event, 0)
	for _, v := range s.events {
		if (v.StartDatetime.After(dateFrom) && v.StartDatetime.Before(dateTo)) ||
			(v.EndDatetime.After(dateFrom) && v.EndDatetime.Before(dateTo)) {
			result = append(result, v)
		}
	}
	return result, nil
}

// check time cross
func (s *Storage) IsTimeFree(_ context.Context, excludeId int, timeFrom, timeTo time.Time) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, v := range s.events {
		if v.Id == excludeId {
			continue
		}
		// start <= 13 && end >= 14
		if (v.StartDatetime.Compare(timeFrom) <= 0 && v.EndDatetime.Compare(timeTo) >= 0) ||
			// start >= 13 && <= 14
			(v.StartDatetime.Compare(timeFrom) > 0 && v.StartDatetime.Compare(timeTo) < 0) ||
			// end between > 13 && <= 14
			(v.EndDatetime.Compare(timeFrom) > 0 && v.EndDatetime.Compare(timeTo) < 0) {
			return false, nil
		}
	}
	return true, nil
}

func (s *Storage) DeleteEventsOlderThan(_ context.Context, date time.Time) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rowsAffected := 0
	for i := 0; i < len(s.events); i++ {
		if s.events[i].StartDatetime.Compare(date) < 0 {
			s.events[i] = s.events[len(s.events)-1]
			s.events = s.events[:len(s.events)-1]
			rowsAffected++
		}
	}

	return int64(rowsAffected), nil
}

func indexById(src []model.Event, eventId int) int {
	// return slices.IndexFunc(src, func(el model.Event) bool { return v.Id == eventId })
	for i, v := range src {
		if v.Id == eventId {
			return i
		}
	}
	return -1
}
