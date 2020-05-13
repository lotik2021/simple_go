package device

import (
	"context"
	"time"

	"bitbucket.movista.ru/maas/maasapi/internal/models"
)

type Device struct {
	tableName          struct{}    `pg:"maasapi.device"`
	ID                 string      `pg:"id,pk" json:"id"`
	Name               string      `pg:"name" json:"user_name"`
	LastName           string      `pg:"last_name" json:"last_name"`
	DeviceOS           string      `pg:"device_os" json:"device_os"`
	DeviceCategory     string      `pg:"device_category" json:"device_category"`
	DeviceType         string      `pg:"device_type" json:"device_type"`
	DeviceInfo         string      `pg:"device_info" json:"device_info"`
	OsPlayerID         string      `pg:"os_player_id" json:"os_player_id"`
	GoogleTransitModes []string    `pg:"google_transit_modes,array" json:"google_transit_modes"`
	NameAskedAt        models.Time `pg:"name_asked_at" json:"name_asked_at"`
	SessionCount       int         `pg:"session_count" json:"session_count"`
	CreatedAt          models.Time `json:"created_at"`
	UpdatedAt          models.Time `json:"updated_at"`
	DeletedAt          time.Time   `pg:",soft_delete"`
}

func (d *Device) BeforeInsert(ctx context.Context) (context.Context, error) {
	d.CreatedAt = models.Time{Time: time.Now()}
	d.UpdatedAt = models.Time{Time: time.Now()}

	d.GoogleTransitModes = []string{"bus", "train", "tram", "subway"}

	return ctx, nil
}

func (d *Device) BeforeUpdate(ctx context.Context) (context.Context, error) {
	d.UpdatedAt = models.Time{Time: time.Now()}
	return ctx, nil
}

type Location struct {
	tableName  struct{}         `pg:"maasapi.device_location"`
	DeviceID   string           `pg:"device_id,notnull,pk"`
	Coordinate *models.GeoPoint `pg:"coordinate,notnull,type:geometry(Geometry,4326)"`
	CreatedAt  models.Time      `pg:",pk"`
}

func (l *Location) BeforeInsert(ctx context.Context) (context.Context, error) {
	l.CreatedAt = models.Time{Time: time.Now()}
	return ctx, nil
}

type TaxiOrder struct {
	tableName struct{}    `pg:"maasapi.device_taxi_order,alias:device_taxi_order"`
	Id        string      `pg:"id,pk" json:"id"`
	DeviceID  string      `pg:"device_id,notnull" json:"device_id"`
	Provider  string      `pg:"provider,notnull" json:"provider"`
	UsedLink  string      `pg:"used_link" json:"used_link"`
	CreatedAt models.Time `json:"created_at"`
	UpdatedAt models.Time `json:"updated_at"`
}

func (o *TaxiOrder) BeforeInsert(ctx context.Context) (context.Context, error) {
	o.CreatedAt = models.Time{Time: time.Now()}
	o.UpdatedAt = models.Time{Time: time.Now()}
	return ctx, nil
}

func (o *TaxiOrder) BeforeUpdate(ctx context.Context) (context.Context, error) {
	o.UpdatedAt = models.Time{Time: time.Now()}
	return ctx, nil
}

type TravelCardPayment struct {
	tableName  struct{}    `pg:"maasapi.device_travel_card_payment,alias:device_travel_card_payment"`
	ID         string      `pg:"id,pk" json:"id"`
	DeviceID   string      `pg:"device_id,notnull" json:"device_id"`
	Amount     float64     `pg:"amount,notnull" json:"amount"`
	CardNumber string      `pg:"card_number" json:"card_number"`
	CardType   string      `pg:"card_type,notnull" json:"card_type"`
	Processed  bool        `pg:"processed" json:"processed"`
	Currency   string      `pg:"currency" json:"currency"`
	CreatedAt  models.Time `json:"created_at"`
	UpdatedAt  models.Time `json:"updated_at"`
}

func (p *TravelCardPayment) BeforeInsert(ctx context.Context) (context.Context, error) {
	p.CreatedAt = models.Time{Time: time.Now()}
	p.UpdatedAt = models.Time{Time: time.Now()}
	return ctx, nil
}

func (p *TravelCardPayment) BeforeUpdate(ctx context.Context) (context.Context, error) {
	p.UpdatedAt = models.Time{Time: time.Now()}
	return ctx, nil
}

type RegenerateCodeAttempt struct {
	tableName struct{}    `pg:"maasapi.auth_regenerate_code_attempt"`
	ID        string      `pg:"id,pk" json:"id"`
	DeviceID  string      `pg:"device_id,notnull" json:"device_id"`
	Phone     string      `pg:"phone,notnull" json:"phone"`
	CreatedAt models.Time `json:"created_at"`
	NotBefore models.Time `pg:"nbf" json:"nbf"`
	DeletedAt time.Time   `pg:",soft_delete"`
}

func (p *RegenerateCodeAttempt) BeforeInsert(ctx context.Context) (context.Context, error) {
	p.CreatedAt = models.Time{Time: time.Now()}
	return ctx, nil
}
