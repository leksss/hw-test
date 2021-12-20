package interfaces

import (
	"context"

	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/domain/entities"
)

type Storage interface {
	CreateEvent(ctx context.Context, event entities.Event) (string, error)
	UpdateEvent(ctx context.Context, eventID string, event entities.Event) error
	DeleteEvent(ctx context.Context, eventID string) error
	GetEventList(ctx context.Context, filter entities.EventListFilter) ([]*entities.Event, error)
}
