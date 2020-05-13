package dialog

import (
	"context"
	"time"

	"bitbucket.movista.ru/maas/maasapi/internal/models"
)

type DeviceSession struct {
	tableName          struct{}         `pg:"maasapi.device_session,alias:device_session"`
	ID                 string           `pg:"id,pk" json:"id"`
	DeviceID           string           `pg:"device_id,notnull" json:"device_id"`
	Origin             *models.GeoPoint `pg:"origin,type:geometry(Geometry,4326)" json:"origin"`
	Destination        *models.GeoPoint `pg:"destination,type:geometry(Geometry,4326)" json:"destination"`
	OriginPlaceId      string           `pg:"origin_place_id" json:"origin_place_id"`
	DestinationPlaceId string           `pg:"destination_place_id" json:"destination_place_id"`
	ArrivalTime        *models.Time     `pg:"arrival_time" json:"arrival_time"`
	DepartureTime      *models.Time     `pg:"departure_time" json:"departure_time"`
	State              string           `pg:"state,notnull" json:"state"`
	TripTypes          []string         `pg:"trip_types,array" json:"trip_types,omitempty"`
	CreatedAt          models.Time      `json:"created_at"`
	UpdatedAt          models.Time      `json:"updated_at"`
}

func (l *DeviceSession) BeforeInsert(ctx context.Context) (context.Context, error) {
	l.CreatedAt = models.Time{Time: time.Now()}
	l.UpdatedAt = models.Time{Time: time.Now()}
	return ctx, nil
}

func (l *DeviceSession) BeforeUpdate(ctx context.Context) (context.Context, error) {
	l.UpdatedAt = models.Time{Time: time.Now()}
	return ctx, nil
}

type DeviceSessionData struct {
	tableName      struct{}    `pg:"maasapi.device_session_data,alias:device_session_data"`
	ID             string      `pg:"id,pk" json:"id"`
	SessionID      string      `pg:"session_id,notnull" json:"session_id"`
	ActionID       string      `pg:"action_id" json:"action_id"`
	ActionName     string      `pg:"action_name,notnull" json:"action_name"`
	UserResponse   string      `pg:"user_response" json:"user_response"`
	DialogResponse string      `pg:"dialog_response" json:"dialog_response"`
	UserEntryData  string      `pg:"user_entry_data" json:"user_entry_data,omitempty"`
	Objects        string      `pg:"objects" json:"objects,omitempty"`
	Actions        string      `pg:"actions" json:"actions,omitempty"`
	CreatedAt      models.Time `json:"created_at"`
}

func (s *DeviceSessionData) BeforeInsert(ctx context.Context) (context.Context, error) {
	s.CreatedAt = models.Time{Time: time.Now()}
	return ctx, nil
}
