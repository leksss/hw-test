package memorystorage

import (
	"context"
	"sync"

	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/domain/entities"
	uuid "github.com/satori/go.uuid"
)

type Storage struct {
	mu     sync.RWMutex //nolint
	events []entities.Event
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

func (s *Storage) CreateEvent(ctx context.Context, event entities.Event) (string, error) {
	uuID := uuid.NewV4()
	event.ID = uuID.String()
	s.events = append(s.events, event)
	return uuID.String(), nil
}

func (s *Storage) UpdateEvent(ctx context.Context, id string, event entities.Event) error {
	for i, e := range s.events {
		if e.ID == id {
			event.ID = id
			s.events[i] = event
			break
		}
	}
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id string) error {
	for i, e := range s.events {
		if e.ID == id {
			s.events = append(s.events[:i], s.events[i+1:]...)
			break
		}
	}
	return nil
}

func (s *Storage) GetEventList(ctx context.Context, params map[string]string) ([]entities.Event, error) {
	return s.events, nil
}
