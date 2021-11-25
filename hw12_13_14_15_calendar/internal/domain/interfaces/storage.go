package interfaces

import (
	"context"

	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/domain/entities"
)

const EventListLimit = 20

type DatabaseConf struct {
	Host     string
	User     string
	Password string
	Name     string
}

type Storage interface {
	CreateEvent(ctx context.Context, event entities.Event) (string, error)
	UpdateEvent(ctx context.Context, eventID string, event entities.Event) error
	DeleteEvent(ctx context.Context, eventID string) error
	GetEventList(ctx context.Context, limit, offset int64) ([]entities.Event, error)
}
