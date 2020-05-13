package search

import (
	"fmt"
	"math"
	"time"

	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/device"
	"bitbucket.movista.ru/maas/maasapi/internal/google"
	"bitbucket.movista.ru/maas/maasapi/internal/history"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"bitbucket.movista.ru/maas/maasapi/internal/uapi/fapi"
	"bitbucket.movista.ru/maas/maasapi/internal/uapi/place"
	"github.com/google/uuid"
	"googlemaps.github.io/maps"
)

const (
	DRIVING        = "driving"
	PUBLIC         = "public"
	WALKING        = "walking"
	MOVISTA        = "movista"
	MOVISTA_BUS    = "movista-bus"
	MOVISTA_TRAIN  = "movista-train"
	MOVISTA_FLIGHT = "movista-flight"
	TAXI           = "taxi"

	WARNING = "warning"

	TRANSIT_MODE = "TRANSIT"
	WALKING_MODE = "WALKING"

	GOOGLEAPI = "google"

	MIN_WALKING_DISTANCE = 20   // meters
	MAX_WALKING_DISTANCE = 1500 // meters
	MAX_WALKING_DURATION = 15   // minutes

	GROUP_OF_TAXI          = "trip_taxi"
	GROUP_OF_DRIVING       = "trip_auto"
	GROUP_OF_MOVISTA       = "movista_trip"
	PATH_GROUP             = "path_group"
	GOOGLE_ROUTE           = "google_routes"
	GOOGLE_ROUTE_TRANSPORT = "google_route_transport"
	MOVISTA_ROUTE          = "movista_route"
	MOVISTA_PLACEHOLDER    = "movista_placeholder"
	TEXT_DATA              = "text_data"

	GOOGLE_FOOT           = "google_trip_foot"
	GOOGLE_BUS            = "google_trip_bus"
	GOOGLE_SUBWAY         = "google_trip_subway"
	GOOGLE_TRAM           = "google_trip_tram"
	GOOGLE_COMMUTER_TRAIN = "google_trip_commutertrain"
	GOOGLE_HEAVY_RAIL     = "google_trip_rail"
	GOOGLE_FERRY          = "google_trip_ferry"
	GOOGLE_TROLLEYBUS     = "google_trip_trolleybus"
	GOOGLE_SHARE_TAXI     = "google_trip_sharetaxi"
	GOOGLE_TAXI           = "google_trip_taxi"
	GOOGLE_DRIVING        = "google_route_auto"

	// generic data & icons const
	SUBWAY_TITLE         = "Метро"
	AEROEXPRESS_TITLE    = "Аэроэкспресс"
	MCD_TITLE            = "МЦД"
	BUS_TITLE            = "Автобус"
	HEAVY_RAIL_TITLE     = "Поезд"
	TRAM_TITLE           = "Трамвай"
	COMMUTER_TRAIN_TITLE = "Электричка"
	TROLLEYBUS_TITLE     = "Троллейбус"
	SHARE_TAXI_TITLE     = "Маршрутка"
	FERRY_TITLE          = "Паром"

	BUS_COLOR            = "#1CB81F"
	HEAVY_RAIL_COLOR     = "#0087f5"
	TRAM_COLOR           = "#f7331f"
	COMMUTER_TRAIN_COLOR = "#0087f5"
	TROLLEYBUS_COLOR     = "#5ebdff"
	SHARE_TAXI_COLOR     = "#c9c9c9"
	FERRY_COLOR          = "#696969"
	AEROEXPRESS_COLOR    = "#DA3732"
	MCD_D1_COLOR         = "#EBA93B"
	MCD_D2_COLOR         = "#DB4E84"
)

type FindRoutesUsecaseIn struct {
	DeviceID                                        string
	Origin, Destination                             models.GeoPoint
	OriginPlaceId, DestinationPlaceId, Region       string
	DepartureTime, ArrivalTime                      models.Time
	OriginMovistaPlaceId, DestinationMovistaPlaceId *int
	TripTypes                                       map[string]bool
	ShowPathGroups                                  bool
}

