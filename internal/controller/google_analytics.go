package controller

import (
	"io/ioutil"
	"net/http"

	"bitbucket.movista.ru/maas/maasapi/internal/analytics"
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"github.com/labstack/echo/v4"
)

func sendGA(c echo.Context) error {
	ctx := common.NewContext(c)

	var err error
	body := []byte{}
	if c.Request().Body != nil {
		body, err = ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}
	}

	analytics.SendGoogleAnalytics(ctx, string(body))
	return c.JSON(http.StatusOK, nil)
}
