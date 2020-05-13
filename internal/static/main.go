package static

import (
	"fmt"
	"net/http"

	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/parnurzeal/gorequest"
)

var (
	staticClient *gorequest.SuperAgent
)

func init() {
	staticClient = common.DefaultRequest.Clone().Timeout(config.C.RequestTimeout)
}

func LoadContent(c echo.Context) error {
	var (
		ctx = common.NewContext(c)
		in  struct {
			Name string `json:"name"`
		}
	)

	err := common.BindAndValidateReq(c, &in)
	if err != nil {
		return err
	}

	if in.Name == "aboutCustomer" {
		body, _, err := common.SendRequest(ctx, staticClient.Clone().Get(config.C.Static.AboutCustomer))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, models.UapiResponse{
			Data: echo.Map{
				"content": string(body),
			},
		})
	}

	return fmt.Errorf("no such file: %s", in.Name)
}
