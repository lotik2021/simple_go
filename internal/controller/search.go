package controller

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"bitbucket.movista.ru/maas/maasapi/internal/device"
	"bitbucket.movista.ru/maas/maasapi/internal/favorite"
	"bitbucket.movista.ru/maas/maasapi/internal/google"
	"bitbucket.movista.ru/maas/maasapi/internal/history"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"bitbucket.movista.ru/maas/maasapi/internal/search"
	"bitbucket.movista.ru/maas/maasapi/internal/uapi/place"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func searchFindAutocompletePlaces(c echo.Context) (err error) {
	var (
		req struct {
			SessionToken string `json:"token"`
			Region       string `json:"region" validate:"required"`
			Input        string `json:"input"`
		}
		ctx = common.NewContext(c)
	)

	err = common.BindAndValidateReq(c, &req)
	if err != nil {
		return
	}

	loc := device.GetLastLocation(ctx)

	sessionUUID, _ := uuid.Parse(req.SessionToken)
	result, err := google.FindAutocomplete(ctx, &sessionUUID, req.Input, req.Region, loc)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"sessionToken": sessionUUID.String(),
		"result":       result,
	})
}

type FindRoutesRequest struct {
	Origin      models.GeoPoint `json:"origin" validate:"required"`
	Destination models.GeoPoint `json:"destination" validate:"required"`

	OriginPlaceId             string   `json:"origin_google_place_id"`
	DestinationPlaceId        string   `json:"destination_google_place_id"`
	OriginMovistaPlaceId      *int     `json:"origin_movista_place_id"`
	DestinationMovistaPlaceId *int     `json:"destination_movista_place_id"`
	TripTypes                 []string `json:"trip_types"`
	DepartureTime             string   `json:"departure_time,omitempty" validate:"optional-correct-rfc3339-date"`
	ArrivalTime               string   `json:"arrival_time,omitempty" validate:"optional-correct-rfc3339-date"`
	Region                    string   `json:"region" validate:"required"`
}

func parseTime(dt string) (time.Time, error) {
	zeroTime := time.Unix(0, 0)
	if dt == "" {
		return zeroTime, nil
	}

	t, err := time.Parse(time.RFC3339, dt)
	if err != nil {
		return zeroTime, err
	}

	return t, nil
}

func searchFindRoutesByLocationV2(c echo.Context) (err error) {
	var request FindRoutesRequest

	ctx := common.NewContext(c)

	err = common.BindAndValidateReq(c, &request)
	if err != nil {
		return
	}

	departureTime, err := parseTime(request.DepartureTime)
	if err != nil {
		return
	}

	arrivalTime, err := parseTime(request.ArrivalTime)
	if err != nil {
		return
	}

	tripTypes := make(map[string]bool)
	for i := 0; i < len(request.TripTypes); i++ {
		tripTypes[request.TripTypes[i]] = true
	}

	// TODO: глянуть на досуге влияет ли region == "en" на ошибки
	if strings.ToLower(request.Region) != "ru" {
		request.Region = "ru"
	}

	req := search.FindRoutesUsecaseIn{
		DeviceID:                  ctx.DeviceID,
		Origin:                    request.Origin,
		Destination:               request.Destination,
		OriginPlaceId:             request.OriginPlaceId,
		DestinationPlaceId:        request.DestinationPlaceId,
		OriginMovistaPlaceId:      request.OriginMovistaPlaceId,
		DestinationMovistaPlaceId: request.DestinationMovistaPlaceId,
		Region:                    request.Region,
		DepartureTime:             *models.NewTime(departureTime),
		ArrivalTime:               *models.NewTime(arrivalTime),
		TripTypes:                 tripTypes,
		ShowPathGroups:            true,
	}

	//TODO: find by location in DB if exists -> use ID
	result, err := search.FindRoutes(ctx, &req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"result": result,
	})
}

