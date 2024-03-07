package errors

import "errors"

var (
	ErrEventNotFound  = errors.New("event not found")
	ErrDateTimeBusy   = errors.New("date & time busy")
	ErrIncorrectDates = errors.New("incorrect dates")
)
