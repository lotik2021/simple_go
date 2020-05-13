package controller

import (
	"net/http"

	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"bitbucket.movista.ru/maas/maasapi/internal/uapi/payment"
	"github.com/labstack/echo/v4"
)

func paymentPay(c echo.Context) (err error) {
	var (
		req payment.PayRequest
		ctx = common.NewContext(c)
	)

	err = common.BindAndValidateReq(c, &req)
	if err != nil {
		return
	}

	resp, err := payment.Pay(ctx, req)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": resp,
	})
}

func paymentOrderState(c echo.Context) (err error) {
	var (
		req payment.OrderStateRequest
		ctx = common.NewContext(c)
	)

	err = common.BindAndValidateReq(c, &req)
	if err != nil {
		return
	}

	resp, err := payment.OrderState(ctx, req)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": resp,
	})
}

func additionalPay(c echo.Context) error {
	var (
		req payment.AdditionalPayRequest
		ctx = common.NewContext(c)
	)

	err := common.BindAndValidateReq(c, &req)
	if err != nil {
		return err
	}

	payResp, err := payment.AdditionalPay(ctx, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, models.UapiResponse{
		Data: payResp,
	})
}
