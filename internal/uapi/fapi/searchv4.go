package fapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"time"

	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"bitbucket.movista.ru/maas/maasapi/internal/ratelimit"
	"bitbucket.movista.ru/maas/maasapi/internal/uapi/place"
	"github.com/labstack/echo/v4"
)

const (
	typedTripTrainObjectType  = "tripTrain"
	typedTripFlightObjectType = "tripFlight"
	typedTripBusObjectType    = "tripBus"

	TripTrain  = "train"
	TripFlight = "flight"
	TripBus    = "bus"

	MOVISTA_BUS    = "movista-bus"
	MOVISTA_TRAIN  = "movista-train"
	MOVISTA_FLIGHT = "movista-flight"
)

var (
	iosIcons     = make(map[string]string)
	androidIcons = make(map[string]string)
)

func init() {
	iosIcons = map[string]string{
		TripBus:    config.C.Icons.Trips.Bus.Ios,
		TripFlight: config.C.Icons.Trips.Avia.Ios,
		TripTrain:  config.C.Icons.Trips.Train.Ios,
	}
	androidIcons = map[string]string{
		TripBus:    config.C.Icons.Trips.Bus.Android,
		TripFlight: config.C.Icons.Trips.Avia.Android,
		TripTrain:  config.C.Icons.Trips.Train.Android,
	}
}

type searchResponseV4 struct {
	RoutesInfo struct {
		Trips map[string]struct {
			// typed searchModel
			ObjectType string `json:"objectType"`
			TripTrain  *trip  `json:"tripTrain"`
			TripBus    *trip  `json:"tripBus"`
			TripFlight *trip  `json:"tripFlight"`
			// untyped searchModel
			TripType      string `json:"tripType"`
			DepartureTime string `json:"departure"`
			ArrivalTime   string `json:"arrival"`
		} `json:"trips"`
		Routes []struct {
			ID       string   `json:"id"`
			TripIDs  []string `json:"tripIds"`
			Price    float64  `json:"minAgrPrice"`
			Duration int      `json:"routeDuration"`
		} `json:"routes"`
	} `json:"routesInfo"`
	SearchParams struct {
		DepartureTime string `json:"departureBegin"`
		From          int    `json:"from"`
		To            int    `json:"to"`
	} `json:"searchParams"`
	CountResult int `json:"countResult"`
}

type trip struct {
	DepartureTime string `json:"departure"`
	ArrivalTime   string `json:"arrival"`
	TripType      string `json:"tripType"`
}

func unmapTripTypes(tripTypes map[string]bool) []string {
	longTypes := make([]string, 0)
	if tripTypes[MOVISTA_BUS] {
		longTypes = append(longTypes, TripBus)
	}
	if tripTypes[MOVISTA_FLIGHT] {
		longTypes = append(longTypes, TripFlight)
	}
	if tripTypes[MOVISTA_TRAIN] {
		longTypes = append(longTypes, TripTrain)
	}

	if len(longTypes) == 0 {
		longTypes = []string{"bus", "flight", "train"}
	}

	return longTypes

}

func findRoutesV4(ctx common.Context, departureTime models.Time, from, to int, tripTypes map[string]bool) (*models.RawResponse, error) {

	var (
		req = struct {
			CurrencyCode    string           `json:"currencyCode"`
			CultureCode     string           `json:"cultureCode"`
			DepartureBegin  string           `json:"departureBegin"`
			From            int              `json:"from"`
			To              int              `json:"to"`
			Customers       []map[string]int `json:"customers"`
			TripTypes       []string         `json:"tripTypes"`
			AirServiceClass string           `json:"airServiceClass"`
		}{
			CurrencyCode:    "RUB",
			CultureCode:     "RU",
			DepartureBegin:  departureTime.Format("2006-01-02"),
			From:            from,
			To:              to,
			Customers:       []map[string]int{{"id": 0}},
			TripTypes:       unmapTripTypes(tripTypes),
			AirServiceClass: "economy",
		}
	)

	rawResp, err := common.UapiPost(ctx, faClient, map[string]interface{}{"searchParams": req},
		config.C.FapiAdapter.Urls.SearchV4)
	if err != nil {
		return nil, err
	}

	return rawResp, nil
}

type ShortMovistaRoute struct {
	ID                string   `json:"-"`
	Price             float64  `json:"price"`
	Trips             []string `json:"trips"`
	IosIconUrl        []string `json:"trip_icons_ios"`
	AndroidIconUrl    []string `json:"trip_icons_android"`
	Duration          int      `json:"duration"`
	DepartureTime     string   `json:"departure_time"`
	ArrivalTime       string   `json:"arrival_time"`
	RedirectUrl       string   `json:"redirect_url"`
	NumberOfTransfers int      `json:"number_of_transfers"`
}

