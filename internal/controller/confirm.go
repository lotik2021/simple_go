package controller

import (
	"net/http"

	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/uapi/auth"
	"github.com/labstack/echo/v4"
)

func confirmEmail(c echo.Context) error {
	var (
		ctx = common.NewContext(c)
		in  auth.ConfirmEmailRequest
	)

	if err := common.BindAndValidateReq(c, &in); err != nil {
		return err
	}

	rawResp, err := auth.ConfirmEmail(ctx, in)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, rawResp)
}
