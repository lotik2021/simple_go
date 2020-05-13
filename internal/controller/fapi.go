package controller

import (
	"io/ioutil"
	"net/http"

	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/history"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"bitbucket.movista.ru/maas/maasapi/internal/ratelimit"
	"bitbucket.movista.ru/maas/maasapi/internal/uapi/fapi"
	"github.com/labstack/echo/v4"
)

func searchSearchAsync(c echo.Context) (err error) {
	var (
		asyncReq fapi.AsyncRequestV5
		ctx      = common.NewContext(c)
	)

	err = ratelimit.Apply(ctx, ratelimit.MethodSearchAsync)
	if err != nil {
		return
	}

	if err = common.BindAndValidateReq(c, &asyncReq); err != nil {
		return err
	}

	asyncResp, err := fapi.SearchAsyncV5(ctx, asyncReq)
	if err != nil {
		return err
	}

	history.SaveMovistaWebSearch(ctx, asyncReq.SearchParams, asyncResp.UID)

	return c.JSON(http.StatusOK, models.UapiResponse{
		Data: asyncResp,
	})

}

func searchGetSearchStatus(c echo.Context) (err error) {
	var (
		request struct {
			UID string `json:"uid" validate:"required"`
		}
		ctx = common.NewContext(c)
	)

	err = common.BindAndValidateReq(c, &request)
	if err != nil {
		return err
	}

	result, err := fapi.GetSearchStatusV5(ctx, request.UID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

func searchGetSearchResults(c echo.Context) (err error) {
	var (
		request struct {
			UID string `json:"uid"`
		}
		ctx = common.NewContext(c)
	)

	err = common.BindAndValidateReq(c, &request)
	if err != nil {
		return err
	}

	result, err := fapi.GetSearchResultsV5(ctx, request.UID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

func searchGetPathGroup(c echo.Context) (err error) {
	var (
		request struct {
			UID         string `json:"uid"`
			PathGroupId string `json:"pathGroupId"`
		}
		ctx = common.NewContext(c)
	)

	err = common.BindAndValidateReq(c, &request)
	if err != nil {
		return err
	}

	result, err := fapi.GetPathGroupV5(ctx, request.UID, request.PathGroupId)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

func searchGetSegmentRoutes(c echo.Context) (err error) {

	var (
		body []byte
		ctx  = common.NewContext(c)
	)

	body, err = ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	result, err := fapi.GetSegmentRoutesV5(ctx, body)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

func searchRefundOrderAsync(c echo.Context) (err error) {
	var (
		request struct {
			JobID string `json:"job_id"`
		}
		ctx = common.NewContext(c)
	)

	err = common.BindAndValidateReq(c, &request)
	if err != nil {
		return err
	}

	result, err := fapi.RefundOrderAsync(ctx, request.JobID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

func searchGetTrainSchedule(c echo.Context) (err error) {
	var (
		request struct {
			Origin      models.GeoPoint `json:"origin"`
			Destination models.GeoPoint `json:"destination"`
			Departure   string          `json:"departure"`
			Count       int             `json:"count,omitempty"`
		}
		ctx = common.NewContext(c)
	)

	err = common.BindAndValidateReq(c, &request)
	if err != nil {
		return err
	}

	result, err := fapi.GetTrainSchedule(ctx, request.Origin, request.Destination, request.Departure, request.Count)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"data": result})
}