func searchFindRoutesByLocation(c echo.Context) (err error) {
	var request FindRoutesRequest

	ctx := common.NewContext(c)

	err = common.BindAndValidateReq(c, &request)
	if err != nil {
		return
	}

	fmt.Printf("Token: %s\n", c.Get("token"))

	if request.Origin.IsZero() && request.OriginPlaceId == "" && request.OriginMovistaPlaceId == nil {
		return fmt.Errorf("empty origin point")
	}

	if request.Destination.IsZero() && request.DestinationPlaceId == "" && request.DestinationMovistaPlaceId == nil {
		return fmt.Errorf("empty destination point")
	}

	var (
		departureTime = time.Unix(0, 0)
		arrivalTime   = time.Unix(0, 0)
	)

	if request.DepartureTime != "" && request.ArrivalTime == "" {
		dt, _ := time.Parse(time.RFC3339, request.DepartureTime)
		if dt.After(time.Now()) {
			departureTime = dt
		} else {
			departureTime = time.Now()
		}
	}

	if request.ArrivalTime != "" && request.DepartureTime == "" {
		at, _ := time.Parse(time.RFC3339, request.ArrivalTime)
		if at.After(time.Now()) {
			arrivalTime = at
		}
	}

	if request.ArrivalTime == "" && request.DepartureTime == "" {
		departureTime = time.Now()
	}

	if request.ArrivalTime != "" && request.DepartureTime != "" {
		dt, _ := time.Parse(time.RFC3339, request.DepartureTime)
		if dt.After(time.Now()) {
			departureTime = dt
		}
	}

	tripTypes := make(map[string]bool)
	for i := 0; i < len(request.TripTypes); i++ {
		tripTypes[request.TripTypes[i]] = true
	}

	req := search.FindRoutesUsecaseIn{
		DeviceID:                  ctx.DeviceID,
		Origin:                    request.Origin,
		Destination:               request.Destination,
		OriginPlaceId:             request.OriginPlaceId,
		DestinationPlaceId:        request.DestinationPlaceId,
		OriginMovistaPlaceId:      request.OriginMovistaPlaceId,
		DestinationMovistaPlaceId: request.DestinationMovistaPlaceId,
		Region:                    request.Region,
		DepartureTime:             *models.NewTime(departureTime),
		ArrivalTime:               *models.NewTime(arrivalTime),
		TripTypes:                 tripTypes,
	}

	// TODO: глянуть на досуге влияет ли region == "en" на ошибки
	if strings.ToLower(req.Region) != "ru" {
		req.Region = "ru"
	}

	//TODO: find by location in DB if exists -> use ID
	result, err := search.FindRoutes(ctx, &req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"result": result,
	})
}

func searchFindPlaceByName(c echo.Context) (err error) {
	var (
		req struct {
			Name       string   `json:"name"`
			Count      int      `json:"count" validate:"required"`
			PlaceTypes []string `json:"place_types,omitempty"`
		}
		ctx = common.NewContext(c)
	)

	err = common.BindAndValidateReq(c, &req)
	if err != nil {
		return
	}

	if req.Name == "" {
		return c.JSON(http.StatusOK, echo.Map{
			"result": make([]string, 0),
		})
	}

	if req.Count == 0 {
		req.Count = 3
	}

	result, err := place.FindByName(ctx, req.Name, req.Count, req.PlaceTypes)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, echo.Map{
		"result": result,
	})
}

func searchFindPlaceByIP(c echo.Context) (err error) {
	var (
		ctx = common.NewContext(c)
	)

	url := config.C.IpGeo.BaseURL + fmt.Sprintf("/%s", c.RealIP()) + ("?lang=ru")

	if strings.Contains(c.RealIP(), "192.168") {
		result, err := place.FindByName(ctx, "Москва", 1, nil)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, echo.Map{
			"data": result[0],
		})
	}

	placeTitle, err := common.FindPlaceTitleByIP(ctx, url)
	if err != nil {
		return err
	}

	result, err := place.FindByName(ctx, placeTitle, 1, nil)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": result[0],
	})
}

