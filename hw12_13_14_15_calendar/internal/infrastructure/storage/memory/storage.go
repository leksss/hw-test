package memorystorage

import (
	"context"
	"sync"

	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/domain/entities"
	uuid "github.com/satori/go.uuid"
)

type Storage struct {
	mu     sync.RWMutex //nolint
	events []*entities.Event
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context) error {
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return nil
}

func (s *Storage) CreateEvent(ctx context.Context, event *entities.Event) (string, error) {
	UUID := uuid.NewV4()
	event.UUID = UUID.String()
	s.events = append(s.events, event)
	return UUID.String(), nil
}

func (s *Storage) UpdateEvent(ctx context.Context, uuid string, event *entities.Event) (int64, error) {
	rowsCnt := int64(0)
	for i, e := range s.events {
		if e.UUID == uuid {
			event.UUID = uuid
			s.events[i] = event
			rowsCnt++
			break
		}
	}
	return rowsCnt, nil
}

func (s *Storage) DeleteEvent(ctx context.Context, uuid string) error {
	for i, e := range s.events {
		if e.UUID == uuid {
			s.events = append(s.events[:i], s.events[i+1:]...)
			break
		}
	}
	return nil
}

func (s *Storage) GetEventList(ctx context.Context, filter entities.Filter) ([]*entities.Event, error) {
	return s.events, nil
}
