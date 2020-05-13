package session

import (
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"context"
	"time"
)

type UserDevice struct {
	tableName struct{}    `pg:"maasapi.user_device_session"`
	ID        int         `pg:"id,pk" json:"id"`
	DeviceID  string      `pg:"device_id,notnull" json:"device_id"`
	UserID    int         `pg:"user_id,notnull" json:"user_id"`
	CreatedAt models.Time `json:"-"`
	DeletedAt time.Time   `pg:",soft_delete"`
}

func (s *UserDevice) BeforeInsert(ctx context.Context) (context.Context, error) {
	s.CreatedAt = models.Time{Time: time.Now()}

	return ctx, nil
}
