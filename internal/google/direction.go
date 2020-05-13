package google

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"context"
	"fmt"
	"googlemaps.github.io/maps"
)

func FindPublicRoutes(ctx common.Context, from, to string, transitModes []maps.TransitMode, mode maps.Mode, departureTime, arrivalTime models.Time) (routes []maps.Route, err error) {
	if len(transitModes) == 0 {
		transitModes = []maps.TransitMode{maps.TransitModeTrain, maps.TransitModeSubway, maps.TransitModeBus, maps.TransitModeTram}
	}

	routesRequest := &maps.DirectionsRequest{
		Origin:                   from,
		Destination:              to,
		Mode:                     mode,
		Alternatives:             true,
		TransitRoutingPreference: maps.TransitRoutingPreferenceFewerTransfers,
		TransitMode:              transitModes,
		Language:                 "ru",
	}

	if arrivalTime.Unix() > 0 {
		routesRequest.ArrivalTime = fmt.Sprintf("%d", arrivalTime.Unix())
	} else if departureTime.Unix() > 0 {
		routesRequest.DepartureTime = fmt.Sprintf("%d", departureTime.Unix())
	}

	routes, _, err = mapsClient.Directions(context.Background(), routesRequest)
	if err != nil {
		return
	}

	return
}

func FindDrivingRoutes(ctx common.Context, from, to string, mode maps.Mode, model maps.TrafficModel, alternatives bool, departureTime, arrivalTime models.Time) (routes []maps.Route, err error) {
	routesRequest := &maps.DirectionsRequest{
		Origin:       from,
		Destination:  to,
		Mode:         mode,
		Alternatives: alternatives,
		Language:     "ru",
	}

	if arrivalTime.Unix() > 0 {
		// Это такой прикол гугла, что при указании arrival_time не возвращается duration_in_traffic
		// и это время вообще не точное, поэтому вместо указания arrival_time,
		// указываем departure_time и пессиместическую traffic_model
		// для получения максимально возможного времени в пути
		routesRequest.DepartureTime = fmt.Sprintf("%d", arrivalTime.Unix())
		routesRequest.TrafficModel = maps.TrafficModelPessimistic
	} else if departureTime.Unix() > 0 {
		routesRequest.DepartureTime = fmt.Sprintf("%d", departureTime.Unix())
		routesRequest.TrafficModel = model
	}

	routes, _, err = mapsClient.Directions(context.Background(), routesRequest)
	if err != nil {
		return
	}

	return
}

func FindWalkingRoutes(ctx common.Context, from, to string, mode maps.Mode, departureTime, arrivalTime models.Time) (routes []maps.Route, err error) {
	routesRequest := &maps.DirectionsRequest{
		Origin:       from,
		Destination:  to,
		Mode:         mode,
		Alternatives: false,
		Language:     "ru",
	}

	if arrivalTime.Unix() > 0 {
		routesRequest.ArrivalTime = fmt.Sprintf("%d", arrivalTime.Unix())
	} else if departureTime.Unix() > 0 {
		routesRequest.DepartureTime = fmt.Sprintf("%d", departureTime.Unix())
	}

	routes, _, err = mapsClient.Directions(context.Background(), routesRequest)
	if err != nil {
		return
	}

	return
}
