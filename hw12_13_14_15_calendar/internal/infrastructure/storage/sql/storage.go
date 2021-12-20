package sqlstorage

import (
	"context"
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // nolint
	"github.com/jmoiron/sqlx"
	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/domain/entities"
	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/domain/errors"
	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/domain/interfaces"
	uuid "github.com/satori/go.uuid"
)

const DefaultListLimit = 20

type Storage struct {
	db  *sqlx.DB
	log interfaces.Log
}

func New(db *sqlx.DB, log interfaces.Log) *Storage {
	return &Storage{
		db:  db,
		log: log,
	}
}

func (s *Storage) CreateEvent(ctx context.Context, event entities.Event) (string, error) {
	if err := s.isTimeForEventAvailable(ctx, "", event); err != nil {
		return "", err
	}

	uuID := uuid.NewV4()
	sql := `INSERT INTO event (id, owner_id, title, started_at, ended_at, text, notify_for) 
			VALUES (:EventID, :OwnerID, :Title, :StartedAt, :EndedAt, :Text, :NotifyFor)`
	arg := map[string]interface{}{
		"EventID":   uuID.String(),
		"OwnerID":   event.OwnerID,
		"Title":     event.Title,
		"StartedAt": event.StartedAt,
		"EndedAt":   event.EndedAt,
		"Text":      event.Text,
		"NotifyFor": event.NotifyFor,
	}
	_, err := s.db.NamedExecContext(ctx, sql, arg)
	s.logQuery(sql, arg)
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
	arg := map[string]interface{}{
		"OwnerID":   event.OwnerID,
		"Title":     event.Title,
		"StartedAt": event.StartedAt,
		"EndedAt":   event.EndedAt,
		"Text":      event.Text,
		"NotifyFor": event.NotifyFor,
		"EventID":   eventID,
	}
	_, err := s.db.NamedExecContext(ctx, sql, arg)
	s.logQuery(sql, arg)
	if err != nil {
		return err
	}
	return err
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
	arg := map[string]interface{}{
		"eventID": eventID,
	}
	_, err = s.db.NamedExecContext(ctx, sql, arg)
	s.logQuery(sql, arg)
	return err
}

func (s *Storage) GetEventList(ctx context.Context, filter entities.EventListFilter) ([]*entities.Event, error) {
	if filter.Limit == 0 {
		filter.Limit = DefaultListLimit
	}

	var sql string
	var arg map[string]interface{}
	if filter.EventID == "" {
		sql = `SELECT * FROM event LIMIT :limit OFFSET :offset`
		arg = map[string]interface{}{
			"limit":  filter.Limit,
			"offset": filter.Offset,
		}
	} else {
		sql = `SELECT * FROM event WHERE id = :EventID`
		arg = map[string]interface{}{
			"EventID": filter.EventID,
		}
	}

	rows, err := s.db.NamedQueryContext(ctx, sql, arg)
	s.logQuery(sql, arg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*entities.Event
	var event eventDB
	for rows.Next() {
		err = rows.StructScan(&event)
		if err != nil {
			return nil, err
		}

		startedAt := event.StartedAt.Time
		endedAt := event.EndedAt.Time

		events = append(events, &entities.Event{
			EventID:   event.ID,
			OwnerID:   event.OwnerID,
			Title:     event.Title,
			StartedAt: &startedAt,
			EndedAt:   &endedAt,
			Text:      event.Text,
			NotifyFor: event.NotifyFor,
		})
	}
	return events, nil
}

func (s *Storage) getEventByID(ctx context.Context, eventID string) (*entities.Event, error) {
	sql := `SELECT * FROM event	WHERE id=:EventID`
	arg := map[string]interface{}{
		"EventID": eventID,
	}
	rows, err := s.db.NamedQueryContext(ctx, sql, arg)
	s.logQuery(sql, arg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, nil
	}

	var eventDB eventDB
	err = rows.StructScan(&eventDB)
	if err != nil {
		return nil, err
	}

	return &entities.Event{
		EventID:   eventDB.ID,
		OwnerID:   eventDB.OwnerID,
		Title:     eventDB.Title,
		StartedAt: &eventDB.StartedAt.Time,
		EndedAt:   &eventDB.EndedAt.Time,
		Text:      eventDB.Text,
		NotifyFor: eventDB.NotifyFor,
	}, nil
}

func (s *Storage) isTimeForEventAvailable(ctx context.Context, eventID string, event entities.Event) error {
	sql := `SELECT id FROM event
			WHERE id!=:EventID AND owner_id=:OwnerID AND 
				(:StartedAt BETWEEN started_at AND ended_at OR started_at BETWEEN :StartedAt AND :EndedAt)
			LIMIT 1`
	arg := map[string]interface{}{
		"OwnerID":   event.OwnerID,
		"StartedAt": event.StartedAt,
		"EndedAt":   event.EndedAt,
		"EventID":   eventID,
	}
	rows, err := s.db.NamedQueryContext(ctx, sql, arg)
	s.logQuery(sql, arg)
	if err != nil {
		return err
	}
	defer rows.Close()

	var rowEventID string
	if rows.Next() {
		err = rows.Scan(&rowEventID)
		if err != nil {
			return err
		}
	}
	if rowEventID != "" {
		return errors.ErrDateBusy
	}
	return nil
}

func (s *Storage) logQuery(sql string, arg map[string]interface{}) {
	byteArg, _ := json.MarshalIndent(arg, "", "  ")
	s.log.Info(fmt.Sprintf("%s %s", sql, string(byteArg)))
}
