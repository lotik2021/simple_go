package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"bitbucket.movista.ru/maas/maasapi/internal/weather"
	"github.com/labstack/echo/v4"
)

func getWeather(c echo.Context) (err error) {
	var (
		req struct {
			Location models.GeoPoint `json:"location"`
		}
		ctx = common.NewContext(c)
	)

	err = common.BindAndValidateReq(c, &req)
	if err != nil {
		return models.NewInternalDialogError(err)
	}

	resp, err := weather.GetWeather(ctx, req.Location)
	if err != nil {
		return models.NewInternalDialogError(err)
	}

	result := &weather.OpenWeatherResponse{}

	err = json.Unmarshal(resp, result)
	if err != nil {
		err = fmt.Errorf("cannot unmarshal OpenWeather response - %s", resp)
		return
	}

	result.GetIcons()

	return c.JSON(http.StatusOK, result)
}

func getForecast(c echo.Context) (err error) {
	var (
		in struct {
			Location models.GeoPoint `json:"location" validate:"required"`
			From     string          `json:"from" validate:"required,optional-correct-rfc3339-date"`
			To       string          `json:"to" validate:"required,optional-correct-rfc3339-date"`
		}
		ctx = common.NewContext(c)
	)

	err = common.BindAndValidateReq(c, &in)
	if err != nil {
		return models.NewInternalDialogError(err)
	}

	var fromUnix, toUnix int64

	if t, err := time.Parse(time.RFC3339, in.From); err == nil {
		fromUnix = t.Unix()
	} else {
		return models.NewInternalDialogError(err)
	}

	if t, err := time.Parse(time.RFC3339, in.To); err == nil {
		toUnix = t.Unix()
	} else {
		return models.NewInternalDialogError(err)
	}

	resp, err := weather.GetForecast(ctx, in.Location)
	if err != nil {
		return models.NewInternalDialogError(err)
	}

	forecast := &weather.ForecastResponse{}
	err = json.Unmarshal(resp, forecast)
	if err != nil {
		return models.NewInternalDialogError(err)
	}

	if forecast == nil {
		fmt.Errorf("forecast response is empty - %s", resp)
	}

	startI := 0
	endI := 0

	for i := 0; i < forecast.Cnt; i++ {
		if forecast.List[i].Dt >= fromUnix {
			startI = i
			break
		}
	}

	for i := startI + 1; i < forecast.Cnt; i++ {
		if forecast.List[i].Dt > toUnix {
			endI = i
			break
		}
	}

	forecast.List = forecast.List[startI:endI]

	for i := range forecast.List {
		forecast.List[i].GetIcons()
	}

	forecast.Cnt = len(forecast.List)

	return c.JSON(http.StatusOK, forecast)
}
