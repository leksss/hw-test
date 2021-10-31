package memorystorage

import (
	"context"
	"testing"

	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/domain/entities"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	storage := New()
	eventID, err := storage.CreateEvent(context.Background(), entities.Event{
		Title: "Тестовое событие",
	})
	require.NoError(t, err)
	require.Equal(t, 36, len(eventID))

	err = storage.UpdateEvent(context.Background(), eventID, entities.Event{
		Title: "Тестовое событие UPDATED",
	})
	require.NoError(t, err)

	events, err := storage.GetEventList(context.Background(), map[string]string{})
	require.NoError(t, err)
	require.Equal(t, 1, len(events))
	require.Equal(t, "Тестовое событие UPDATED", events[0].Title)

	err = storage.DeleteEvent(context.Background(), eventID)
	require.NoError(t, err)

	events, err = storage.GetEventList(context.Background(), map[string]string{})
	require.NoError(t, err)
	require.Equal(t, 0, len(events))
}
