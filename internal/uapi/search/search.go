package search

import (
	"net/http"

	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/parnurzeal/gorequest"
)

var (
	routeClient *gorequest.SuperAgent
)

func init() {
	routeClient = common.DefaultRequest.Clone().Timeout(config.C.RequestTimeout)
}

func RouteDetailsV3(c echo.Context) error {
	var (
		ctx = common.NewContext(c)
		in  struct {
			Number              string `json:"number"`
			FromId              uint32 `json:"fromId"`
			ToId                uint32 `json:"toId"`
			DepartureDate       string `json:"departureDate"`
			ProviderServiceCode string `json:"providerServiceCode"`
		}
	)

	if err := common.BindAndValidateReq(c, &in); err != nil {
		return err
	}

	resp, err := common.UapiAuthorizedPost(ctx, routeClient, in, config.C.Search.Urls.RouteDetailsV3)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}
