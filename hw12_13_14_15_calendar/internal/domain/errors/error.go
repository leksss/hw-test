package errors

var ErrDateBusy = EventError("another event exists for this date")

type EventError string

func (ee EventError) Error() string {
	return string(ee)
}
