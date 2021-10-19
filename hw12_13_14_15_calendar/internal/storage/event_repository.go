package storage

import (
	sqlstorage "github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/storage/sql"
)

type Repository interface {
	Add(event Event) error
	Update(ID int, event Event) error
	Delete(ID int) error
	GetList(params map[string]string) []Event
}

type EventRepository struct {
	storage *sqlstorage.Storage
}

func New(storage *sqlstorage.Storage) *EventRepository {
	return &EventRepository{
		storage: storage,
	}
}

func (e *EventRepository) Add(event Event) error {
	return nil
}

func (e *EventRepository) Update(ID int, event Event) error {
	return nil
}

func (e *EventRepository) Delete(ID int) error {
	return nil
}

func (e *EventRepository) GetList(params map[string]string) []Event {
	return nil
}