func isMovistaSearch(triptypes map[string]bool) bool {
	if triptypes[MOVISTA_BUS] || triptypes[MOVISTA_TRAIN] || triptypes[MOVISTA_FLIGHT] || triptypes[MOVISTA] {
		return true
	}

	return false
}

func isGoogleSearch(triptypes map[string]bool) bool {
	if triptypes[DRIVING] || triptypes[PUBLIC] || triptypes[WALKING] || triptypes[TAXI] {
		return true
	}

	return false
}

func getMovistaGeoPoint(ctx common.Context, movistaID int) (*models.GeoPoint, error) {
	pl, err := place.FindByID(ctx, movistaID)
	if err != nil {
		return nil, err
	}
	return &models.GeoPoint{Longitude: pl.Lon, Latitude: pl.Lat}, nil
}

func fillCoordinates(ctx common.Context, point models.GeoPoint, googleID string, movistaID *int,
	region string) (models.GeoPoint, error) {

	if point.IsZero() && googleID == "" && movistaID == nil {
		return models.GeoPoint{}, fmt.Errorf("empty location")
	}

	if point.IsZero() && movistaID != nil {
		pt, err := getMovistaGeoPoint(ctx, *movistaID)
		if err == nil {
			return *pt, nil
		}
	}

	if point.IsZero() && googleID != "" {
		loc, err := google.FindByID(ctx, uuid.Nil, googleID, region)
		if err == nil {
			return *loc.Location, nil
		}
	}

	return point, nil
}

