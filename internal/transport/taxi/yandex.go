package taxi

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"net/url"
	"sort"
)

const (
	YANDEX_ICON         = "icon_yandex"
	YANDEX_DISPLAY_NAME = "Яндекс.Такси"
)

func CalculateRouteYandex(ctx common.Context, origin, destination models.GeoPoint) (trip *Trip, err error) {

	resp, err := getYandexRoute(ctx, origin, destination)
	if err != nil {
		return
	}

	sort.Slice(resp.Tariffs, func(i, j int) bool {
		return resp.Tariffs[i].Price < resp.Tariffs[j].Price
	})

	trip = &Trip{
		ID:                  uuid.New().String(),
		ObjectId:            TAXI_ROUTE,
		MinPrice:            resp.Tariffs[0].Price,
		FromLocation:        &origin,
		ToLocation:          &destination,
		ProviderDescription: YANDEX_DISPLAY_NAME,
		ProviderIcon:        YANDEX_ICON,
		IosIconUrl:          config.C.Icons.Taxi.Yandex.Ios,
		AndroidIconUrl:      config.C.Icons.Taxi.Yandex.Android,
		Fares:               make([]Fare, 0),
	}

	for _, fare := range resp.Tariffs {
		if fare.Price != 0 {
			trip.Fares = append(trip.Fares, Fare{
				ID:    cast.ToString(fare.ClassLevel),
				Name:  fare.ClassText,
				Price: fare.PriceText,
			})
		}
	}

	trip.DeepLink = makeYandexDeepLink(trip.ID, origin, destination)

	return
}

func makeYandexDeepLink(refOrderId string, origin, destination models.GeoPoint) models.DeepLink {
	//https://taxi-routeinfo.taxi.yandex.net/taxi_info?clid=movista&apikey=3b320e58d9fb4fb2a9bf5a04dde8c8ca&rll=37.611717,55.742268~37.679609,55.722653&class=econom,business,comfortplus,minivan,vip&req=&lang=ru
	//https://3.redirect.appmetrica.yandex.com/route?start-lat=<широта>&start-lon=<долгота>&end-lat=<широта>&end-lon=<долгота>&level=<тариф>&ref=<источник>&appmetrica_tracking_id=1178268795219780156

	var (
		query = fmt.Sprintf("start-lat=%f&start-lon=%f&end-lat=%f&end-lon=%f&ref=%s&appmetrica_tracking_id=%s",
			origin.Latitude,
			origin.Longitude,
			destination.Latitude,
			destination.Longitude,
			config.C.Taxi.Yandex.RefId,
			config.C.Taxi.Yandex.DeepLink.TrackingId)
	)

	deepLinkUrl, _ := url.Parse(config.C.Taxi.Yandex.DeepLink.Prefix + query)

	return models.DeepLink{
		Id:           refOrderId,
		Ios:          deepLinkUrl.String(),
		Android:      deepLinkUrl.String(),
		IosStore:     deepLinkUrl.String(),
		AndroidStore: deepLinkUrl.String(),
	}
}

func getYandexRoute(ctx common.Context, origin, destination models.GeoPoint) (resp *YandexResponse, err error) {
	//37.611717,55.742268~37.679609,55.722653
	rll := fmt.Sprintf("%f,%f~%f,%f", origin.Longitude, origin.Latitude, destination.Longitude, destination.Latitude)

	priceRequest := common.DefaultRequest.Clone().
		Get(config.C.Taxi.Yandex.Urls.Price).
		Param("clid", config.C.Taxi.Yandex.Clid).
		Param("rll", rll).
		Param("class", "econom,business,comfortplus,minivan,vip").
		Param("lang", "ru").
		Param("apikey", config.C.Taxi.Yandex.ApiKey)

	body, _, err := common.SendRequest(ctx, priceRequest)
	if err != nil {
		err = fmt.Errorf("yandex request failed - %w", err)
		return
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		err = fmt.Errorf("cannot unmarshal yandex response - %s", string(body))
		return
	}

	if len(resp.Tariffs) == 0 {
		err = fmt.Errorf("response from yandex contains 0 tariffs - %s", string(body))
		return
	}

	return
}
