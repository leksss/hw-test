package errors

import "errors"

var (
	ErrDateBusy                 = errors.New("another event exists for this date")
	ErrNoEventFound             = errors.New("no event found")
	ErrEventOwnerIDIsRequired   = errors.New("ownerID is required")
	ErrEventTitleIsRequired     = errors.New("title is required")
	ErrEventStartedAtIsRequired = errors.New("startedAt is required")
	ErrEventEndedAtIsRequired   = errors.New("endedAt is required")
	ErrEventIDIsRequired        = errors.New("EventID is required")

	//ErrEventNotFound   = errors.New("event not found")
	//ErrNoAffectedEvent = errors.New("no affected event")
)