func searchFindPlaceByID(c echo.Context) (err error) {
	var (
		req struct {
			SessionToken string `json:"token"`
			Region       string `json:"region" validate:"required"`
		}
		ctx = common.NewContext(c)
	)

	placeID := c.Param("placeID")
	if placeID == "" {
		return fmt.Errorf("empty placeID param")
	}

	err = common.BindAndValidateReq(c, &req)
	if err != nil {
		return err
	}

	sessionUUID, _ := uuid.Parse(req.SessionToken)

	result, err := google.FindByID(ctx, sessionUUID, placeID, req.Region)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, echo.Map{
		"result": result,
	})
}

func searchGetPlaceHistory(c echo.Context) (err error) {
	var (
		req struct {
			Count int `json:"count" validate:"required"`
		}
		ctx = common.NewContext(c)
	)

	err = common.BindAndValidateReq(c, &req)
	if err != nil {
		return
	}

	result, err := history.GetGooglePlaceMentions(ctx)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, echo.Map{
		"result": result,
	})
}

func searchGetPlaceFavorite(c echo.Context) (err error) {
	var (
		req struct {
			Count int `json:"count"`
		}
		ctx = common.NewContext(c)
	)

	err = common.BindAndValidateReq(c, &req)
	if err != nil {
		return
	}

	favs, err := favorite.FindFavoriteByDeviceIDWithGooglePlaceInfo(ctx)
	if err != nil {
		return
	}

	places := make([]*google.Place, 0)

	for _, p := range favs {
		response := &google.Place{
			ID:                   p.GooglePlace.ID,
			PlaceTypes:           p.GooglePlace.PlaceTypes,
			MainText:             p.GooglePlace.MainText,
			SecondaryText:        p.GooglePlace.SecondaryText,
			Coordinate:           p.GooglePlace.Coordinate,
			Name:                 p.Name,
			Type:                 p.Type,
			IosIconURLDarkTheme:  p.GooglePlace.IosIconURLDarkTheme,
			IosIconURL:           p.GooglePlace.IosIconURLDarkTheme,
			IosIconURLLightTheme: p.GooglePlace.IosIconURLLightTheme,
			AndroidIconURL:       p.GooglePlace.AndroidIconURL,
		}

		places = append(places, response)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"result": places,
	})
}

func searchSavePlace(c echo.Context) (err error) {
	var (
		req struct {
			UserPlaceType string           `json:"type"`
			Location      *models.GeoPoint `json:"location" validate:"required"`
			PlaceID       string           `json:"place_id"`
			PlaceTypes    []string         `json:"types"`
			MainText      string           `json:"main_text"`
			SecondaryText string           `json:"secondary_text"`
		}
		ctx = common.NewContext(c)
	)

	err = common.BindAndValidateReq(c, &req)
	if err != nil {
		return
	}

	placeID, err := google.Save(ctx, &google.Place{
		ID:            req.PlaceID,
		MainText:      req.MainText,
		SecondaryText: req.SecondaryText,
		Coordinate:    req.Location,
		PlaceTypes:    req.PlaceTypes,
	})
	if err != nil {
		return
	}

	if req.UserPlaceType == "history" {
		err = history.SaveGooglePlaceMention(ctx, placeID)
		return
	}

	_, _, err = favorite.Create(ctx, &favorite.DeviceFavorite{
		DeviceID: ctx.DeviceID,
		PlaceID:  placeID,
		Type:     req.UserPlaceType,
	})
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, echo.Map{})
}

func searchAsyncStatus(c echo.Context) (err error) {
	return c.JSON(http.StatusOK, echo.Map{"status": config.C.AsyncAvailable})
}

func placesByIds(c echo.Context) error {
	var (
		req place.PlaceIdsRequest
		ctx = common.NewContext(c)
	)

	err := common.BindAndValidateReq(c, &req)
	if err != nil {
		return err
	}

	resp, err := place.PlacesByIds(ctx, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}