func FindRoutes(ctx common.Context, req *FindRoutesUsecaseIn) (objects []models.DataObject, err error) {

	var (
		googleFrom, googleTo       string
		movistaFromID, movistaToID int
	)

	req.Origin, err = fillCoordinates(ctx, req.Origin, req.OriginPlaceId, req.OriginMovistaPlaceId, req.Region)
	if err != nil {
		return nil, err
	}

	req.Destination, err = fillCoordinates(ctx, req.Destination, req.DestinationPlaceId, req.DestinationMovistaPlaceId, req.Region)
	if err != nil {
		return nil, err
	}

	if req.OriginMovistaPlaceId != nil && req.DestinationMovistaPlaceId != nil {
		req.TripTypes = map[string]bool{
			MOVISTA: true,
			DRIVING: true,
		}
	}

	if req.OriginPlaceId != "" {
		googleFrom = "place_id:" + req.OriginPlaceId
	} else if !req.Origin.IsZero() {
		googleFrom = fmt.Sprint(req.Origin.Latitude, req.Origin.Longitude)
	}

	if req.DestinationPlaceId != "" {
		googleTo = "place_id:" + req.DestinationPlaceId
	} else if !req.Destination.IsZero() {
		googleTo = fmt.Sprint(req.Destination.Latitude, req.Destination.Longitude)
	}

	distance := distanceInKmBetweenEarthCoordinates(req.Origin, req.Destination)

	if len(req.TripTypes) == 0 {
		if distance < 150 {
			req.TripTypes = map[string]bool{
				PUBLIC:  true,
				DRIVING: true,
				TAXI:    true,
			}
			if distance < 3 {
				req.TripTypes[WALKING] = true
			}
		} else {
			req.TripTypes = map[string]bool{
				MOVISTA: true,
				DRIVING: true,
			}
		}
	}

	dataObjectsMap := make(map[string][]models.DataObject)

	var asyncRequestID string

	if isMovistaSearch(req.TripTypes) {
		if req.ShowPathGroups {
			departure := req.DepartureTime.Format("2006-01-02")
			result, uid, err := fapi.SearchSyncV5(ctx, req.Origin, req.Destination, departure, req.TripTypes)
			if err != nil {
				return nil, err
			}

			asyncRequestID = uid

			if len(result.PathGroups) > 0 {
				for _, pathGroup := range result.PathGroups {
					places, err := place.FindByIds(ctx, []int{
						result.SearchParams.From,
						result.SearchParams.To,
					})
					if err != nil {
						return nil, err
					}
					var fromPlace, toPlace place.Place
					var ok bool
					if fromPlace, ok = places[result.SearchParams.From]; !ok {
						return nil, fmt.Errorf("cannot find place by id: %d", result.SearchParams.From)
					}
					if toPlace, ok = places[result.SearchParams.To]; !ok {
						return nil, fmt.Errorf("cannot find place by id: %d", result.SearchParams.To)
					}
					result.SearchParams.FromPlace = &fromPlace
					result.SearchParams.ToPlace = &toPlace
					pathGroup.SearchParams = result.SearchParams
					pathGroup.SearchUID = uid
					pathGroup.CreatedAt = models.NewTime(time.Now())
					dataObjectsMap[MOVISTA] = append(dataObjectsMap[MOVISTA], models.DataObject{
						ObjectId: PATH_GROUP,
						Data:     pathGroup,
					})
				}
			}
		} else {
			movistaRoutes := findMovistaFastestLowestNearestRoutes(ctx, req.DepartureTime, req.Origin, req.Destination, req.TripTypes)
			if movistaRoutes.Count > 0 {
				dataObjectsMap[MOVISTA] = convertMovistaFastestLowestNearestRoutesToDataObjects(movistaRoutes)
			}
			movistaFromID, movistaToID = movistaRoutes.OriginID, movistaRoutes.DestinationID
		}
	}

	if req.TripTypes[PUBLIC] {
		publicRoutes, resultIsWalkingRoute := findPublicRoutes(ctx, googleFrom, googleTo, req.DepartureTime, req.ArrivalTime, req.Region)

		if resultIsWalkingRoute {
			req.TripTypes[WALKING] = false
			dataObjectsMap[PUBLIC] = []models.DataObject{
				{ObjectId: TEXT_DATA, Data: models.MessageObject{Text: "Маршруты пешком"}},
				convertPublicRoutesToDataObject(publicRoutes),
			}
		} else {
			dataObjectsMap[PUBLIC] = []models.DataObject{
				{ObjectId: TEXT_DATA, Data: models.MessageObject{Text: "Маршруты общественного транспорта"}},
				convertPublicRoutesToDataObject(publicRoutes),
			}
		}
	}

	if req.TripTypes[WALKING] {
		walkingRoutes := findWalkingRoutes(ctx, googleFrom, googleTo, req.DepartureTime, req.ArrivalTime)
		if len(walkingRoutes) > 0 && walkingRoutes[0].Duration < 2400 {
			dataObjectsMap[WALKING] = []models.DataObject{
				{ObjectId: TEXT_DATA, Data: models.MessageObject{Text: "Маршруты пешком"}},
				convertWalkingRoutesToDataObject(walkingRoutes),
			}
		}
	}

	if req.TripTypes[DRIVING] {
		drivingRoutes := findDrivingRoutes(ctx, googleFrom, googleTo, req.Origin, req.Destination, req.DepartureTime, req.ArrivalTime, req.Region)
		if len(drivingRoutes) > 0 {
			dataObjectsMap[DRIVING] = []models.DataObject{
				{ObjectId: TEXT_DATA, Data: models.MessageObject{Text: "Маршруты на машине"}},
				convertDrivingRoutesToDataObject(drivingRoutes),
			}
		}
	}

	if req.TripTypes[TAXI] {
		taxiTrips, deeplinks := findTaxiRoutes(ctx, req.Origin, req.Destination, req.DepartureTime, req.ArrivalTime, req.Region)
		if taxiTrips != nil && len(taxiTrips) > 0 {
			dataObjectsMap[TAXI] = []models.DataObject{
				{ObjectId: TEXT_DATA, Data: models.MessageObject{Text: "Маршруты на такси"}},
				convertTaxiRoutesToDataObject(taxiTrips),
			}
			minWarningDuration := 25 * time.Minute
			greaterThanArrivalTime := req.ArrivalTime.Sub(time.Now()) > minWarningDuration
			greaterThanDepartureTime := req.DepartureTime.Sub(time.Now()) > minWarningDuration
			if greaterThanDepartureTime || greaterThanArrivalTime {
				dataObjectsMap[WARNING] = []models.DataObject{
					{ObjectId: TEXT_DATA, Data: models.MessageObject{Text: "Реальная стоимость поездки может отличаться в зависимости от дорожной обстановки и спроса на такси"}},
				}
			}
			if deeplinks != nil && len(deeplinks) > 0 {
				device.SaveTaxiOrder(ctx, deeplinks)
			}
		}
	}

	objects = append(objects, dataObjectsMap[MOVISTA]...)
	objects = append(objects, dataObjectsMap[WALKING]...)
	objects = append(objects, dataObjectsMap[PUBLIC]...)
	objects = append(objects, dataObjectsMap[TAXI]...)
	objects = append(objects, dataObjectsMap[WARNING]...)

	if len(dataObjectsMap[DRIVING]) > 0 {
		objects = append(objects, dataObjectsMap[DRIVING]...)
	}

	if len(objects) == 0 {
		objects = append(objects, models.DataObject{ObjectId: TEXT_DATA, Data: models.MessageObject{Text: "Маршруты не найдены"}})
	}

	tripTypesSlice := make([]string, 0)
	for k, v := range req.TripTypes {
		if v {
			tripTypesSlice = append(tripTypesSlice, k)
		}
	}

	if isMovistaSearch(req.TripTypes) {
		searchData := &history.MovistaSearch{
			DeviceID:          ctx.DeviceID,
			FromID:            movistaFromID,
			ToID:              movistaToID,
			Origin:            &req.Origin,
			Destination:       &req.Destination,
			FromGooglePlaceID: req.OriginPlaceId,
			ToGooglePlaceID:   req.DestinationPlaceId,
			DepartureTime:     req.DepartureTime,
			ArrivalTime:       req.ArrivalTime,
			TripTypes:         tripTypesSlice,
		}
		if req.ShowPathGroups {
			searchData.ID = asyncRequestID
		}
		err = history.SaveMovistaSearch(ctx, searchData)
	}
	if isGoogleSearch(req.TripTypes) {
		err = history.SaveGoogleSearch(ctx, &history.GoogleSearch{
			DeviceID:           ctx.DeviceID,
			Origin:             &req.Origin,
			Destination:        &req.Destination,
			OriginPlaceID:      req.OriginPlaceId,
			DestinationPlaceID: req.DestinationPlaceId,
			DepartureTime:      req.DepartureTime,
			ArrivalTime:        req.ArrivalTime,
			TripTypes:          tripTypesSlice,
		})
	}

	return

}

