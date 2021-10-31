package sqlstorage

import (
	"context"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // nolint
	"github.com/jmoiron/sqlx"
	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/domain/entities"
	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/domain/interfaces"
	uuid "github.com/satori/go.uuid"
)

type Storage struct {
	db   *sqlx.DB
	conf interfaces.DatabaseConf
}

func New(conf interfaces.DatabaseConf) *Storage {
	return &Storage{
		conf: conf,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	dsn := fmt.Sprintf("%s:%s@(%s:3306)/%s", s.conf.User, s.conf.Password, s.conf.Host, s.conf.Name)
	db, err := sqlx.ConnectContext(ctx, "mysql", dsn)
	if err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	if err := s.db.Close(); err != nil {
		return err
	}
	return nil
}

func (s *Storage) CreateEvent(ctx context.Context, event entities.Event) (string, error) {
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

func (s *Storage) UpdateEvent(ctx context.Context, id string, event entities.Event) error {
	sql := `UPDATE event SET 
                 owner_id=:OwnerID, 
                 title=:Title, 
                 started_at=:StartedAt, 
                 ended_at=:EndedAt, 
                 text=:Text, 
                 notify_for=:NotifyFor 
			WHERE id=:Id`
	_, err := s.db.NamedExecContext(ctx, sql, map[string]interface{}{
		"OwnerID":   event.OwnerID,
		"Title":     event.Title,
		"StartedAt": event.StartedAt,
		"EndedAt":   event.EndedAt,
		"Text":      event.Text,
		"NotifyFor": event.NotifyFor,
		"Id":        id,
	})
	return err
}

func (s *Storage) DeleteEvent(ctx context.Context, id string) error {
	sql := `DELETE FROM event WHERE id=:id`
	_, err := s.db.NamedExecContext(ctx, sql, map[string]interface{}{
		"id": id,
	})
	return err
}

func (s *Storage) GetEventList(ctx context.Context, params map[string]string) ([]entities.Event, error) {
	sql := `SELECT * FROM event LIMIT 10`
	rows, err := s.db.NamedQueryContext(ctx, sql, map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []entities.Event
	event := eventDB{}
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
