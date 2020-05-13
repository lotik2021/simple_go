package favorite

import (
	"bitbucket.movista.ru/maas/maasapi/internal/google"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"context"
	"time"
)

type WithGooglePlaceInfo struct {
	ID          int          `json:"id"`
	DeviceID    string       `json:"device_id"`
	Type        string       `json:"favorite_place_type"`
	Name        string       `json:"name"`
	InActions   bool         `json:"in_actions"`
	GooglePlace google.Place `json:"google_place,omitempty"`
	UpdatedAt   models.Time  `json:"updated_at"`
}

type DeviceFavorite struct {
	tableName struct{}    `pg:"maasapi.device_favorite"`
	ID        int         `pg:"id,pk" json:"id"`
	DeviceID  string      `pg:"device_id,notnull" json:"device_id"`
	PlaceID   string      `pg:"place_id,notnull" json:"place_id"`
	Type      string      `pg:"type,notnull" json:"type"`
	Name      string      `pg:"name" json:"name"`
	InActions bool        `pg:"in_actions" json:"in_actions"`
	CreatedAt models.Time `json:"created_at"`
	UpdatedAt models.Time `json:"updated_at"`
}

func (f *DeviceFavorite) BeforeInsert(ctx context.Context) (context.Context, error) {
	f.CreatedAt = models.Time{Time: time.Now()}
	f.UpdatedAt = models.Time{Time: time.Now()}

	return ctx, nil
}

func (f *DeviceFavorite) BeforeUpdate(ctx context.Context) (context.Context, error) {
	f.UpdatedAt = models.Time{Time: time.Now()}
	return ctx, nil
}

type UserFavorite struct {
	tableName struct{} `pg:"maasapi.user_favorite"`
	UserID    int      `pg:"user_id" json:"-"`
	DeviceFavorite
}

func (f *UserFavorite) BeforeInsert(ctx context.Context) (context.Context, error) {
	f.DeviceFavorite.BeforeInsert(ctx)

	return ctx, nil
}

func (f *UserFavorite) BeforeUpdate(ctx context.Context) (context.Context, error) {
	f.DeviceFavorite.BeforeUpdate(ctx)
	return ctx, nil
}
