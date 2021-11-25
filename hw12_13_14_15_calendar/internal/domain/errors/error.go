package errors

import "errors"

var (
	ErrDateBusy        = errors.New("another event exists for this date")
	ErrNoAffectedEvent = errors.New("no affected event")
	ErrNoEventFound    = errors.New("no event found")
)
