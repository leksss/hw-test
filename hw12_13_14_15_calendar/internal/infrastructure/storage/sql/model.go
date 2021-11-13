package sqlstorage

import "database/sql"

type eventDB struct {
	ID        string       `db:"uuid"`
	OwnerID   string       `db:"owner_id"`
	Title     string       `db:"title"`
	StartedAt sql.NullTime `db:"started_at"`
	EndedAt   sql.NullTime `db:"ended_at"`
	Text      string       `db:"text"`
	NotifyFor uint64       `db:"notify_for"`
}
