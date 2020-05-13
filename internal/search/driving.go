package search

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"bitbucket.movista.ru/maas/maasapi/internal/google"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"googlemaps.github.io/maps"
)

func convertDrivingRoutesToDataObject(mroutes []MRoute) (dataObject models.DataObject) {

	if len(mroutes) == 0 {
		return
	}

	routes := make([]models.DataObject, 0, len(mroutes))
	for _, v := range mroutes {
		routes = append(routes, models.DataObject{ObjectId: GOOGLE_DRIVING, Id: v.ID, Data: v})
	}

	dataObject = models.DataObject{
		ObjectId: GROUP_OF_DRIVING,
		Data:     models.RouteObject{Routes: routes},
	}

	return
}

func findDrivingRoutes(ctx common.Context, from, to string, origin, destination models.GeoPoint, departureTime, arrivalTime models.Time, region string) (result []MRoute) {
	response, err := google.FindDrivingRoutes(ctx, from, to, maps.TravelModeDriving, maps.TrafficModelBestGuess, true, departureTime, arrivalTime)
	if err != nil {
		return
	}

	if len(response) == 0 {
		return
	}

	result = make([]MRoute, 0)

	for _, route := range response {
		leg := route.Legs[0] //TODO ? что делать, если будет несколько точек

		deeplinkYN := models.DeepLink{
			Title:    "Yandex Navigator",
			Ios:      fmt.Sprintf("yandexnavi://build_route_on_map?lat_from=%f&lon_from=%f&lat_to=%f&lon_to=%f", origin.Latitude, origin.Longitude, destination.Latitude, destination.Longitude),
			Android:  fmt.Sprintf("yandexnavi://build_route_on_map?lat_from=%f&lon_from=%f&lat_to=%f&lon_to=%f", origin.Latitude, origin.Longitude, destination.Latitude, destination.Longitude),
			IconName: "app_icon_yandex_navi",
		}

		deeplinkYM := models.DeepLink{
			Title:    "Yandex Maps",
			Ios:      fmt.Sprintf("yandexmaps://maps.yandex.ru/?rtext=%f,%f~%f,%f&rtt=auto", origin.Latitude, origin.Longitude, destination.Latitude, destination.Longitude),
			Android:  fmt.Sprintf("yandexmaps://maps.yandex.ru/?rtext=%f,%f~%f,%f&rtt=auto", origin.Latitude, origin.Longitude, destination.Latitude, destination.Longitude),
			IconName: "app_icon_yandex_maps",
		}

		deeplinkGM := models.DeepLink{
			Title:    "Google Maps",
			Ios:      fmt.Sprintf("comgooglemaps://?saddr=%f,%f&daddr=%f,%f&directionsmode=driving", origin.Latitude, origin.Longitude, destination.Latitude, destination.Longitude),
			Android:  fmt.Sprintf("https://www.google.com/maps/dir/?api=1&origin=%f,%f&destination=%f,%f", origin.Latitude, origin.Longitude, destination.Latitude, destination.Longitude),
			IconName: "app_icon_google_maps",
		}

		deeplinkAM := models.DeepLink{
			Title: "Apple Maps",
			Ios:   fmt.Sprintf("http://maps.apple.com/?saddr=%f,%f&daddr=%f,%f&z=10&t=s", origin.Latitude, origin.Longitude, destination.Latitude, destination.Longitude),
		}

		startTime := models.NewTime(time.Now())
		endTime := startTime.Add(time.Duration(int64(leg.Duration)) * time.Second)

		mroute := MRoute{
			ID:                uuid.New().String(),
			FromLocation:      latLngToGeoPoint(leg.StartLocation),
			FromAddress:       leg.StartAddress,
			ToLocation:        latLngToGeoPoint(leg.EndLocation),
			ToAddress:         leg.EndAddress,
			Duration:          int64(leg.Duration.Seconds()), //?
			DurationInTraffic: int64(leg.DurationInTraffic.Seconds()),
			StartTime:         startTime,
			EndTime:           &endTime,
			Distance:          int64(leg.Distance.Meters),
			Polyline:          route.OverviewPolyline.Points,
			DeepLinks:         []models.DeepLink{deeplinkYN, deeplinkYM, deeplinkGM, deeplinkAM},
			Summary:           route.Summary,
			IosIconUrl:        config.C.Icons.Driving.IconUrl.Ios,
			AndroidIconUrl:    config.C.Icons.Driving.IconUrl.Android,
		}

		if departureTime.Unix() > 0 {
			mroute.StartTime = &departureTime
			var endTime models.Time
			if mroute.DurationInTraffic != 0 {
				endTime = departureTime.Add(time.Duration(mroute.DurationInTraffic) * time.Second)
			} else {
				endTime = departureTime.Add(time.Duration(mroute.Duration) * time.Second)
			}

			mroute.EndTime = &endTime
		}

		result = append(result, mroute)
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].DurationInTraffic != 0 && result[j].DurationInTraffic != 0 {
			return result[i].DurationInTraffic < result[j].DurationInTraffic
		}
		return result[i].Duration < result[j].Duration
	})

	return
}
