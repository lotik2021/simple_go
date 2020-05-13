package taxi

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"bitbucket.movista.ru/maas/maasapi/internal/logger"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	geo "github.com/kellydunn/golang-geo"
	"github.com/parnurzeal/gorequest"
	"net/http"
	"sort"
)

const (
	CITYMOBILE_ICON         = "icon_citymobil"
	CITYMOBILE_DISPLAY_NAME = "Ситимобил"
)

var (
	moscow = geo.NewPolygon([]*geo.Point{
		geo.NewPoint(56.419255, 35.485840),
		geo.NewPoint(56.562081, 36.870117),
		geo.NewPoint(56.936313, 37.694092),
		geo.NewPoint(56.773801, 38.265381),
		geo.NewPoint(55.978291, 38.551025),
		geo.NewPoint(55.758993, 39.391479),
		geo.NewPoint(55.796145, 39.825439),
		geo.NewPoint(55.313525, 40.144043),
		geo.NewPoint(54.311060, 38.649902),
		geo.NewPoint(54.834445, 38.045654),
		geo.NewPoint(54.897812, 37.221680),
		geo.NewPoint(55.216285, 36.996460),
		geo.NewPoint(55.341711, 36.441650),
		geo.NewPoint(55.206862, 36.205444),
		geo.NewPoint(55.297857, 35.332031),
	})

	yaroslavl = geo.NewPolygon([]*geo.Point{
		geo.NewPoint(58.362963, 37.496338),
		geo.NewPoint(58.900069, 39.067383),
		geo.NewPoint(58.381721, 41.066895),
		geo.NewPoint(56.601491, 39.171753),
		geo.NewPoint(56.834074, 38.226929),
		geo.NewPoint(57.378094, 38.358765),
		geo.NewPoint(57.466999, 37.825928),
		geo.NewPoint(57.975471, 37.386475),
		geo.NewPoint(58.106560, 37.749023),
		geo.NewPoint(58.303736, 37.380981),
	})

	tolyatti = geo.NewPolygon([]*geo.Point{
		geo.NewPoint(53.590348, 49.214630),
		geo.NewPoint(53.551929, 49.611511),
		geo.NewPoint(53.468430, 49.664383),
		geo.NewPoint(53.513884, 49.195404),
	})

	samara = geo.NewPolygon([]*geo.Point{
		geo.NewPoint(54.640560, 51.361084),
		geo.NewPoint(54.320695, 52.470703),
		geo.NewPoint(51.827185, 50.784302),
		geo.NewPoint(53.355814, 47.977295),
		geo.NewPoint(53.813424, 50.026245),
		geo.NewPoint(54.490537, 50.130615),
	})
)

func CalculateRouteCityMobile(ctx common.Context, origin, destination models.GeoPoint) (trip *Trip, err error) {

	resp, err := getCityMobileRoute(ctx, origin, destination)
	if err != nil {
		return
	}

	// сотировать по минимальной цене
	sort.Slice(resp.Tariffs, func(i, j int) bool {
		return resp.Tariffs[i].TotalPrice < resp.Tariffs[j].TotalPrice
	})

	trip = &Trip{
		ID:                  uuid.New().String(),
		ObjectId:            TAXI_ROUTE,
		MinPrice:            resp.Tariffs[0].TotalPrice,
		FromLocation:        &origin,
		ToLocation:          &destination,
		ProviderDescription: CITYMOBILE_DISPLAY_NAME,
		ProviderIcon:        CITYMOBILE_ICON,
		IosIconUrl:          config.C.Icons.Taxi.Citymobil.Ios,
		AndroidIconUrl:      config.C.Icons.Taxi.Citymobil.Android,
		Fares:               make([]Fare, 0),
		CalculationID:       resp.ID,
	}

	for _, fare := range resp.Tariffs {
		trip.Fares = append(trip.Fares, Fare{
			ID:    fmt.Sprintf("%d", fare.ID),
			Name:  fare.Info.Name,
			Price: fmt.Sprintf("%d₽", int(fare.TotalPrice)),
		})
	}

	trip.DeepLink = makeCityMobileDeepLink(trip.ID, origin, destination)

	return
}

