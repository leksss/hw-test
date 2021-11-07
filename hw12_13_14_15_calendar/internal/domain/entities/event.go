package entities

import "time"

type Event struct {
	ID        string     // Уникальный идентификатор события UUID
	OwnerID   string     // ID пользователя, владельца события
	Title     string     // Заголовок - короткий текст
	StartedAt *time.Time // Дата и время начала события
	EndedAt   *time.Time // Дата и время окончания события
	Text      string     // Описание события - длинный текст, опционально
	NotifyFor int64      // За сколько времени высылать уведомление, опционально
}
