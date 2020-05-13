package history

import (
	"context"
	"time"

	"bitbucket.movista.ru/maas/maasapi/internal/models"
)

type Customer struct {
	tableName       struct{} `pg:"maasapi.customer"`
	Id              int      `pg:"id,pk"`
	Age             int      `pg:"age"`
	SeatRequired    bool     `pg:"seat_required"`
	MovistaSearchId string
}

var DefaultCustomer = Customer{Id: 0, Age: 36, SeatRequired: true}

type MovistaSearch struct {
	tableName         struct{}         `pg:"maasapi.device_movista_search_history,alias:device_movista_search_history"`
	ID                string           `pg:"id,pk" json:"id"`
	DeviceID          string           `pg:"device_id,notnull" json:"device_id"`
	FromID            int              `pg:"from_id,notnull" json:"from_id"`
	ToID              int              `pg:"to_id,notnull" json:"to_id"`
	Origin            *models.GeoPoint `pg:"origin,type:geometry(Geometry,4326)" json:"origin"`
	Destination       *models.GeoPoint `pg:"destination,type:geometry(Geometry,4326)" json:"destination"`
	FromGooglePlaceID string           `pg:"from_google_place_id" json:"from_google_place_id"`
	ToGooglePlaceID   string           `pg:"to_google_place_id" json:"to_google_place_id"`
	DepartureTime     models.Time      `pg:"departure_time,notnull" json:"departure_time"`
	ArrivalTime       models.Time      `pg:"arrival_time" json:"arrival_time"`
	TripTypes         []string         `pg:"trip_types,array" json:"trip_types,omitempty"`
	CurrencyCode      string           `pg:"currency_code"`
	CultureCode       string           `pg:"culture_code"`
	ComfortType       string           `pg:"comfort_type"`
	Customers         []Customer
	CreatedAt         models.Time `pg:"created_at" json:"created_at"`
}

func (h *MovistaSearch) BeforeInsert(ctx context.Context) (context.Context, error) {
	h.CreatedAt = models.Time{Time: time.Now()}
	return ctx, nil
}

type GoogleSearch struct {
	tableName          struct{}         `pg:"maasapi.device_google_search_history,alias:device_google_search_history"`
	ID                 string           `pg:"id,pk" json:"id"`
	DeviceID           string           `pg:"device_id,notnull" json:"device_id"`
	Origin             *models.GeoPoint `pg:"origin,type:geometry(Geometry,4326),notnull" json:"origin"`
	Destination        *models.GeoPoint `pg:"destination,type:geometry(Geometry,4326),notnull" json:"destination"`
	OriginPlaceID      string           `pg:"origin_place_id" json:"origin_place_id"`
	DestinationPlaceID string           `pg:"destination_place_id" json:"destination_place_id"`
	DepartureTime      models.Time      `pg:"departure_time" json:"departure_time"`
	ArrivalTime        models.Time      `pg:"arrival_time" json:"arrival_time"`
	TripTypes          []string         `pg:"trip_types,array" json:"trip_types"`
	CreatedAt          models.Time      `pg:"created_at" json:"created_at"`
}

func (h *GoogleSearch) BeforeInsert(ctx context.Context) (context.Context, error) {
	h.CreatedAt = models.Time{Time: time.Now()}
	return ctx, nil
}

type GooglePlaceHistory struct {
	tableName        struct{}    `pg:"maasapi.device_google_place_history,alias:dgph"`
	DeviceID         string      `pg:"device_id,notnull" json:"device_id"`
	PlaceID          string      `pg:"place_id,notnull" json:"coordinate"`
	NumberOfSearches int         `pg:"number_of_searches" json:"number_of_searches"`
	CreatedAt        models.Time `pg:"created_at" json:"created_at"`
	UpdatedAt        models.Time `pg:"updated_at" json:"updated_at"`
}

func (h *GooglePlaceHistory) BeforeInsert(ctx context.Context) (context.Context, error) {
	h.CreatedAt = models.Time{Time: time.Now()}
	h.UpdatedAt = models.Time{Time: time.Now()}
	return ctx, nil
}

func (h *GooglePlaceHistory) BeforeUpdate(ctx context.Context) (context.Context, error) {
	h.UpdatedAt = models.Time{Time: time.Now()}
	return ctx, nil
}
