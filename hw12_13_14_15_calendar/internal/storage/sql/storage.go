package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/errors"
	"github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/model"
)

type Storage struct { // TODO
	db  *sqlx.DB
	dsn string
}

func NewStorage(host string, port int16, user, pwd, db string) *Storage {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, pwd, db)
	return &Storage{
		dsn: dsn,
	}
}

func (s *Storage) Connect(ctx context.Context) (err error) {
	s.db, err = sqlx.Open("pgx", s.dsn)
	if err != nil {
		return err
	}
	return s.db.PingContext(ctx)
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) InsertEvent(ctx context.Context, event model.Event) (int, error) {
	q := `
		insert into events(title, start_date_time, end_date_time, description, user_id, notify_before_min)
		values (:title, :start_date_time, :end_date_time, :description, :user_id, :notify_before_min)
		returning id
	`
	rows, err := s.db.NamedQueryContext(ctx, q, event)
	if err != nil {
		return -1, err
	}
	defer rows.Close()
	if !rows.Next() {
		return -1, fmt.Errorf("received 0 rows")
	}
	var id int
	rows.Scan(&id)
	return id, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event model.Event) error {
	q := `
		update events
		set
			title = :title,
			start_date_time = :start_date_time,
			end_date_time = :end_date_time,
			description = :description,
			user_id = :user_id,
			notify_before_min = :notify_before_min
		where
			id = :id
	`
	result, err := s.db.NamedExecContext(ctx, q, event)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected < 1 {
		return errors.ErrEventNotFound
	}

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, eventId int) error {
	q := `
		delete
		from events
		where
			id = $1
	`
	result, err := s.db.ExecContext(ctx, q, eventId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected < 1 {
		return errors.ErrEventNotFound
	}

	return nil
}

func (s *Storage) GetEvent(ctx context.Context, eventId int) (model.Event, error) {
	q := `
		select
			id,
			title,
			start_date_time,
			end_date_time,
			description,
			user_id,
			notify_before_min
		from events
		where id = $1
	`
	ev := model.Event{}
	err := s.db.GetContext(ctx, &ev, q, eventId)
	if err == sql.ErrNoRows {
		return ev, errors.ErrEventNotFound
	}

	return ev, err
}

func (s *Storage) GetEvents(ctx context.Context, dateFrom, dateTo time.Time) ([]model.Event, error) {
	q := `
		select
			id,
			title,
			start_date_time,
			end_date_time,
			description,
			user_id,
			notify_before_min
		from events
		where start_date_time >= :date_from and start_date_time < :date_to
	`
	return namedSelect(s.db, ctx, q, map[string]interface{}{
		"date_from": dateFrom,
		"date_to":   dateTo,
	})
}

func (s *Storage) IsTimeFree(ctx context.Context, excludeId int, dateFrom, dateTo time.Time) (bool, error) {
	q := `
		select count(*) as cnt
		from events
		where
			id != :excludeId
			and (
				(start_date_time <= :date_from && end_date_time >= :date_to) ||
				-- start >= 13 && <= 14
				(start_date_time > :date_from && start_date_time < :date_to) ||
				-- end between > 13 && <= 14
				(end_date_time > :date_from && end_date_time < :date_to)
			)
	`
	rows, err := s.db.NamedQueryContext(ctx, q, map[string]interface{}{
		"date_from": dateFrom,
		"date_to":   dateTo,
	})
	if err != nil {
		return false, err
	}
	defer rows.Close()
	if !rows.Next() {
		return false, fmt.Errorf("received 0 rows")
	}
	var cnt int
	rows.Scan(&cnt)
	return cnt == 0, nil
}

func namedSelect(db *sqlx.DB, ctx context.Context, query string, args interface{}) ([]model.Event, error) {
	rows, err := db.NamedQueryContext(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]model.Event, 0)
	ev := model.Event{}
	for rows.Next() {
		err := rows.StructScan(&ev)
		if err != nil {
			return nil, err
		}
		result = append(result, ev)
	}
	return result, nil
}
