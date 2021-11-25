package sqlstorage

import (
	"context"

	_ "github.com/go-sql-driver/mysql" // nolint
	"github.com/jmoiron/sqlx"
	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/domain/entities"
	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/domain/errors"
	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/domain/interfaces"
	uuid "github.com/satori/go.uuid"
)

type Storage struct {
	db   *sqlx.DB
	conf interfaces.DatabaseConf
}

func New(conf interfaces.DatabaseConf, db *sqlx.DB) *Storage {
	return &Storage{
		conf: conf,
		db:   db,
	}
}

func (s *Storage) CreateEvent(ctx context.Context, event entities.Event) (string, error) {
	if err := s.isTimeForEventAvailable(ctx, "", event); err != nil {
		return "", err
	}

	uuID := uuid.NewV4()
	sql := `INSERT INTO event (id, owner_id, title, started_at, ended_at, text, notify_for) 
			VALUES (:ID, :OwnerID, :Title, :StartedAt, :EndedAt, :Text, :NotifyFor)`
	_, err := s.db.NamedExecContext(ctx, sql, map[string]interface{}{
		"ID":        uuID.String(),
		"OwnerID":   event.OwnerID,
		"Title":     event.Title,
		"StartedAt": event.StartedAt,
		"EndedAt":   event.EndedAt,
		"Text":      event.Text,
		"NotifyFor": event.NotifyFor,
	})
	if err != nil {
		return "", err
	}
	return uuID.String(), nil
}

func (s *Storage) UpdateEvent(ctx context.Context, eventID string, event entities.Event) error {
	if err := s.isTimeForEventAvailable(ctx, eventID, event); err != nil {
		return err
	}

	sql := `UPDATE event SET 
                 owner_id=:OwnerID, 
                 title=:Title, 
                 started_at=:StartedAt, 
                 ended_at=:EndedAt, 
                 text=:Text, 
                 notify_for=:NotifyFor 
			WHERE id=:EventID`
	result, err := s.db.NamedExecContext(ctx, sql, map[string]interface{}{
		"OwnerID":   event.OwnerID,
		"Title":     event.Title,
		"StartedAt": event.StartedAt,
		"EndedAt":   event.EndedAt,
		"Text":      event.Text,
		"NotifyFor": event.NotifyFor,
		"EventID":   eventID,
	})
	if err != nil {
		return err
	}
	cnt, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if cnt == 0 {
		return errors.ErrNoAffectedEvent
	}
	return nil
}

func (s *Storage) isTimeForEventAvailable(ctx context.Context, eventID string, event entities.Event) error {
	sql := `SELECT id FROM event
			WHERE id!=:EventID AND owner_id=:OwnerID AND 
				(:StartedAt BETWEEN started_at AND ended_at OR started_at BETWEEN :StartedAt AND :EndedAt)
			LIMIT 1`
	rows, err := s.db.NamedQueryContext(ctx, sql, map[string]interface{}{
		"OwnerID":   event.OwnerID,
		"StartedAt": event.StartedAt,
		"EndedAt":   event.EndedAt,
		"EventID":   eventID,
	})
	if err != nil {
		return err
	}
	defer rows.Close()

	var rowEventID int64
	if rows.Next() {
		err = rows.Scan(&rowEventID)
		if err != nil {
			return err
		}
	}
	if rowEventID > 0 {
		return errors.ErrDateBusy
	}
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, eventID string) error {
	event, err := s.getEventByID(ctx, eventID)
	if err != nil {
		return err
	}
	if event == nil {
		return errors.ErrNoEventFound
	}
	sql := `DELETE FROM event WHERE id=:eventID`
	_, err = s.db.NamedExecContext(ctx, sql, map[string]interface{}{
		"eventID": eventID,
	})
	return err
}

func (s *Storage) GetEventList(ctx context.Context, limit, offset int64) ([]entities.Event, error) {
	sql := `SELECT * FROM event LIMIT :limit OFFSET :offset`
	rows, err := s.db.NamedQueryContext(ctx, sql, map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []entities.Event
	var event eventDB
	for rows.Next() {
		err = rows.StructScan(&event)
		if err != nil {
			return nil, err
		}
		events = append(events, entities.Event{
			ID:        event.ID,
			OwnerID:   event.OwnerID,
			Title:     event.Title,
			StartedAt: &event.StartedAt.Time,
			EndedAt:   &event.EndedAt.Time,
			Text:      event.Text,
			NotifyFor: event.NotifyFor,
		})
	}
	return events, nil
}

func (s *Storage) getEventByID(ctx context.Context, eventID string) (*entities.Event, error) {
	sql := `SELECT * FROM event	WHERE id=:EventID`
	rows, err := s.db.NamedQueryContext(ctx, sql, map[string]interface{}{
		"EventID": eventID,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var event *entities.Event
	if rows.Next() {
		err = rows.StructScan(&event)
		if err != nil {
			return nil, err
		}
	}
	return event, nil
}
