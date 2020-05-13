package controller

import (
	"net/http"

	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/content"
	"github.com/labstack/echo/v4"
)

func queryCms(c echo.Context) error {
	var (
		ctx = common.NewContext(c)
		in  content.CmsContent
	)

	if err := common.BindAndValidateReq(c, &in); err != nil {
		return err
	}

	rawResp, err := content.QueryCmsV1(ctx, in)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, rawResp)
}
