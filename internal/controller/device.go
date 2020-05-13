package controller

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/device"
	"github.com/labstack/echo/v4"
	"net/http"
)

func setOSPlayerID(c echo.Context) (err error) {
	var (
		req struct {
			OSPlayerID string `json:"os_player_id" validate:"required"`
		}
		ctx = common.NewContext(c)
	)

	err = common.BindAndValidateReq(c, &req)
	if err != nil {
		return err
	}

	err = device.Update(ctx, &device.Device{ID: ctx.DeviceID, OsPlayerID: req.OSPlayerID})
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": "success",
	})
}

func getDeviceTransportSettings(c echo.Context) (err error) {
	ctx := common.NewContext(c)

	resp, err := device.GetGoogleTransitModeSettingsWithInfo(ctx)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, echo.Map{
		"transport": resp,
	})
}

func setUsedLinkInTaxiOrder(c echo.Context) (err error) {
	var (
		req struct {
			ID       string `json:"id" validate:"required"`
			UsedLink string `json:"used_link" validate:"required"`
		}
		ctx = common.NewContext(c)
	)

	err = common.BindAndValidateReq(c, &req)
	if err != nil {
		return
	}

	err = device.SetUsedLinkInTaxiOrder(ctx, req.ID, req.UsedLink)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": "ok",
	})
}

func getDeviceProfile(c echo.Context) (err error) {

	ctx := common.NewContext(c)

	v, err := device.GetOneWithSettingsAndFavorites(ctx)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, v)
}

func updateDeviceGoogleTransportSettings(c echo.Context) (err error) {
	var (
		req struct {
			Transports []struct {
				Name   string `json:"name" validate:"required"`
				Status bool   `json:"status" validate:"required"`
			} `json:"transports" validate:"required"`
		}
		ctx = common.NewContext(c)
	)

	err = common.BindAndValidateReq(c, &req)
	if err != nil {
		return
	}

	transitModes := make([]string, 0)

	for _, t := range req.Transports {
		if t.Status {
			transitModes = append(transitModes, t.Name)
		}
	}

	err = device.Update(ctx, &device.Device{ID: ctx.DeviceID, GoogleTransitModes: transitModes})
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": "ok",
	})
}
