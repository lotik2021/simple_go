package weather

import (
	"fmt"

	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
)

func getByLocation(ctx common.Context, location models.GeoPoint, what string) (resp []byte, err error) {
	url := config.C.OpenWeather.BaseURL + fmt.Sprintf("/%s?lat=%f&lon=%f&units=%s&lang=%s&appid=%s", what, location.Latitude, location.Longitude, units, lang, config.C.OpenWeather.ApiKey)

	req := wclient.Clone().Get(url)

	resp, _, err = common.SendRequest(ctx, req)
	if err != nil {
		return
	}

	return
}

func GetWeather(ctx common.Context, location models.GeoPoint) ([]byte, error) {
	return getByLocation(ctx, location, "weather")
}

func GetForecast(ctx common.Context, location models.GeoPoint) ([]byte, error) {
	return getByLocation(ctx, location, "forecast")
}
