package taxi

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/google"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"fmt"
	"sort"
	"sync"
	"time"

	"bitbucket.movista.ru/maas/maasapi/internal/logger"
	"googlemaps.github.io/maps"
)

var (
	cityMobileToken string
	TAXI_ROUTE      = "taxi_route"
)

func init() {
	err := authCityMobile()
	if err != nil {
		logger.Log.WithError(err).Error("cannot authorize in citymobile")
	}
}

func CalculateRoute(ctx common.Context, origin, destination models.GeoPoint, departureTime, arrivalTime models.Time, region string) (taxiTrips []*Trip, deeplinks map[string]models.DeepLink) {
	var (
		drivingRoute                      maps.Route
		cityTrip, maximTrip, yandexTrip   Trip
		originAddress, destinationAddress string
	)

	if departureTime.Before(time.Now()) {
		departureTime = *models.NewTime(time.Now())
	}

	taxiTrips = make([]*Trip, 0)
	deeplinks = make(map[string]models.DeepLink)

	wg := &sync.WaitGroup{}
	wg.Add(6)

	go func(addr *string) {
		defer wg.Done()
		place, err := google.FindByLocation(ctx, &origin, region)
		if err != nil {
			return
		}

		*addr = place.MainText

	}(&originAddress)

	go func(addr *string) {
		defer wg.Done()

		place, err := google.FindByLocation(ctx, &destination, region)
		if err != nil {
			return
		}

		*addr = place.MainText

	}(&destinationAddress)

	go func(route *maps.Route) {
		defer wg.Done()
		from := fmt.Sprint(origin.Latitude, origin.Longitude)
		to := fmt.Sprint(destination.Latitude, destination.Longitude)
		routes, err := google.FindDrivingRoutes(ctx, from, to, maps.TravelModeDriving, maps.TrafficModelBestGuess, false, departureTime, *models.NewTime(time.Unix(0, 0)))
		if err != nil || len(routes) == 0 {
			*route = maps.Route{}
			return
		}

		*route = routes[0]

	}(&drivingRoute)

	go func(trip *Trip) {
		defer wg.Done()
		if citymobileTrip, err := CalculateRouteCityMobile(ctx, origin, destination); err != nil {
			logger.Log.WithError(err).Error("citymobile error")
		} else if citymobileTrip != nil {
			*trip = *citymobileTrip
		}
	}(&cityTrip)

	go func(trip *Trip) {
		defer wg.Done()
		if maximTrip, err := CalculateRouteMaxim(ctx, origin, destination); err != nil {
			logger.Log.WithError(err).Error("maxim error")
		} else if maximTrip != nil {
			*trip = *maximTrip
		}
	}(&maximTrip)

	go func(trip *Trip) {
		defer wg.Done()
		if yandexTrip, err := CalculateRouteYandex(ctx, origin, destination); err != nil {
			logger.Log.WithError(err).Error("yandex error")
		} else if yandexTrip != nil {
			*trip = *yandexTrip
		}
	}(&yandexTrip)

	wg.Wait()

	firstLeg := drivingRoute.Legs[0]
	startTime := departureTime
	duration := int64(firstLeg.DurationInTraffic.Seconds())
	endTime := startTime.Add(time.Duration(duration) * time.Second)
	distance := int64(firstLeg.Distance.Meters)
	polyline := drivingRoute.OverviewPolyline.Points

	if time.Now().Before(arrivalTime.Time) {
		startTime = arrivalTime.Add(-firstLeg.DurationInTraffic)
		endTime = arrivalTime
	}

	if cityTrip.ID != "" {
		taxiTrips = append(taxiTrips, &cityTrip)
	}
	if maximTrip.ID != "" {
		maximTrip.DeepLink = makeMaximDeepLink(maximTrip.ID, origin, destination, originAddress, destinationAddress)
		taxiTrips = append(taxiTrips, &maximTrip)
	}
	if yandexTrip.ID != "" {
		taxiTrips = append(taxiTrips, &yandexTrip)
	}

	if len(taxiTrips) == 0 {
		return
	}

	for _, taxiTrip := range taxiTrips {
		taxiTrip.FromAddress = originAddress
		taxiTrip.ToAddress = destinationAddress
		taxiTrip.Duration = duration
		taxiTrip.Distance = distance
		taxiTrip.StartTime = &startTime
		taxiTrip.EndTime = &endTime
		taxiTrip.Polyline = polyline

		deeplinks[taxiTrip.ProviderDescription] = taxiTrip.DeepLink
	}

	sort.Slice(taxiTrips, func(i, j int) bool {
		return taxiTrips[i].MinPrice < taxiTrips[j].MinPrice
	})

	return
}