type FastestLowestNearestRes struct {
	Routes                      map[string]ShortMovistaRoute
	Count                       int
	OriginName, DestinationName string
	OriginID, DestinationID     int
	DepartureTime               models.Time
}

func FindFastestLowestNearestRoutesV4(ctx common.Context, departureTime models.Time, from, to models.GeoPoint, tripTypes map[string]bool) (res *FastestLowestNearestRes, err error) {

	res = &FastestLowestNearestRes{
		Routes:        make(map[string]ShortMovistaRoute),
		DepartureTime: departureTime,
	}

	placeIdFrom, originName, err := place.FindOneByLocation(ctx, &from)
	if err != nil {
		return
	}

	res.OriginID = placeIdFrom
	res.OriginName = originName

	placeIdTo, destinationName, err := place.FindOneByLocation(ctx, &to)
	if err != nil {
		return
	}

	res.DestinationID = placeIdTo
	res.DestinationName = destinationName

	rawResp, err := findRoutesV4(ctx, departureTime, placeIdFrom, placeIdTo, tripTypes)
	if err != nil {
		err = fmt.Errorf("error from fapiadapter v4 search %w", err)
		return
	}

	var resp searchResponseV4
	err = json.Unmarshal(rawResp.Data, &resp)
	if err != nil {
		err = fmt.Errorf("cannot unmarshal fapiadapter response - %w", err)
		return
	}

	// FIXME? не надо создавать ошибку
	if len(rawResp.Error) > 0 {
		err = fmt.Errorf("error in fapiadapter response %+v", rawResp.Error)
	}

	res.Count = resp.CountResult

	if res.Count == 0 {
		return
	}
	//https://movista.ru/tickets/train-flight-bus/Москва/Париж/65537/70577?customers%5B0%5D%5Bage%5D=36&customers%5B0%5D%5BseatRequired%5D=true&departure=2019-12-28&returnDeparture=2020-01-18&serviceClass=Economy
	newRedirectUrl := fmt.Sprintf("%s/tickets/train-flight-bus/%s/%s/%d/%d?customers[0][age]=36&customers[0][seatRequired]=true&departure=%s&serviceClass=Economy",
		config.C.Web.BaseURL,
		url.QueryEscape(res.OriginName),
		url.QueryEscape(res.DestinationName),
		resp.SearchParams.From,
		resp.SearchParams.To,
		departureTime.Format("2006-01-02"),
	)

	//convert to untyped result for future tests
	for k, v := range resp.RoutesInfo.Trips {
		if v.ObjectType == "" {
			continue
		}

		var tripType, departureTime, arrivalTime string

		switch v.ObjectType {
		case typedTripFlightObjectType:
			tripType = TripFlight
			departureTime = v.TripFlight.DepartureTime
			arrivalTime = v.TripFlight.ArrivalTime
		case typedTripBusObjectType:
			tripType = TripBus
			departureTime = v.TripBus.DepartureTime
			arrivalTime = v.TripBus.ArrivalTime
		case typedTripTrainObjectType:
			tripType = TripTrain
			departureTime = v.TripTrain.DepartureTime
			arrivalTime = v.TripTrain.ArrivalTime
		}

		v.TripType = tripType
		v.DepartureTime = departureTime
		v.ArrivalTime = arrivalTime

		resp.RoutesInfo.Trips[k] = v
	}

	// fastest
	sort.Slice(resp.RoutesInfo.Routes, func(i, j int) bool {
		return resp.RoutesInfo.Routes[i].Duration < resp.RoutesInfo.Routes[j].Duration
	})

	fTripsInfo := obtainTripsInfo(resp.RoutesInfo.Routes[0].TripIDs, &resp)

	res.Routes["fastest"] = ShortMovistaRoute{
		ID:                resp.RoutesInfo.Routes[0].ID,
		Price:             resp.RoutesInfo.Routes[0].Price,
		Trips:             fTripsInfo.TripTypes,
		IosIconUrl:        fTripsInfo.iosIcons,
		AndroidIconUrl:    fTripsInfo.androidIcons,
		Duration:          resp.RoutesInfo.Routes[0].Duration,
		DepartureTime:     fTripsInfo.DepartureTime,
		ArrivalTime:       fTripsInfo.ArrivalTime,
		RedirectUrl:       newRedirectUrl,
		NumberOfTransfers: fTripsInfo.NumberOfTransfers,
	}
	// lowest (cheapest)
	sort.Slice(resp.RoutesInfo.Routes, func(i, j int) bool {
		return resp.RoutesInfo.Routes[i].Price < resp.RoutesInfo.Routes[j].Price
	})

	lTripsInfo := obtainTripsInfo(resp.RoutesInfo.Routes[0].TripIDs, &resp)
	res.Routes["lowest"] = ShortMovistaRoute{
		ID:                resp.RoutesInfo.Routes[0].ID,
		Price:             resp.RoutesInfo.Routes[0].Price,
		Trips:             lTripsInfo.TripTypes,
		IosIconUrl:        lTripsInfo.iosIcons,
		AndroidIconUrl:    lTripsInfo.androidIcons,
		Duration:          resp.RoutesInfo.Routes[0].Duration,
		DepartureTime:     lTripsInfo.DepartureTime,
		ArrivalTime:       lTripsInfo.ArrivalTime,
		RedirectUrl:       newRedirectUrl,
		NumberOfTransfers: lTripsInfo.NumberOfTransfers,
	}

	res.Routes["nearest"] = findNearest(&resp, newRedirectUrl)

	return
}

