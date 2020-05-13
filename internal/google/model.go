package google

import (
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"context"
	"time"
)

type TransitModeName struct {
	tableName      struct{} `pg:"maasapi.google_transit_mode_name,alias:google_transit_mode_name"`
	Name           string   `pg:"name,pk" json:"name"`
	DisplayName    string   `pg:"display_name,notnull" json:"transport_name"`
	IconName       string   `pg:"icon_name" json:"icon_name"`
	IosIconUrl     string   `pg:"ios_icon_url" json:"ios_icon_url"`
	AndroidIconUrl string   `pg:"android_icon_url" json:"android_icon_url"`
}

type ShortPlace struct {
	PlaceID     string           `json:"place_id"`
	Description string           `json:"description"`
	Location    *models.GeoPoint `json:"location"`
}

type Place struct {
	tableName            struct{}         `pg:"maasapi.google_place,alias:gp,discard_unknown_columns"`
	ID                   string           `pg:"id,pk" json:"place_id"`
	MainText             string           `pg:"main_text" json:"main_text"`
	SecondaryText        string           `pg:"secondary_text" json:"secondary_text"`
	Coordinate           *models.GeoPoint `pg:"coordinate,notnull,type:geometry(Geometry,4326)" json:"location"`
	PlaceTypes           []string         `pg:"place_types,array" json:"types"`
	IosIconURL           string           `pg:"-" json:"ios_icon_url"` // tmp fix for Android 2.1
	IosIconURLLightTheme string           `pg:"-" json:"ios_icon_url_dark"`
	IosIconURLDarkTheme  string           `pg:"-" json:"ios_icon_url_light"`
	AndroidIconURL       string           `pg:"-" json:"android_icon_url"`
	Type                 string           `pg:"-" json:"type"`
	Name                 string           `pg:"-" json:"name"`
	CreatedAt            models.Time      `pg:"created_at" json:"created_at"`
	UpdatedAt            models.Time      `pg:"updated_at" json:"updated_at"`
}

func (p *Place) BeforeInsert(ctx context.Context) (context.Context, error) {
	p.CreatedAt = models.Time{Time: time.Now()}
	p.UpdatedAt = models.Time{Time: time.Now()}

	return ctx, nil
}

func (p *Place) BeforeUpdate(ctx context.Context) (context.Context, error) {
	p.UpdatedAt = models.Time{Time: time.Now()}
	return ctx, nil
}

func (p *Place) GetIcons() {
	for _, t := range p.PlaceTypes {
		if t == "subway_station" {
			p.IosIconURLLightTheme = config.C.Icons.Places.Metro.Ios
			p.IosIconURLDarkTheme = config.C.Icons.Places.DarkTheme.Metro.Ios
			p.AndroidIconURL = config.C.Icons.Places.Metro.Android
			p.IosIconURL = p.IosIconURLDarkTheme
			return
		} else if t == "train_station" {
			p.IosIconURLLightTheme = config.C.Icons.Places.Train.Ios
			p.IosIconURLDarkTheme = config.C.Icons.Places.DarkTheme.Train.Ios
			p.AndroidIconURL = config.C.Icons.Places.Train.Android
			p.IosIconURL = p.IosIconURLDarkTheme
			return
		} else if t == "bus_station" {
			p.IosIconURLLightTheme = config.C.Icons.Places.Bus.Ios
			p.IosIconURLDarkTheme = config.C.Icons.Places.DarkTheme.Bus.Ios
			p.AndroidIconURL = config.C.Icons.Places.Bus.Android
			p.IosIconURL = p.IosIconURLDarkTheme
			return
		} else if t == "airport" {
			p.IosIconURLLightTheme = config.C.Icons.Places.Airport.Ios
			p.IosIconURLDarkTheme = config.C.Icons.Places.DarkTheme.Airport.Ios
			p.AndroidIconURL = config.C.Icons.Places.Airport.Android
			p.IosIconURL = p.IosIconURLDarkTheme
			return
		} else if t == "transit_station" {
			p.IosIconURLLightTheme = config.C.Icons.Places.Hub.Ios
			p.IosIconURLDarkTheme = config.C.Icons.Places.DarkTheme.Hub.Ios
			p.AndroidIconURL = config.C.Icons.Places.Hub.Android
			p.IosIconURL = p.IosIconURLDarkTheme
			return
		} else if t == "locality" {
			p.IosIconURLLightTheme = config.C.Icons.Places.City.Ios
			p.IosIconURLDarkTheme = config.C.Icons.Places.DarkTheme.City.Ios
			p.AndroidIconURL = config.C.Icons.Places.City.Android
			p.IosIconURL = p.IosIconURLDarkTheme
			return
		}
	}

	// Сюда попадём если не нашли тип ранее
	p.IosIconURLLightTheme = config.C.Icons.Places.Default.Ios
	p.IosIconURLDarkTheme = config.C.Icons.Places.DarkTheme.Default.Ios
	p.AndroidIconURL = config.C.Icons.Places.Default.Android
	p.IosIconURL = p.IosIconURLDarkTheme

	return
}
