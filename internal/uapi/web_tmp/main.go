package web_tmp

import (
	"net/http"

	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/parnurzeal/gorequest"
)

var (
	webTmpClient *gorequest.SuperAgent
)

func init() {
	webTmpClient = common.DefaultRequest.Clone().Timeout(config.C.RequestTimeout)
}

func GetEventFromOrderV2(c echo.Context) error {
	ctx := common.NewContext(c)
	url := config.C.WebTmp.Urls.OrderEvent + c.Param("eventId")
	resp, err := common.UapiAuthorizedGet(ctx, webTmpClient, url)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, resp)
}
