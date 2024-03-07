package memorystorage

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/model"
	"github.com/stretchr/testify/require"
)

func TestInsertEvent(t *testing.T) {
	inmem := NewStorage()
	require.Equal(t, 0, len(inmem.events))

	event := model.Event{
		StartDatetime: time.Date(2024, 1, 5, 20, 0, 0, 0, time.Local),
		EndDatetime:   time.Date(2024, 1, 5, 20, 15, 0, 0, time.Local),
	}
	id, err := inmem.InsertEvent(context.Background(), event)
	event.Id = id

	require.NoError(t, err)
	require.Equal(t, 1, id)

	require.Contains(t, inmem.events, event)
}

func TestUpdateEvent(t *testing.T) {
	inmem := NewStorage()
	ctx := context.Background()

	event := model.Event{
		StartDatetime: time.Date(2024, 1, 5, 20, 0, 0, 0, time.Local),
		EndDatetime:   time.Date(2024, 1, 5, 20, 15, 0, 0, time.Local),
	}
	id, err := inmem.InsertEvent(ctx, event)
	event.Id = id
	require.NoError(t, err)

	updated := event
	updated.StartDatetime = time.Date(2024, 1, 5, 20, 30, 0, 0, time.Local)
	updated.EndDatetime = time.Date(2024, 1, 5, 20, 45, 0, 0, time.Local)

	err = inmem.UpdateEvent(ctx, updated)
	require.NoError(t, err)

	require.Contains(t, inmem.events, updated)
	require.NotContains(t, inmem.events, event)
}

func TestDeleteEvent(t *testing.T) {
	inmem := NewStorage()
	ctx := context.Background()

	event := model.Event{
		StartDatetime: time.Date(2024, 1, 5, 20, 0, 0, 0, time.Local),
		EndDatetime:   time.Date(2024, 1, 5, 20, 15, 0, 0, time.Local),
	}
	id, err := inmem.InsertEvent(ctx, event)
	event.Id = id
	require.NoError(t, err)

	err = inmem.DeleteEvent(ctx, 1)
	require.NoError(t, err)

	require.Equal(t, 0, len(inmem.events))
}

func TestGet(t *testing.T) {
	inmem := NewStorage()
	ctx := context.Background()

	expected := generateEvs(2024, 1, 31)
	for i := range expected {
		id, err := inmem.InsertEvent(ctx, expected[i])
		require.NoError(t, err)
		expected[i].Id = id
	}
	require.Equal(t, 31*2, len(expected))

	// once
	ev, err := inmem.GetEvent(ctx, 10)
	require.NoError(t, err)
	require.Equal(t, expected[10-1], ev)

	// 2 days
	dateFrom := time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local)
	dateTo := dateFrom.AddDate(0, 0, 2)
	evs, err := inmem.GetEvents(ctx, dateFrom, dateTo)
	require.NoError(t, err)

	expected = generateEvs(2024, 1, 2)
	require.Equal(t, expected, evs)
}

func TestIsTimeFree(t *testing.T) {
	inmem := NewStorage()
	ctx := context.Background()

	setup := []model.Event{
		{
			StartDatetime: time.Date(2024, 1, 1, 10, 0, 0, 0, time.Local),
			EndDatetime:   time.Date(2024, 1, 1, 11, 30, 0, 0, time.Local),
		},
		{
			StartDatetime: time.Date(2024, 1, 1, 12, 0, 0, 0, time.Local),
			EndDatetime:   time.Date(2024, 1, 1, 13, 0, 0, 0, time.Local),
		},
	}
	for i := range setup {
		_, err := inmem.InsertEvent(ctx, setup[i])
		require.NoError(t, err)
	}

	tests := []struct {
		in   model.Event
		free bool
	}{
		{
			in: model.Event{
				StartDatetime: time.Date(2024, 1, 1, 8, 0, 0, 0, time.Local),
				EndDatetime:   time.Date(2024, 1, 1, 10, 0, 0, 0, time.Local),
			},
			free: true,
		},
		{
			in: model.Event{
				StartDatetime: time.Date(2024, 1, 1, 11, 30, 0, 0, time.Local),
				EndDatetime:   time.Date(2024, 1, 1, 12, 0, 0, 0, time.Local),
			},
			free: true,
		},
		{
			in: model.Event{
				StartDatetime: time.Date(2024, 1, 1, 13, 0, 0, 0, time.Local),
				EndDatetime:   time.Date(2024, 1, 1, 14, 0, 0, 0, time.Local),
			},
			free: true,
		},
		//
		{
			in: model.Event{
				StartDatetime: time.Date(2024, 1, 1, 10, 0, 0, 0, time.Local),
				EndDatetime:   time.Date(2024, 1, 1, 11, 30, 0, 0, time.Local),
			},
			free: false,
		},
		{
			in: model.Event{
				StartDatetime: time.Date(2024, 1, 1, 12, 0, 0, 0, time.Local),
				EndDatetime:   time.Date(2024, 1, 1, 13, 0, 0, 0, time.Local),
			},
			free: false,
		},
		{
			in: model.Event{
				StartDatetime: time.Date(2024, 1, 1, 11, 0, 0, 0, time.Local),
				EndDatetime:   time.Date(2024, 1, 1, 12, 0, 0, 0, time.Local),
			},
			free: false,
		},
		{
			in: model.Event{
				StartDatetime: time.Date(2024, 1, 1, 12, 30, 0, 0, time.Local),
				EndDatetime:   time.Date(2024, 1, 1, 14, 0, 0, 0, time.Local),
			},
			free: false,
		},
		{
			in: model.Event{
				StartDatetime: time.Date(2024, 1, 1, 8, 30, 0, 0, time.Local),
				EndDatetime:   time.Date(2024, 1, 1, 10, 1, 0, 0, time.Local),
			},
			free: false,
		},
		{
			in: model.Event{
				StartDatetime: time.Date(2024, 1, 1, 12, 59, 59, 0, time.Local),
				EndDatetime:   time.Date(2024, 1, 1, 14, 0, 0, 0, time.Local),
			},
			free: false,
		},
	}

	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			t.Parallel()

			free, err := inmem.IsTimeFree(ctx, 0, tt.in.StartDatetime, tt.in.EndDatetime)
			require.NoError(t, err)
			require.Equal(t, tt.free, free)
		})
	}
}

