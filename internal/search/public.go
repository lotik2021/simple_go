package search

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"bitbucket.movista.ru/maas/maasapi/internal/device"
	"bitbucket.movista.ru/maas/maasapi/internal/google"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"bitbucket.movista.ru/maas/maasapi/internal/transport/metro"
	"fmt"
	"sort"
	"strings"
	"time"

	"bitbucket.movista.ru/maas/maasapi/internal/logger"
	"bitbucket.movista.ru/maas/maasapi/internal/transport/taxi"
	"github.com/google/uuid"
	"googlemaps.github.io/maps"
)

func convertPublicRoutesToDataObject(mroutes []MRoute) (dataObject models.DataObject) {

	sort.SliceStable(mroutes, func(i, j int) bool {
		mi, mj := mroutes[i], mroutes[j]
		switch {
		case mi.Duration != mj.Duration:
			return mi.Duration < mj.Duration
		default:
			return mi.StartTime.Before(mj.StartTime.Time)
		}
	})

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

func findPublicRoutes(ctx common.Context, from, to string, departureTime, arrivalTime models.Time, region string) (result []MRoute, onlyWalking bool) {
	modes := device.GetGoogleTransitModeSettings(ctx)
	routes, err := google.FindPublicRoutes(ctx, from, to, modes, maps.TravelModeTransit, departureTime, arrivalTime)
	if err != nil {
		return
	}

	// TODO: глянуть на досуге влияет ли region == "en" на ошибки
	region = "ru"

	result = make([]MRoute, 0)

	for _, route := range routes {
		leg := route.Legs[0] //TODO ? что делать, если будет несколько точек

		routeStartTime, _ := getRouteTimes(leg, departureTime, arrivalTime)

		var (
			routeDurationInSeconds int64
			routeDistanceInMeters  int64
		)

		mroute := MRoute{
			ID:           uuid.New().String(),
			FromAddress:  leg.StartAddress,
			FromLocation: latLngToGeoPoint(leg.StartLocation),
			ToAddress:    leg.EndAddress,
			ToLocation:   latLngToGeoPoint(leg.EndLocation),
			Fare:         route.Fare,
		}

		agencies := make([]Agency, 0)

		// проверка на то, что public содержит только один step и это walking
		if len(leg.Steps) == 1 && leg.Steps[0].TravelMode == WALKING_MODE {
			onlyWalking = true
		} else {
			onlyWalking = false
		}

		stepDepartureTime := models.NewTime(leg.DepartureTime)
		for idx, step := range leg.Steps {
			var (
				trip      *Trip
				taxiTrips []*taxi.Trip
				deeplinks map[string]models.DeepLink
			)

			// заменяем пеший step на такси, если дистанция больше MAX_WALKING_DISTANCE и длительность больше MAX_WALKING_DURATION
			if step.TravelMode == WALKING_MODE && step.Distance.Meters >= MAX_WALKING_DISTANCE && step.Duration.Minutes() >= MAX_WALKING_DURATION {
				origin := latLngToGeoPoint(step.StartLocation)
				destination := latLngToGeoPoint(step.EndLocation)
				taxiTrips, deeplinks = findTaxiRoutes(ctx, origin, destination, models.Time{Time: time.Unix(0, 0)}, *stepDepartureTime, region)
			} else if step.TravelMode == TRANSIT_MODE {
				trip = getTransitTrip(step, &agencies)
			} else if step.TravelMode == WALKING_MODE {
				trip = getWalkingTrip(leg, idx, step)
			} else {
				logger.Log.Errorf("unknown step travelMode %s", step.TravelMode)
				continue
			}

			if trip != nil {
				if idx == 0 {
					routeStartTime = *trip.DepartureTime
				}

				routeDurationInSeconds = routeDurationInSeconds + trip.Duration
				routeDistanceInMeters = routeDistanceInMeters + trip.Distance

				stepDepartureTime = trip.ArrivalTime

				mroute.Trips = append(mroute.Trips, models.DataObject{ObjectId: trip.ObjectID, Data: trip})
			} else if len(taxiTrips) > 0 {
				// добавляем 3 минуты на посадку
				for _, tt := range taxiTrips {
					*tt.StartTime = tt.StartTime.Add(1 * time.Minute)
					*tt.EndTime = tt.EndTime.Add(1 * time.Minute)
				}

				if idx == 0 {
					routeStartTime = *taxiTrips[0].StartTime
				}

				routeDurationInSeconds = routeDurationInSeconds + taxiTrips[0].Duration
				routeDistanceInMeters = routeDistanceInMeters + taxiTrips[0].Distance

				*stepDepartureTime = stepDepartureTime.Add(step.Duration)

				if len(deeplinks) > 0 {
					device.SaveTaxiOrder(ctx, deeplinks)
				}

				mroute.Trips = append(mroute.Trips, convertTaxiRoutesToDataObject(taxiTrips))
			}
		}

		if len(agencies) > 0 {
			mroute.Agencies = agencies
		}

		mroute.StartTime = &routeStartTime

		endTime := routeStartTime.Add(time.Duration(routeDurationInSeconds) * time.Second)
		mroute.EndTime = &endTime
		mroute.Duration = routeDurationInSeconds / 60
		mroute.Distance = routeDistanceInMeters

		result = append(result, mroute)
	}

	return
}

func getTransitTrip(step *maps.Step, agencies *[]Agency) (mtrip *Trip) {

	var (
		objectId           string
		title              string
		color              string
		iconDetails        string
		iconShort          string
		androidIconDetails string
		androidIconShort   string
	)

	switch step.TransitDetails.Line.Vehicle.Type {
	case "BUS":
		objectId = GOOGLE_BUS
		title = BUS_TITLE
		color = BUS_COLOR
		iconDetails = config.C.Icons.GenericData.DetailsIcons.Bus.Ios
		iconShort = config.C.Icons.GenericData.ShortIcons.Bus.Ios
		androidIconDetails = config.C.Icons.GenericData.DetailsIcons.Bus.Android
		androidIconShort = config.C.Icons.GenericData.ShortIcons.Bus.Android
	case "SUBWAY":
		objectId = GOOGLE_SUBWAY
		title = SUBWAY_TITLE
		iconDetails = config.C.Icons.GenericData.DetailsIcons.Metro.Ios
		androidIconDetails = config.C.Icons.GenericData.DetailsIcons.Metro.Android
	case "TRAM":
		objectId = GOOGLE_TRAM
		title = TRAM_TITLE
		color = TRAM_COLOR
		iconDetails = config.C.Icons.GenericData.DetailsIcons.Tram.Ios
		iconShort = config.C.Icons.GenericData.ShortIcons.Tram.Ios
		androidIconDetails = config.C.Icons.GenericData.DetailsIcons.Tram.Android
		androidIconShort = config.C.Icons.GenericData.ShortIcons.Tram.Android
	case "COMMUTER_TRAIN":
		objectId = GOOGLE_COMMUTER_TRAIN
		title = COMMUTER_TRAIN_TITLE
		color = COMMUTER_TRAIN_COLOR
		iconDetails = config.C.Icons.GenericData.DetailsIcons.Train.Ios
		iconShort = config.C.Icons.GenericData.ShortIcons.CommuterTrain.Ios
		androidIconDetails = config.C.Icons.GenericData.DetailsIcons.Train.Android
		androidIconShort = config.C.Icons.GenericData.ShortIcons.CommuterTrain.Android
	case "HEAVY_RAIL":
		objectId = GOOGLE_HEAVY_RAIL
		title = HEAVY_RAIL_TITLE
		color = HEAVY_RAIL_COLOR
		iconDetails = config.C.Icons.GenericData.DetailsIcons.Train.Ios
		iconShort = config.C.Icons.GenericData.ShortIcons.Train.Ios
		androidIconDetails = config.C.Icons.GenericData.DetailsIcons.Train.Android
		androidIconShort = config.C.Icons.GenericData.ShortIcons.Train.Android
	case "FERRY":
		objectId = GOOGLE_FERRY
		title = FERRY_TITLE
		color = FERRY_COLOR
		iconDetails = config.C.Icons.GenericData.DetailsIcons.Ferry.Ios
		iconShort = config.C.Icons.GenericData.ShortIcons.Ferry.Ios
		androidIconDetails = config.C.Icons.GenericData.DetailsIcons.Ferry.Android
		androidIconShort = config.C.Icons.GenericData.ShortIcons.Ferry.Android
	case "TROLLEYBUS":
		objectId = GOOGLE_TROLLEYBUS
		title = TROLLEYBUS_TITLE
		color = TROLLEYBUS_COLOR
		iconDetails = config.C.Icons.GenericData.DetailsIcons.Trolleybus.Ios
		iconShort = config.C.Icons.GenericData.ShortIcons.Trolleybus.Ios
		androidIconDetails = config.C.Icons.GenericData.DetailsIcons.Trolleybus.Android
		androidIconShort = config.C.Icons.GenericData.ShortIcons.Trolleybus.Android
	case "SHARE_TAXI":
		objectId = GOOGLE_SHARE_TAXI
		title = SHARE_TAXI_TITLE
		color = SHARE_TAXI_COLOR
		iconDetails = config.C.Icons.GenericData.DetailsIcons.Bus.Ios
		iconShort = config.C.Icons.GenericData.ShortIcons.ShareTaxi.Ios
		androidIconDetails = config.C.Icons.GenericData.DetailsIcons.Bus.Android
		androidIconShort = config.C.Icons.GenericData.ShortIcons.ShareTaxi.Android
	}

	td := step.TransitDetails
	lineColor, shortName, iosIconUrl, androidIconUrl := metro.GetLineColor(td.Line.Name, td.Line.ShortName, td.Line.Agencies[0].URL.String())
	if objectId == GOOGLE_SUBWAY {
		if iosIconUrl != "" {
			iconShort = iosIconUrl
			androidIconShort = androidIconUrl
			color = lineColor
		} else {
			lineColor = step.TransitDetails.Line.Color
			color = lineColor
			iconShort = config.C.Icons.GenericData.ShortIcons.MetroDefault.Ios
			androidIconShort = config.C.Icons.GenericData.ShortIcons.MetroDefault.Android
		}
	}
	if td.Line.ShortName == "D1" || td.Line.ShortName == "D2" {
		switch {
		case td.Line.ShortName == "D1":
			title = MCD_TITLE
			color = MCD_D1_COLOR
			iconDetails = config.C.Icons.GenericData.DetailsIcons.Mcd.Ios
			iconShort = config.C.Icons.GenericData.ShortIcons.McdD1.Ios
			androidIconDetails = config.C.Icons.GenericData.DetailsIcons.Mcd.Android
			androidIconShort = config.C.Icons.GenericData.ShortIcons.McdD1.Android

		case td.Line.ShortName == "D2":
			title = MCD_TITLE
			color = MCD_D2_COLOR
			iconDetails = config.C.Icons.GenericData.DetailsIcons.Mcd.Ios
			iconShort = config.C.Icons.GenericData.ShortIcons.McdD2.Ios
			androidIconDetails = config.C.Icons.GenericData.DetailsIcons.Mcd.Android
			androidIconShort = config.C.Icons.GenericData.ShortIcons.McdD2.Android
		}
	}
	for _, aero := range td.Line.Agencies {
		if aero.Name == AEROEXPRESS_TITLE {
			title = AEROEXPRESS_TITLE
			color = AEROEXPRESS_COLOR
			iconDetails = config.C.Icons.GenericData.DetailsIcons.Aeroexpress.Ios
			iconShort = config.C.Icons.GenericData.ShortIcons.Aeroexpress.Ios
			androidIconDetails = config.C.Icons.GenericData.DetailsIcons.Aeroexpress.Android
			androidIconShort = config.C.Icons.GenericData.ShortIcons.Aeroexpress.Android
		}
	}

	gd := &GenericDataInfo{
		Color:              color,
		Title:              title,
		ShortDescription:   td.Line.ShortName,
		FullDescription:    td.Line.Name,
		IconDetails:        iconDetails,
		IconShort:          iconShort,
		AndroidIconDetails: androidIconDetails,
		AndroidIconShort:   androidIconShort,
	}

	if objectId == GOOGLE_BUS || objectId == GOOGLE_TRAM || objectId == GOOGLE_SHARE_TAXI || objectId == GOOGLE_TROLLEYBUS {
		gd.FullDescription = ""
	}

	if title == MCD_TITLE {
		gd.FullDescription = strings.Replace(gd.FullDescription, `МЦД-1 `, "", 1)
		gd.FullDescription = strings.Replace(gd.FullDescription, `МЦД-2 `, "", 1)
		gd.FullDescription = strings.Replace(gd.FullDescription, `"`, "", 2)
	}

	if objectId == GOOGLE_COMMUTER_TRAIN && title != MCD_TITLE && title != AEROEXPRESS_TITLE {
		gd.ShortDescription = td.TripShortName
	}

	mtrip = &Trip{
		ObjectID:      objectId,
		TripType:      strings.ToLower(step.TransitDetails.Line.Vehicle.Type),
		Duration:      int64(step.Duration.Seconds()),
		Polyline:      step.Polyline.Points,
		Distance:      int64(step.Distance.Meters),
		FromLocation:  latLngToGeoPoint(td.DepartureStop.Location),
		DepartureStop: &td.DepartureStop.Name,
		DepartureTime: models.NewTime(td.DepartureTime),
		ToLocation:    latLngToGeoPoint(td.ArrivalStop.Location),
		ArrivalStop:   &td.ArrivalStop.Name,
		ArrivalTime:   models.NewTime(td.ArrivalTime),
		LineNumber:    &td.Line.ShortName,
		LineName:      &td.Line.Name,
		LineColor:     &lineColor,
		GenericData:   gd,
		Headsign:      &td.Headsign,
		NumberStops:   &td.NumStops,
		TripShortName: &td.TripShortName,
	}

	if objectId == GOOGLE_SUBWAY {
		mtrip.LineName = &shortName
	}

	toLoc := step.TransitDetails.ArrivalStop.Location
	mtrip.ToLocation = latLngToGeoPoint(toLoc)

	fromLoc := step.TransitDetails.DepartureStop.Location
	mtrip.FromLocation = latLngToGeoPoint(fromLoc)

	for _, x := range td.Line.Agencies {
		*agencies = append(*agencies, Agency{Name: x.Name, Phone: strings.Replace(x.Phone, "011 ", "+", 1), Url: x.URL.String()})
	}

	return
}

func getWalkingTrip(leg *maps.Leg, idx int, step *maps.Step) (mtrip *Trip) {

	mtrip = &Trip{
		ObjectID:   GOOGLE_FOOT,
		SourceCode: GOOGLEAPI,
		TripType:   "foot",
		Polyline:   step.Polyline.Points,
		Distance:   int64(step.Distance.Meters),
		Duration:   int64(step.Duration.Seconds()),
	}

	if mtrip.Duration < 60 {
		mtrip.Duration = 60
	}

	if idx == 0 {
		if leg.DepartureTime.IsZero() {
			now := time.Now()
			mtrip.DepartureTime = models.NewTime(now)
		} else {
			mtrip.DepartureTime = models.NewTime(leg.DepartureTime)
		}

		mtrip.FromLocation = latLngToGeoPoint(step.StartLocation)

		if len(leg.Steps) == 1 {
			if leg.ArrivalTime.IsZero() {
				d, _ := time.ParseDuration(fmt.Sprintf("%ds", int(mtrip.Duration)))
				mtrip.ArrivalTime = models.NewTime(time.Now().Add(d))
			} else {
				mtrip.ArrivalTime = models.NewTime(leg.ArrivalTime)
			}
			mtrip.ToLocation = latLngToGeoPoint(step.EndLocation)
		} else {
			tdNext := leg.Steps[1].TransitDetails
			if tdNext != nil {
				mtrip.ArrivalTime = models.NewTime(tdNext.DepartureTime)
				mtrip.ToLocation = latLngToGeoPoint(tdNext.DepartureStop.Location)
			}
		}
	} else if idx == len(leg.Steps)-1 {
		tdPrev := leg.Steps[idx-1].TransitDetails
		if tdPrev != nil {
			mtrip.ArrivalTime = models.NewTime(tdPrev.DepartureTime)
			mtrip.FromLocation = latLngToGeoPoint(tdPrev.ArrivalStop.Location)
		}

		mtrip.DepartureTime = models.NewTime(leg.ArrivalTime)
		mtrip.ToLocation = latLngToGeoPoint(step.EndLocation)
	} else {
		tdPrev := leg.Steps[idx-1].TransitDetails
		if tdPrev != nil {
			mtrip.DepartureTime = models.NewTime(tdPrev.ArrivalTime)
			mtrip.FromLocation = latLngToGeoPoint(tdPrev.ArrivalStop.Location)
		}

		tdNext := leg.Steps[idx+1].TransitDetails
		if tdNext != nil {
			mtrip.ArrivalTime = models.NewTime(tdNext.DepartureTime)
			mtrip.ToLocation = latLngToGeoPoint(tdNext.DepartureStop.Location)
		}
	}

	return
}
