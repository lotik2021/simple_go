package search

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/google"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"github.com/google/uuid"
	"googlemaps.github.io/maps"
)

func convertWalkingRoutesToDataObject(mroutes []MRoute) (dataObject models.DataObject) {
	if len(mroutes) == 0 {
		return
	}

	routes := make([]models.DataObject, 0, len(mroutes))
	for _, v := range mroutes {
		routes = append(routes, models.DataObject{ObjectId: GOOGLE_ROUTE_TRANSPORT, Id: v.ID, Data: v})
	}

	dataObject = models.DataObject{
		ObjectId: GOOGLE_ROUTE,
		Data:     models.RouteObject{Routes: routes},
	}

	return
}

func findWalkingRoutes(ctx common.Context, from, to string, departureTime, arrivalTime models.Time) (result []MRoute) {
	response, err := google.FindWalkingRoutes(ctx, from, to, maps.TravelModeWalking, departureTime, arrivalTime)
	if err != nil {
		return
	}

	if len(response) == 0 {
		return
	}

	result = make([]MRoute, 0)

	for _, groute := range response {
		leg := groute.Legs[0]

		startTime, endTime := getRouteTimes(leg, departureTime, arrivalTime)

		mroute := MRoute{
			ID:           uuid.New().String(),
			StartTime:    &startTime,
			EndTime:      &endTime,
			Distance:     int64(leg.Meters),
			Duration:     int64(leg.Duration.Minutes()),
			FromAddress:  leg.StartAddress,
			FromLocation: latLngToGeoPoint(leg.StartLocation),
			ToAddress:    leg.EndAddress,
			ToLocation:   latLngToGeoPoint(leg.EndLocation),
			Fare:         groute.Fare,
		}

		wtrip := Trip{
			SourceCode:    GOOGLEAPI,
			TripType:      "foot",
			DepartureTime: models.NewTime(leg.DepartureTime),
			ArrivalTime:   models.NewTime(leg.ArrivalTime),
			FromLocation:  latLngToGeoPoint(leg.StartLocation),
			ToLocation:    latLngToGeoPoint(leg.EndLocation),
			Polyline:      groute.OverviewPolyline.Points,
			Distance:      int64(leg.Distance.Meters),
			Duration:      int64(leg.Duration.Seconds()),
		}

		if wtrip.Duration < 60 {
			wtrip.Duration = 60
		}

		mroute.Trips = []models.DataObject{{ObjectId: GOOGLE_FOOT, Data: wtrip}}

		result = append(result, mroute)
	}

	return
}