func getRouteTimes(leg *maps.Leg, departureTime, arrivalTime models.Time) (startTime, endTime models.Time) {

	if !leg.DepartureTime.IsZero() && !leg.ArrivalTime.IsZero() {
		startTime = models.Time{Time: leg.DepartureTime}
		endTime = models.Time{Time: leg.ArrivalTime}
	} else {
		if departureTime.After(time.Now()) {
			startTime = departureTime
			endTime = startTime.Add(leg.Duration)
		} else if arrivalTime.After(time.Now()) {
			endTime = arrivalTime
			startTime = endTime.Add(-leg.Duration)
		} else {
			startTime = models.Time{Time: time.Now()}
			endTime = startTime.Add(leg.Duration)
		}
	}

	return
}

func latLngToGeoPoint(latLng maps.LatLng) models.GeoPoint {
	return models.GeoPoint{Latitude: latLng.Lat, Longitude: latLng.Lng}
}

func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func distanceInKmBetweenEarthCoordinates(coord1, coord2 models.GeoPoint) int {
	earthRadiusKm := 6371.0
	dLat := degreesToRadians(coord2.Latitude - coord1.Latitude)
	dLon := degreesToRadians(coord2.Longitude - coord1.Longitude)
	lat1 := degreesToRadians(coord1.Latitude)
	lat2 := degreesToRadians(coord2.Latitude)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return int(earthRadiusKm * c)
}
