package memorystorage

import (
	"context"
	"sync"

	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/domain/entities"
	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/domain/errors"
	uuid "github.com/satori/go.uuid"
)

type Storage struct {
	mu        sync.RWMutex
	eventsMap map[string]entities.Event
}

func New() *Storage {
	return &Storage{
		eventsMap: make(map[string]entities.Event),
	}
}

func (s *Storage) CreateEvent(ctx context.Context, event entities.Event) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	uuID := uuid.NewV4()
	event.EventID = uuID.String()
	s.eventsMap[uuID.String()] = event

	return uuID.String(), nil
}

func (s *Storage) UpdateEvent(ctx context.Context, eventID string, event entities.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.eventsMap[eventID]; ok {
		event.EventID = eventID
		s.eventsMap[eventID] = event
	} else {
		return errors.ErrNoEventFound
	}
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, eventID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.eventsMap[eventID]; ok {
		delete(s.eventsMap, eventID)
	} else {
		return errors.ErrNoEventFound
	}
	return nil
}

func (s *Storage) GetEventList(ctx context.Context, filter entities.EventListFilter) ([]*entities.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := make([]*entities.Event, 0)
	for _, event := range s.eventsMap {
		events = append(events, &event)
	}
	return events, nil
}
