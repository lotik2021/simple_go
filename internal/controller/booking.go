package controller

import (
	"fmt"
	"net/http"

	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/uapi/fapi"
	"github.com/labstack/echo/v4"
)

func checkBookingV4(c echo.Context) error {
	var (
		req struct {
			UID string `json:"uid" validate:"required"`
		}
		ctx = common.NewContext(c)
	)

	err := common.BindAndValidateReq(c, &req)
	if err != nil {
		return err
	}

	rawResp, err := fapi.CheckBookingV4(ctx, req.UID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, rawResp)
}

func checkBookingV2(c echo.Context) error {
	var (
		ctx = common.NewContext(c)
	)

	uid := c.Param("uid")
	if uid == "" {
		return fmt.Errorf("uid is empty")
	}

	rawResp, err := fapi.CheckBookingV4(ctx, uid)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, rawResp)
}
