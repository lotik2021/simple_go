package controller

import (
	"net/http"

	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/dictionary"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"github.com/labstack/echo/v4"
)

const (
	RU = "ru"
)

func Citizenships(c echo.Context) error {
	var (
		ctx = common.NewContext(c)
		in  struct {
			CultureCode string `json:"culture_code,omitempty"`
		}
	)

	if err := common.BindAndValidateReq(c, &in); err != nil {
		return err
	}

	if in.CultureCode == "" {
		in.CultureCode = RU
	}

	citShips, err := dictionary.GetCitizenships(ctx, in.CultureCode)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &models.UapiResponse{
		Data: map[string]interface{}{
			"сitizenships": citShips,
		},
	},
	)
}
