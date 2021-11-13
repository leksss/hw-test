package sqlstorage

import (
	"context"
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // nolint
	"github.com/jmoiron/sqlx"
	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/domain/entities"
	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/domain/interfaces"
	"github.com/leksss/hw-test/hw12_13_14_15_calendar/internal/infrastructure/logger"
	uuid "github.com/satori/go.uuid"
)

const DefaultLimit = 20

type Storage struct {
	db   *sqlx.DB
	conf interfaces.DatabaseConf
	log  logger.Log
}

func New(conf interfaces.DatabaseConf, log logger.Log) *Storage {
	return &Storage{
		conf: conf,
		log:  log,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	dsn := fmt.Sprintf("%s:%s@(%s:3306)/%s?parseTime=true", s.conf.User, s.conf.Password, s.conf.Host, s.conf.Name)
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

func (s *Storage) CreateEvent(ctx context.Context, event *entities.Event) (string, error) {
	uuID := uuid.NewV4()
	sql := `INSERT INTO event (uuid, owner_id, title, started_at, ended_at, text, notify_for) 
			VALUES (:UUID, :OwnerID, :Title, :StartedAt, :EndedAt, :Text, :NotifyFor)`
	arg := map[string]interface{}{
		"UUID":      uuID.String(),
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

func (s *Storage) UpdateEvent(ctx context.Context, uuid string, event *entities.Event) (int64, error) {
	sql := `UPDATE event SET 
                 owner_id=:ownerID, 
                 title=:title, 
                 started_at=:startedAt, 
                 ended_at=:endedAt, 
                 text=:text, 
                 notify_for=:notifyFor 
			WHERE uuid=:uuid`
	arg := map[string]interface{}{
		"ownerID":   event.OwnerID,
		"title":     event.Title,
		"startedAt": event.StartedAt,
		"endedAt":   event.EndedAt,
		"text":      event.Text,
		"notifyFor": event.NotifyFor,
		"uuid":      uuid,
	}
	res, err := s.db.NamedExecContext(ctx, sql, arg)
	s.logQuery(sql, arg)
	if err != nil {
		return 0, err
	}
	rowsCnt, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsCnt, err
}

func (s *Storage) DeleteEvent(ctx context.Context, uuid string) error {
	sql := `DELETE FROM event WHERE uuid=:uuid`
	arg := map[string]interface{}{
		"uuid": uuid,
	}
	_, err := s.db.NamedExecContext(ctx, sql, arg)
	s.logQuery(sql, arg)
	return err
}

func (s *Storage) GetEventList(ctx context.Context, filter entities.Filter) ([]*entities.Event, error) {
	if filter.Limit == 0 {
		filter.Limit = DefaultLimit
	}

	var sql string
	var arg map[string]interface{}
	if filter.UUID == "" {
		sql = `SELECT * FROM event LIMIT :limit OFFSET :offset`
		arg = map[string]interface{}{
			"limit":  filter.Limit,
			"offset": filter.Offset,
		}
	} else {
		sql = `SELECT * FROM event WHERE uuid = :uuid`
		arg = map[string]interface{}{
			"uuid": filter.UUID,
		}
	}

	rows, err := s.db.NamedQueryContext(ctx, sql, arg)
	s.logQuery(sql, arg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*entities.Event
	event := eventDB{}
	for rows.Next() {
		err = rows.StructScan(&event)
		if err != nil {
			return nil, err
		}

		events = append(events, &entities.Event{
			UUID:      event.ID,
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

func (s *Storage) logQuery(sql string, arg map[string]interface{}) {
	byteArg, _ := json.MarshalIndent(arg, "", "  ")
	s.log.Info(fmt.Sprintf("%s %s", sql, string(byteArg)))
}
