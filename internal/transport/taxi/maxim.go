package taxi

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"math/rand"
	"net/url"
	"sort"
	"strconv"
	"time"
)

const (
	MAXIM_ICON         = "icon_maxim"
	MAXIM_DISPLAY_NAME = "maxim"
)

func CalculateRouteMaxim(ctx common.Context, origin, destination models.GeoPoint) (trip *Trip, err error) {

	resp, err := getMaximRoute(ctx, origin, destination)
	if err != nil {
		return
	}

	rand.Seed(time.Now().UnixNano())
	refOrderId := strconv.Itoa(rand.Intn(999999999-100000000) + 100000000)

	sort.Slice(resp, func(i, j int) bool {
		return resp[i].Price < resp[j].Price
	})

	trip = &Trip{
		ID:                  refOrderId,
		ObjectId:            TAXI_ROUTE,
		MinPrice:            resp[0].Price,
		FromLocation:        &origin,
		ToLocation:          &destination,
		ProviderDescription: MAXIM_DISPLAY_NAME,
		ProviderIcon:        MAXIM_ICON,
		IosIconUrl:          config.C.Icons.Taxi.Maxim.Ios,
		AndroidIconUrl:      config.C.Icons.Taxi.Maxim.Android,
		Fares:               make([]Fare, 0),
	}

	for _, fare := range resp {
		if fare.Price != 0 {
			trip.Fares = append(trip.Fares, Fare{
				ID:   cast.ToString(fare.TariffTypeId),
				Name: fare.TariffTypeName, Price: cast.ToString(fare.Price) + "₽",
			})
		}
	}

	if len(trip.Fares) == 0 {
		err = errors.Errorf("missing any tariffs with price")
		return
	}

	return
}

func makeMaximDeepLink(refOrderId string, origin, destination models.GeoPoint, originAddress, destinationAddress string) models.DeepLink {

	var (
		linkIOS        = config.C.Taxi.Maxim.DeepLink.TrackingLinkIos
		androidAndroid = config.C.Taxi.Maxim.DeepLink.TrackingLinkAndroid
		query          = fmt.Sprintf("refOrgId=%s&refOrderId=%s&startLatitude=%f&startLongitude=%f&&endLatitude=%f&endLongitude=%f&startAddressName=%s&endAddressName=%s",
			config.C.Taxi.Maxim.DeepLink.RefOrgID,
			refOrderId,
			origin.Latitude,
			origin.Longitude,
			destination.Latitude,
			destination.Longitude,
			originAddress,
			destinationAddress)
	)

	deepLinkQuery, _ := url.Parse(config.C.Taxi.Maxim.DeepLink.Prefix + query)
	deepLinkRaw := config.C.Taxi.Maxim.DeepLink.Prefix + deepLinkQuery.Query().Encode()

	return models.DeepLink{
		Id:           refOrderId,
		Ios:          deepLinkRaw,
		Android:      deepLinkRaw,
		IosStore:     linkIOS,
		AndroidStore: androidAndroid,
	}
}

func getMaximRoute(ctx common.Context, origin, destination models.GeoPoint) (response []*MaximTariff, err error) {
	req := common.DefaultRequest.Clone().
		Get(config.C.Taxi.Maxim.Urls.Price).
		Param("access-token", config.C.Taxi.Maxim.Token).
		Param("startLatitude", cast.ToString(origin.Latitude)).
		Param("startLongitude", cast.ToString(origin.Longitude)).
		Param("endLatitude", cast.ToString(destination.Latitude)).
		Param("endLongitude", cast.ToString(destination.Longitude)).
		Param("appcode", config.C.Taxi.Maxim.AppCode)

	body, _, err := common.SendRequest(ctx, req)
	if err != nil {
		err = fmt.Errorf("maxim request failed - %w", err)
		return
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		err = fmt.Errorf("cannot unmarshal maxim response - %s", string(body))
		return
	}

	if len(response) == 0 {
		err = fmt.Errorf("response from maxim contains 0 tariffs - %s", string(body))
		return
	}

	return
}
