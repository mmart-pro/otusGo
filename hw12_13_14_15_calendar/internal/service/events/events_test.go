package events

import (
	"context"
	"testing"
	"time"

	"github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/errors"
	"github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/model"
	memorystorage "github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debugf(msg string, args ...interface{}) {}
func (m *MockLogger) Infof(msg string, args ...interface{})  {}
func (m *MockLogger) Errorf(msg string, args ...interface{}) {}
func (m *MockLogger) Fatalf(msg string, args ...interface{}) {}

type MockEventsStorage struct {
	mock.Mock
}

func TestCreateEvent(t *testing.T) {
	service := NewEventsService(new(MockLogger), memorystorage.NewStorage())
	ctx := context.Background()

	event := model.Event{
		StartDatetime: time.Date(2024, 1, 5, 20, 0, 0, 0, time.Local),
		EndDatetime:   time.Date(2024, 1, 5, 20, 15, 0, 0, time.Local),
	}

	id, err := service.CreateEvent(ctx, event)
	event.Id = id

	require.NoError(t, err)
	require.Equal(t, 1, id)

	// Проверка, что событие было добавлено в хранилище
	storedEvent, err := service.GetEvent(ctx, id)
	require.NoError(t, err)
	require.Equal(t, event, storedEvent)

	// Проверка ошибки пересечений
	_, err = service.CreateEvent(ctx, event)
	require.Error(t, err)
	require.ErrorIs(t, err, errors.ErrDateTimeBusy)

	// Проверка ошибки некорректного периода
	_, err = service.CreateEvent(ctx, model.Event{
		StartDatetime: time.Date(2024, 1, 5, 20, 15, 0, 0, time.Local),
		EndDatetime:   time.Date(2024, 1, 5, 20, 0, 0, 0, time.Local),
	})
	require.Error(t, err)
	require.ErrorIs(t, err, errors.ErrIncorrectDates)
}

func TestModifyEvent(t *testing.T) {
	service := NewEventsService(new(MockLogger), memorystorage.NewStorage())
	ctx := context.Background()

	event := model.Event{
		StartDatetime: time.Date(2024, 1, 5, 20, 0, 0, 0, time.Local),
		EndDatetime:   time.Date(2024, 1, 5, 20, 15, 0, 0, time.Local),
	}

	id, err := service.CreateEvent(ctx, event)
	event.Id = id
	require.NoError(t, err)

	updated := event
	updated.EndDatetime = time.Date(2024, 1, 5, 20, 25, 0, 0, time.Local)

	err = service.ModifyEvent(ctx, updated)
	require.NoError(t, err)

	// Проверка, что событие было изменено в хранилище
	storedEvent, err := service.GetEvent(ctx, id)
	require.NoError(t, err)
	require.NotEqual(t, event, storedEvent)
	require.Equal(t, updated, storedEvent)

	// Проверка ошибки пересечений
	// добавим второе событие
	secondEvent := model.Event{
		StartDatetime: time.Date(2024, 1, 5, 19, 0, 0, 0, time.Local),
		EndDatetime:   time.Date(2024, 1, 5, 19, 15, 0, 0, time.Local),
	}
	id, err = service.CreateEvent(ctx, secondEvent)
	secondEvent.Id = id
	require.NoError(t, err)

	// Проверка отсутствия ошибки если правим само себя
	err = service.ModifyEvent(ctx, secondEvent)
	require.NoError(t, err)

	// Проверка ошибки событие не найдено
	err = service.ModifyEvent(ctx, model.Event{
		Id:            333,
		StartDatetime: time.Date(2024, 1, 1, 19, 0, 0, 0, time.Local),
		EndDatetime:   time.Date(2024, 1, 1, 19, 15, 0, 0, time.Local),
	})
	require.Error(t, err)
	require.ErrorIs(t, err, errors.ErrEventNotFound)

	// правим ему время, ожидаем ошибку
	secondEvent.EndDatetime = time.Date(2024, 1, 5, 23, 15, 0, 0, time.Local)
	err = service.ModifyEvent(ctx, secondEvent)
	require.Error(t, err)
	require.ErrorIs(t, err, errors.ErrDateTimeBusy)

	// Проверка ошибки некорректного периода
	secondEvent.EndDatetime = time.Date(2024, 1, 5, 13, 15, 0, 0, time.Local)
	err = service.ModifyEvent(ctx, secondEvent)
	require.Error(t, err)
	require.ErrorIs(t, err, errors.ErrIncorrectDates)
}

func TestRemoveEvent(t *testing.T) {
	service := NewEventsService(new(MockLogger), memorystorage.NewStorage())
	ctx := context.Background()

	// ошибка на событие не найдено
	err := service.RemoveEvent(ctx, 1)
	require.Error(t, err)
	require.ErrorIs(t, err, errors.ErrEventNotFound)

	// создаём
	event := model.Event{
		StartDatetime: time.Date(2024, 1, 5, 20, 0, 0, 0, time.Local),
		EndDatetime:   time.Date(2024, 1, 5, 20, 15, 0, 0, time.Local),
	}

	id, err := service.CreateEvent(ctx, event)
	event.Id = id
	require.NoError(t, err)

	// удаляем
	err = service.RemoveEvent(ctx, id)
	require.NoError(t, err)

	// Проверка, что событие было удалено в хранилище
	_, err = service.GetEvent(ctx, id)
	require.Error(t, err)
	require.ErrorIs(t, err, errors.ErrEventNotFound)
}

func TestGetEvent(t *testing.T) {
	service := NewEventsService(new(MockLogger), memorystorage.NewStorage())
	ctx := context.Background()

	// создаём
	event := model.Event{
		StartDatetime: time.Date(2024, 1, 5, 20, 0, 0, 0, time.Local),
		EndDatetime:   time.Date(2024, 1, 5, 20, 15, 0, 0, time.Local),
	}

	id, err := service.CreateEvent(ctx, event)
	event.Id = id
	require.NoError(t, err)

	// проверка
	stored, err := service.GetEvent(ctx, id)
	require.NoError(t, err)
	require.Equal(t, event, stored)

	// ошибка на событие не найдено
	_, err = service.GetEvent(ctx, 333)
	require.Error(t, err)
	require.ErrorIs(t, err, errors.ErrEventNotFound)
}

func TestGetEventFor(t *testing.T) {
	service := NewEventsService(new(MockLogger), memorystorage.NewStorage())
	ctx := context.Background()

	// создаём
	for _, e := range events() {
		_, err := service.CreateEvent(ctx, e)
		require.NoError(t, err)
	}

	// проверка дня
	stored, err := service.GetEventsForDay(ctx, time.Date(2024, 1, 31, 7, 0, 0, 0, time.Local))
	require.NoError(t, err)
	require.Equal(t, 1, len(stored))

	// неделя
	stored, err = service.GetEventsForWeek(ctx, time.Date(2024, 1, 29, 7, 0, 0, 0, time.Local))
	require.NoError(t, err)
	require.Equal(t, 7, len(stored))

	// месяц-1
	stored, err = service.GetEventsForMonth(ctx, time.Date(2024, 1, 1, 7, 0, 0, 0, time.Local))
	require.NoError(t, err)
	require.Equal(t, 3, len(stored))

	// месяц-2
	stored, err = service.GetEventsForMonth(ctx, time.Date(2024, 2, 1, 7, 0, 0, 0, time.Local))
	require.NoError(t, err)
	require.Equal(t, 7, len(stored))
}

func events() []model.Event {
	result := make([]model.Event, 0, 10)
	for i := 0; i < 10; i++ {
		result = append(result,
			model.Event{
				Id:            i + 1,
				StartDatetime: time.Date(2024, 1, i+29, 8, 0, 0, 0, time.Local),
				EndDatetime:   time.Date(2024, 1, i+29, 13, 0, 0, 0, time.Local),
			},
		)
	}
	return result
}
