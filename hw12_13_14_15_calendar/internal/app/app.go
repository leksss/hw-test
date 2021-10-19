package app

import (
	"context"

	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type App struct { // TODO
	logger  Logger
	storage Storage
}

type Logger interface { // TODO
}

type Storage interface { // TODO
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger: logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
