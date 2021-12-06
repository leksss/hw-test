package memorystorage

import (
	"context"
	"testing"

	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/domain/entities"
	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/domain/interfaces"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	events := []entities.Event{
		{
			OwnerID: "327bae0a-8383-4323-93c4-f0501b8380cd",
			Title:   "Тестовое событие 1",
		},
		{
			OwnerID: "327bae0a-8383-4323-93c4-f0501b8380cd",
			Title:   "Тестовое событие 2",
		},
	}

	storage := New()

	eventIDs := make([]string, 0)
	for _, e := range events {
		eventID, err := storage.CreateEvent(context.Background(), e)
		require.NoError(t, err)
		require.Equal(t, 36, len(eventID))
		eventIDs = append(eventIDs, eventID)
	}

	for _, eventID := range eventIDs {
		err := storage.UpdateEvent(context.Background(), eventID, entities.Event{
			Title: "Тестовое событие UPDATED",
		})
		require.NoError(t, err)
	}

	events, err := storage.GetEventList(context.Background(), interfaces.EventListLimit, 0)
	require.Equal(t, 2, len(events))
	for _, event := range events {
		require.NoError(t, err)
		require.Equal(t, "Тестовое событие UPDATED", event.Title)
	}

	for _, eventID := range eventIDs {
		err = storage.DeleteEvent(context.Background(), eventID)
		require.NoError(t, err)
	}

	events, err = storage.GetEventList(context.Background(), interfaces.EventListLimit, 0)
	require.NoError(t, err)
	require.Equal(t, 0, len(events))
}
