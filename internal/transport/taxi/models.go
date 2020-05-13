package taxi

import (
	"bitbucket.movista.ru/maas/maasapi/internal/models"
)

type Fare struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price string `json:"price"`
}

type Trip struct {
	ID                  string           `json:"-"`
	ObjectId            string           `json:"-"`
	MinPrice            float64          `json:"-"`
	FromLocation        *models.GeoPoint `json:"from_location"`
	FromAddress         string           `json:"from_address"`
	ToLocation          *models.GeoPoint `json:"to_location"`
	ToAddress           string           `json:"to_address"`
	Duration            int64            `json:"duration"`
	StartTime           *models.Time     `json:"start_time"`
	EndTime             *models.Time     `json:"end_time"`
	ProviderDescription string           `json:"provider_description"`
	ProviderIcon        string           `json:"provider_icon"`
	IosIconUrl          string           `json:"ios_icon_url"`
	AndroidIconUrl      string           `json:"android_icon_url"`
	Fares               []Fare           `json:"fares"`
	Polyline            string           `json:"polyline"`
	Distance            int64            `json:"distance"`
	CalculationID       string           `json:"calculation_id"`
	DeepLink            models.DeepLink  `json:"deeplinks"`
}

type YandexResponse struct {
	Tariffs []*YandexOptions `json:"options"`
}

type YandexOptions struct {
	ClassLevel   int64   `json:"class_level"`
	ClassName    string  `json:"class_name"`
	ClassText    string  `json:"class_text"`
	MinimumPrice float64 `json:"min_price"`
	Price        float64 `json:"price"`
	PriceText    string  `json:"price_text"`
	WaitingTime  float64 `json:"waiting_time"`
}

type MaximTariff struct {
	Price          float64 `json:"Price"`
	PriceString    string  `json:"PriceString"`
	FeedTime       int64   `json:"FeedTime"`
	CurrencySymbol string  `json:"CurrencySymbol"`
	TariffTypeName string  `json:"TariffTypeName"`
	TariffTypeId   int     `json:"TariffTypeId"`
}

type CitymobileResponse struct {
	ID      string              `json:"id_calculation"`
	Code    int                 `json:"code"`
	Message string              `json:"message"`
	Tariffs []*CitymobileTariff `json:"prices"`
}

type CitymobileTariff struct {
	ID         int                   `json:"id_tariff"`
	TotalPrice float64               `json:"total_price"`
	Info       *CitymobileTariffInfo `json:"tariff_info"`
}

type CitymobileTariffInfo struct {
	CarsDescription         string `json:"car_models"`
	Name                    string `json:"name"`
	CarsCapacityDescription string `json:"car_capacity"`
}
