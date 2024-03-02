package model

import (
	"time"
)

type Event struct {
	Id              int       `db:"id"`
	Title           string    `db:"title"`
	StartDatetime   time.Time `db:"start_date_time"`
	EndDatetime     time.Time `db:"end_date_time"`
	Description     string    `db:"description"`
	UserId          int       `db:"user_id"`
	NotifyBeforeMin int       `db:"notify_before_min"`
}

func (e Event) ValidPeriod() bool {
	return e.StartDatetime.Before(e.EndDatetime)
}