func makeCityMobileDeepLink(refOrderId string, origin, destination models.GeoPoint) models.DeepLink {
	var (
		trackingLinkPrefix string
		city               string
		originPoint        = geo.NewPoint(origin.Latitude, origin.Longitude)
		query              = fmt.Sprintf("from=%f,%f&from_str=&to=%f,%f&to_str=&oid=%s&tariff=644&partner_id=%s",
			origin.Latitude,
			origin.Longitude,
			destination.Latitude,
			destination.Longitude,
			refOrderId,
			config.C.Taxi.CityMobile.DeepLink.PartnerId)
	)

	if tolyatti.Contains(originPoint) {
		city = "tolyatti"
		trackingLinkPrefix = config.C.Taxi.CityMobile.DeepLink.TrackingLinkTolyatti
	} else if samara.Contains(originPoint) {
		city = "samara"
		trackingLinkPrefix = config.C.Taxi.CityMobile.DeepLink.TrackingLinkSamara
	} else if yaroslavl.Contains(originPoint) {
		city = "yaroslavl"
		trackingLinkPrefix = config.C.Taxi.CityMobile.DeepLink.TrackingLinkYaroslavl
	} else if moscow.Contains(originPoint) {
		city = "moscow"
		trackingLinkPrefix = config.C.Taxi.CityMobile.DeepLink.TrackingLinkMoscow
	} else {
		city = "unknown"
		trackingLinkPrefix = config.C.Taxi.CityMobile.DeepLink.TrackingLinkMoscow
	}

	logger.Log.WithFields(logger.Fields{
		"city":               city,
		"originLatLng":       fmt.Sprintf("%f %f", origin.Latitude, origin.Longitude),
		"trackingLinkPrefix": trackingLinkPrefix,
	}).Info()

	deeplinkUrl := config.C.Taxi.CityMobile.DeepLink.Prefix + query
	trackingLink := trackingLinkPrefix + query

	return models.DeepLink{
		Id:           refOrderId,
		Ios:          deeplinkUrl,
		Android:      deeplinkUrl,
		IosStore:     trackingLink,
		AndroidStore: trackingLink,
	}
}

func getCityMobileRoute(ctx common.Context, origin, destination models.GeoPoint) (resp *CitymobileResponse, err error) {
	tariffs := []int{2, 4, 5, 7}

	citymobileRequest := struct {
		TariffGroup          []int   `json:"tariff_group"`
		OriginLatitude       float64 `json:"latitude"`
		OriginLongitude      float64 `json:"longitude"`
		DestinationLatitude  float64 `json:"del_latitude"`
		DestinationLongitude float64 `json:"del_longitude"`
	}{
		TariffGroup:          tariffs,
		OriginLatitude:       origin.Latitude,
		OriginLongitude:      origin.Longitude,
		DestinationLatitude:  destination.Latitude,
		DestinationLongitude: destination.Longitude,
	}

	req := common.DefaultRequest.Clone().
		Post(config.C.Taxi.CityMobile.Urls.Price).
		Set("Authorization", cityMobileToken).
		Type(gorequest.TypeJSON).
		SendStruct(&citymobileRequest)

	body, response, err := common.SendRequest(ctx, req)
	if err != nil {
		if response.StatusCode == http.StatusUnauthorized {
			err = authCityMobile()
			if err != nil {
				return
			}
			return getCityMobileRoute(ctx, origin, destination)
		}
		err = fmt.Errorf("citymobile request failed - %w", err)
		return
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		err = fmt.Errorf("cannot unmarshal maxim response - %s", string(body))
		return
	}

	if resp.Code != 0 {
		err = fmt.Errorf("error code in citymobile response - %d, response - %s", resp.Code, string(body))
		return
	}

	return
}

func authCityMobile() (err error) {
	ctx := common.NewInternalContext()

	authRequest := struct {
		ApiKey string `json:"api_key"`
	}{
		ApiKey: config.C.Taxi.CityMobile.ApiKey,
	}

	var authResponse struct {
		Token string `json:"token"`
		Code  int    `json:"code"`
	}

	req := common.DefaultRequest.Clone().
		Post(config.C.Taxi.CityMobile.Urls.Auth).
		Type(gorequest.TypeJSON).
		SendStruct(&authRequest)

	body, _, err := common.SendRequest(ctx, req)
	if err != nil {
		err = fmt.Errorf("citymobile auth error - %w", err)
		return
	}

	err = json.Unmarshal(body, &authResponse)
	if err != nil {
		err = fmt.Errorf("cannot unmarshal maxim response - %s", string(body))
		return
	}
	// TODO: handle response status code

	if authResponse.Code != 0 {
		err = fmt.Errorf("citymobile error code - %d, body - %s", authResponse.Code, string(body))
		return
	}

	cityMobileToken = "Bearer " + authResponse.Token

	return
}
