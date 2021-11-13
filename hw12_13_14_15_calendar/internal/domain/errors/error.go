package errors

type EventError string

var (
	ErrEventOwnerIDIsRequired   = EventError("ownerID is required")
	ErrEventTitleIsRequired     = EventError("title is required")
	ErrEventStartedAtIsRequired = EventError("startedAt is required")
	ErrEventEndedAtIsRequired   = EventError("endedAt is required")
	ErrEventUUIDIsRequired      = EventError("UUID is required")
	ErrEventNotFound            = EventError("Event not found")
)

func (ee EventError) Error() string {
	return string(ee)
}
