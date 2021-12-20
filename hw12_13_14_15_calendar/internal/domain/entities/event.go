package entities

import "time"

type Event struct {
	EventID   string     // Уникальный идентификатор события EventID
	OwnerID   string     // ID пользователя, владельца события
	Title     string     // Заголовок - короткий текст
	StartedAt *time.Time // Дата и время начала события
	EndedAt   *time.Time // Дата и время окончания события
	Text      string     // Описание события - длинный текст, опционально
	NotifyFor uint64     // За сколько времени высылать уведомление, опционально
}

type EventListFilter struct {
	Limit   int64
	Offset  int64
	EventID string
}
