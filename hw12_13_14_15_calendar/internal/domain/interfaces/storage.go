package interfaces

import (
	"context"

	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/domain/entities"
)

type DatabaseConf struct {
	Host     string
	User     string
	Password string
	Name     string
}

type Storage interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	CreateEvent(ctx context.Context, event *entities.Event) (string, error)
	UpdateEvent(ctx context.Context, uuid string, event *entities.Event) (int64, error)
	DeleteEvent(ctx context.Context, uuid string) error
	GetEventList(ctx context.Context, filter entities.Filter) ([]*entities.Event, error)
}
