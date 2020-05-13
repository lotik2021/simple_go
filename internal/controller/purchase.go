package controller

import (
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"github.com/labstack/echo/v4"
	"net/http"
)

func getPurchaseHints(c echo.Context) (err error) {
	return c.JSON(http.StatusOK, echo.Map{"hints": config.PromptsPlacesTariffs})
}
