package dialog

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/device"
	"bitbucket.movista.ru/maas/maasapi/internal/favorite"
	"bitbucket.movista.ru/maas/maasapi/internal/google"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
	"time"
)

func updateProfile(c echo.Context) (err error) {
	var (
		req struct {
			Name        string `json:"name"`
			NameAskedAt string `json:"last_name_question_time" validate:"optional-correct-rfc3339-date"`
		}
		ctx = common.NewContext(c)
	)

	if err = common.BindAndValidateReq(c, &req); err != nil {
		return models.NewInternalDialogError(err)
	}

	d := device.Device{ID: ctx.DeviceID}

	if req.NameAskedAt != "" {
		t, _ := time.Parse(time.RFC3339, req.NameAskedAt)
		d.NameAskedAt = models.Time{Time: t}
	}

	if req.Name != "" {
		d.Name = strings.Title(req.Name)
	}

	err = device.Update(ctx, &d)
	if err != nil {
		return models.NewInternalDialogError(err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"result": d,
	})
}

func getProfile(c echo.Context) (err error) {
	var (
		ctx = common.NewContext(c)
	)

	d, err := device.GetOne(ctx)
	if err != nil {
		return models.NewInternalDialogError(err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"result": d,
	})
}

type GooglePlace struct {
	Coordinate *models.GeoPoint `json:"coordinate"`
	google.Place
}

type ListItem struct {
	ID          int         `json:"id"`
	DeviceID    string      `json:"device_id"`
	Type        string      `json:"favorite_place_type"`
	Name        string      `json:"name"`
	InActions   bool        `json:"in_actions"`
	GooglePlace GooglePlace `json:"google_place,omitempty"`
}

type getDeviceFavoritesOut struct {
	Result []ListItem `json:"result"`
}

func getDeviceFavorites(c echo.Context) (err error) {
	var (
		req struct {
			Type string `json:"type"`
		}
		ctx = common.NewContext(c)
	)

	if err = common.BindAndValidateReq(c, &req); err != nil {
		return models.NewInternalDialogError(err)
	}

	var list []favorite.DeviceFavorite

	if req.Type != "" {
		list, err = favorite.FindFavoriteByDeviceIDAndType(ctx, req.Type)
	} else {
		list, err = favorite.FindFavorite(ctx)
	}

	if err != nil {
		return models.NewInternalDialogError(err)
	}

	if len(list) == 0 {
		return c.JSON(http.StatusOK, getDeviceFavoritesOut{})
	}

	var googlePlaceIDs []string
	for _, li := range list {
		googlePlaceIDs = append(googlePlaceIDs, li.PlaceID)
	}

	res, err := google.FindByIDsInDB(ctx, googlePlaceIDs)
	if err != nil {
		return models.NewInternalDialogError(err)
	}

	resultList := make([]ListItem, 0)

	for _, item := range list {
		gp, ok := res[item.PlaceID]
		if !ok {
			continue
		}

		li := ListItem{
			ID:        item.ID,
			DeviceID:  item.DeviceID,
			Type:      item.Type,
			Name:      item.Name,
			InActions: item.InActions,
			GooglePlace: GooglePlace{
				Coordinate: gp.Coordinate,
				Place: google.Place{
					ID:                   gp.ID,
					MainText:             gp.MainText,
					SecondaryText:        gp.SecondaryText,
					PlaceTypes:           gp.PlaceTypes,
					IosIconURLLightTheme: gp.IosIconURLLightTheme,
					IosIconURLDarkTheme:  gp.IosIconURLDarkTheme,
					IosIconURL:           gp.IosIconURLDarkTheme,
					AndroidIconURL:       gp.AndroidIconURL,
				},
			},
		}

		resultList = append(resultList, li)
	}

	return c.JSON(http.StatusOK, getDeviceFavoritesOut{Result: resultList})
}
