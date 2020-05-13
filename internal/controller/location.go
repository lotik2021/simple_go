package controller

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/device"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

func pingLocation(c echo.Context) (err error) {
	var (
		req struct {
			Location *models.GeoPoint `json:"location" validate:"required"`
		}
		ctx = common.NewContext(c)
	)

	err = common.BindAndValidateReq(c, &req)
	if err != nil {
		return
	}

	err = device.PingLocation(ctx, req.Location)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": "ok",
	})
}
