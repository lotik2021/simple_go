package search

import (
	"bitbucket.movista.ru/maas/maasapi/internal/models"
)

type MRoute struct {
	ID                string              `json:"-"`
	StartTime         *models.Time        `json:"start_time"`
	EndTime           *models.Time        `json:"end_time"`
	Distance          int64               `json:"distance"`
	Duration          int64               `json:"duration"`
	DurationInTraffic int64               `json:"duration_in_traffic"`
	FromAddress       string              `json:"from_address"`
	FromLocation      models.GeoPoint     `json:"from_location"`
	ToAddress         string              `json:"to_address"`
	ToLocation        models.GeoPoint     `json:"to_location"`
	Agencies          []Agency            `json:"agencies"`
	IosIconUrl        string              `json:"ios_icon_url"`
	AndroidIconUrl    string              `json:"android_icon_url"`
	Fare              interface{}         `json:"fare"`
	Trips             []models.DataObject `json:"trips"`
	Polyline          string              `json:"polyline"`
	DeepLinks         []models.DeepLink   `json:"deeplinks"`
	Summary           string              `json:"summary"`
}

type Agency struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Url   string `json:"url"`
}

type GenericDataInfo struct {
	Color              string `json:"color"`
	Title              string `json:"title"`
	ShortDescription   string `json:"short_description"`
	FullDescription    string `json:"full_description"`
	IconDetails        string `json:"icon_details"`
	IconShort          string `json:"icon_short"`
	AndroidIconDetails string `json:"android_icon_details"`
	AndroidIconShort   string `json:"android_icon_short"`
}

type Trip struct {
	ObjectID      string           `json:"-"`
	SourceCode    string           `json:"sourceCode"`
	FromLocation  models.GeoPoint  `json:"from_location"`
	DepartureStop *string          `json:"departure_stop,omitempty"`
	DepartureTime *models.Time     `json:"departure_time,omitempty"`
	ToLocation    models.GeoPoint  `json:"to_location"`
	ArrivalStop   *string          `json:"arrival_stop,omitempty"`
	ArrivalTime   *models.Time     `json:"arrival_time,omitempty"`
	Duration      int64            `json:"duration"`
	Polyline      string           `json:"polyline"`
	Distance      int64            `json:"distance"`
	TripType      string           `json:"trip_type"`
	LineNumber    *string          `json:"line_number,omitempty"`
	LineName      *string          `json:"line_name,omitempty"`
	LineColor     *string          `json:"line_color,omitempty"`
	Headsign      *string          `json:"headsign,omitempty"`
	GenericData   *GenericDataInfo `json:"generic_data"`
	NumberStops   *uint            `json:"number_stops,omitempty"`
	Agencies      *[]Agency        `json:"agencies,omitempty"`
	Url           string           `json:"url,omitempty"`
	TripShortName *string          `json:"trip_short_name,omitempty"`
}