func TestDeleteEventsOlderThan(t *testing.T) {
	inmem := NewStorage()
	ctx := context.Background()

	_, err := inmem.InsertEvent(ctx, model.Event{
		StartDatetime: time.Date(2024, 1, 5, 20, 0, 0, 0, time.Local),
		EndDatetime:   time.Date(2024, 1, 5, 20, 15, 0, 0, time.Local),
	})
	require.NoError(t, err)

	_, err = inmem.InsertEvent(ctx, model.Event{
		StartDatetime: time.Date(2024, 1, 5, 20, 15, 0, 0, time.Local),
		EndDatetime:   time.Date(2024, 1, 5, 20, 30, 0, 0, time.Local),
	})
	require.NoError(t, err)

	_, err = inmem.InsertEvent(ctx, model.Event{
		StartDatetime: time.Date(2024, 1, 6, 20, 0, 0, 0, time.Local),
		EndDatetime:   time.Date(2024, 1, 6, 20, 15, 0, 0, time.Local),
	})
	require.NoError(t, err)

	rows, err := inmem.DeleteEventsOlderThan(ctx, time.Date(2024, 1, 5, 20, 15, 0, 1, time.Local))
	require.NoError(t, err)

	require.Equal(t, int64(2), rows)
	require.Equal(t, 1, len(inmem.events))
	require.Equal(t, inmem.events[0].StartDatetime, time.Date(2024, 1, 6, 20, 0, 0, 0, time.Local))
}

func TestSetIsNotified(t *testing.T) {
	inmem := NewStorage()
	ctx := context.Background()

	event := model.Event{
		StartDatetime: time.Date(2024, 1, 5, 20, 0, 0, 0, time.Local),
		EndDatetime:   time.Date(2024, 1, 5, 20, 15, 0, 0, time.Local),
		IsNotified:    false,
	}
	id, err := inmem.InsertEvent(ctx, event)
	require.NoError(t, err)

	err = inmem.SetIsNotified(ctx, id)
	require.NoError(t, err)

	require.Equal(t, true, inmem.events[0].IsNotified)
}

func TestParallel(t *testing.T) {
	inmem := NewStorage()
	ctx := context.Background()

	test := func(y int) {
		evs := generateEvs(y, 1, 30)

		for i := range evs {
			id, err := inmem.InsertEvent(ctx, evs[i])
			require.NoError(t, err)
			evs[i].Id = id
			err = inmem.UpdateEvent(ctx, evs[i])
			require.NoError(t, err)
		}

		dateFrom := evs[0].StartDatetime
		dateTo := dateFrom.AddDate(0, 1, 0)
		_, err := inmem.GetEvents(ctx, dateFrom, dateTo)
		require.NoError(t, err)

		for _, e := range evs {
			err := inmem.DeleteEvent(ctx, e.Id)
			require.NoError(t, err)
		}
	}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		y := i
		go func() {
			defer wg.Done()
			test(2000 + y)
		}()
	}
	wg.Wait()

	require.Equal(t, 0, len(inmem.events))
}

func generateEvs(y, n1, n2 int) []model.Event {
	result := make([]model.Event, 0, (n2-n1)*2)
	id := 0
	for i := n1 - 1; i < n2; i++ {
		result = append(result,
			model.Event{
				Id:            id + 1,
				Title:         strconv.Itoa(i*2 + 1),
				StartDatetime: time.Date(y, time.Month(1), i+1, 8, 0, 0, 0, time.Local),
				EndDatetime:   time.Date(y, time.Month(1), i+1, 13, 0, 0, 0, time.Local),
			},
			model.Event{
				Id:            id + 2,
				Title:         strconv.Itoa(i*2 + 2),
				StartDatetime: time.Date(y, time.Month(1), i+1, 14, 0, 0, 0, time.Local),
				EndDatetime:   time.Date(y, time.Month(1), i+1, 20, 0, 0, 0, time.Local),
			})
		id += 2
	}
	return result
}