type obtainTripsInfoResp struct {
	TripTypes                  []string
	iosIcons                   []string
	androidIcons               []string
	DepartureTime, ArrivalTime string
	NumberOfTransfers          int
}

func obtainTripsInfo(routeTripIds []string, resp *searchResponseV4) (info obtainTripsInfoResp) {

	info = obtainTripsInfoResp{
		iosIcons:     make([]string, 0),
		androidIcons: make([]string, 0),
	}

	correctTripIds := make([]string, 0)
	for _, v := range routeTripIds {
		ot := resp.RoutesInfo.Trips[v].TripType
		if ot == TripBus || ot == TripFlight || ot == TripTrain {
			correctTripIds = append(correctTripIds, v)
		}
	}

	for i, v := range correctTripIds {
		trip := resp.RoutesInfo.Trips[v]
		info.TripTypes = append(info.TripTypes, trip.TripType)
		info.iosIcons = append(info.iosIcons, iosIcons[trip.TripType])
		info.androidIcons = append(info.androidIcons, androidIcons[trip.TripType])
		info.NumberOfTransfers++
		if i == 0 {
			info.DepartureTime = trip.DepartureTime
		}
		if i == len(correctTripIds)-1 {
			info.ArrivalTime = trip.ArrivalTime
		}
	}

	return
}

func findNearest(resp *searchResponseV4, redirectUrl string) (route ShortMovistaRoute) {

	route = ShortMovistaRoute{
		RedirectUrl: redirectUrl,
	}

	var (
		now = time.Now().Format(time.RFC3339)

		nearestTripID        string
		nearestTripDeparture string

		nearestRouteTripIDs []string
	)

	for i, v := range resp.RoutesInfo.Trips {
		if v.TripType != TripBus && v.TripType != TripTrain && v.TripType != TripFlight {
			continue
		}

		tmpDepartureTime := v.DepartureTime

		if tmpDepartureTime <= now {
			continue
		}

		if nearestTripID == "" {
			nearestTripDeparture = tmpDepartureTime
			nearestTripID = i

			continue
		}

		if tmpDepartureTime < nearestTripDeparture {
			nearestTripID = i
			nearestTripDeparture = tmpDepartureTime
		}
	}

routesLoop:
	for _, v := range resp.RoutesInfo.Routes {
		for _, tid := range v.TripIDs {
			if tid == nearestTripID {
				route.ID = v.ID
				route.Price = v.Price
				route.Duration = v.Duration
				route.DepartureTime = nearestTripDeparture
				nearestRouteTripIDs = v.TripIDs
				break routesLoop
			}
		}
	}

	nTripsInfo := obtainTripsInfo(nearestRouteTripIDs, resp)

	route.NumberOfTransfers = nTripsInfo.NumberOfTransfers
	route.Trips = nTripsInfo.TripTypes
	route.IosIconUrl = nTripsInfo.iosIcons
	route.AndroidIconUrl = nTripsInfo.androidIcons
	route.ArrivalTime = nTripsInfo.ArrivalTime

	return
}

func SaveSelectedRoutesV4(c echo.Context) error {
	ctx := common.NewContext(c)

	err := ratelimit.Apply(ctx, ratelimit.MethodSaveSelectedRoutes)
	if err != nil {
		return err
	}

	uid := c.Param("uid")
	forwardRouteId := c.Param("forwardRouteId")
	url := config.C.FapiAdapter.Urls.SaveSelectedRoutesV4 + fmt.Sprintf("/%s/%s", uid, forwardRouteId)

	rawResp, err := common.UapiGet(ctx, faClient, url)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, rawResp)
}

func GetSelectedRoutesV4(c echo.Context) error {
	var (
		ctx = common.NewContext(c)
	)

	uid := c.Param("uid")
	if uid == "" {
		return fmt.Errorf("uid is empty")
	}

	url := config.C.FapiAdapter.Urls.GetSelectedRoutesV4 + fmt.Sprintf("/%s", uid)
	rawResp, err := common.UapiGet(ctx, faClient, url)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, rawResp)
}
