package app

import (
	"context"

	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/domain/entities"
	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/domain/interfaces"
	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/infrastructure/logger"
)

type App struct {
	logger  logger.Log
	storage interfaces.Storage
}

func New(logger logger.Log, storage interfaces.Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, title string) (string, error) {
	event := entities.Event{
		Title: title,
	}
	return a.storage.CreateEvent(ctx, event)
}

func (a *App) UpdateEvent(ctx context.Context, id, title string) error {
	event := entities.Event{
		Title: title,
	}
	return a.storage.UpdateEvent(ctx, id, event)
}

func (a *App) DeleteEvent(ctx context.Context, id string) error {
	return a.storage.DeleteEvent(ctx, id)
}

func (a *App) GetEventList(ctx context.Context) ([]entities.Event, error) {
	return a.storage.GetEventList(ctx, map[string]string{})
}
