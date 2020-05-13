package controller

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/favorite"
	"bitbucket.movista.ru/maas/maasapi/internal/google"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type FavoriteIn struct {
	DeviceID      string           `json:"user_id" validate:"required"`
	Location      *models.GeoPoint `json:"location" validate:"required"`
	PlaceID       string           `json:"place_id"`
	PlaceTypes    []string         `json:"types"`
	MainText      string           `json:"main_text"`
	SecondaryText string           `json:"secondary_text"`
	UserPlaceType string           `json:"type" validate:"oneof=work home none"`
	Name          string           `json:"name"`
}

type CreateFavoriteOut struct {
	ID int `json:"id"`
}

func createFavorite(c echo.Context) (err error) {
	var (
		req FavoriteIn
		ctx = common.NewContext(c)
	)

	err = common.BindAndValidateReq(c, &req)
	if err != nil {
		return err
	}

	placeID, err := google.Save(ctx, &google.Place{
		ID:            req.PlaceID,
		Coordinate:    req.Location,
		MainText:      req.MainText,
		SecondaryText: req.SecondaryText,
		PlaceTypes:    req.PlaceTypes,
	})

	if err != nil {
		return
	}

	id, _, err := favorite.Create(ctx, &favorite.DeviceFavorite{
		DeviceID: req.DeviceID,
		PlaceID:  placeID,
		Type:     req.UserPlaceType,
		Name:     req.Name,
	})

	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, CreateFavoriteOut{ID: id})
}

func updateFavorite(c echo.Context) (err error) {
	var (
		req FavoriteIn
		ctx = common.NewContext(c)
	)

	idParam := c.Param("id")
	if idParam == "" {
		return fmt.Errorf("empty id in params")
	}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return
	}

	err = common.BindAndValidateReq(c, &req)
	if err != nil {
		return err
	}

	placeID, err := google.Save(ctx, &google.Place{
		ID:            req.PlaceID,
		Coordinate:    req.Location,
		MainText:      req.MainText,
		SecondaryText: req.SecondaryText,
		PlaceTypes:    req.PlaceTypes,
	})
	if err != nil {
		return
	}

	fav := &favorite.DeviceFavorite{
		ID:       id,
		DeviceID: req.DeviceID,
		PlaceID:  placeID,
		Type:     req.UserPlaceType,
		Name:     req.Name,
	}

	err = favorite.UpdateFavorite(ctx, fav)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": "success",
	})
}

func deleteFavorite(c echo.Context) (err error) {
	var (
		ctx = common.NewContext(c)
	)

	id := c.Param("id")
	if id == "" {
		return fmt.Errorf("empty id in params")
	}

	intId, _ := strconv.Atoi(id)

	err = favorite.DeleteFavorite(ctx, intId)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": "success",
	})
}

type ListOut struct {
	Data []favorite.WithGooglePlaceInfo `json:"data"`
}

func listFavorite(c echo.Context) (err error) {
	var (
		ctx = common.NewContext(c)
	)

	result, err := favorite.FindFavoriteByDeviceIDWithGooglePlaceInfo(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ListOut{Data: result})
}

func toggleFavoriteInActions(c echo.Context) (err error) {
	var (
		ctx = common.NewContext(c)
	)

	id := c.Param("id")
	if id == "" {
		return fmt.Errorf("empty id in params")
	}

	intId, _ := strconv.Atoi(id)

	err = favorite.ToggleFavoriteInActions(ctx, intId)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": "success",
	})
}
