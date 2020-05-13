package controller

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/device"
	"bitbucket.movista.ru/maas/maasapi/internal/transport/travelcard"
	"github.com/labstack/echo/v4"
	"net/http"
)

func payTroika(c echo.Context) (err error) {
	var (
		req struct {
			Amount     int    `json:"amount" validate:"required"`
			CardNumber string `json:"card_number" validate:"required"`
		}
		ctx = common.NewContext(c)
	)

	err = common.BindAndValidateReq(c, &req)
	if err != nil {
		return
	}

	redirectURL, err := device.PayTroika(ctx, req.Amount, req.CardNumber)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, echo.Map{
		"redirect_url": redirectURL,
	})
}

func getTroikaParameters(c echo.Context) (err error) {
	ctx := common.NewContext(c)

	parameters, err := travelcard.GetTroikaParameters(ctx)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, echo.Map{
		"parameters": parameters,
	})
}

func getStrelkaCardParameters(c echo.Context) (err error) {
	ctx := common.NewContext(c)

	parameters, err := travelcard.GetStrelkaParameters(ctx)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, echo.Map{
		"parameters": parameters,
	})
}

func getStrelkaPaymentParameters(c echo.Context) (err error) {

	var (
		req struct {
			CardNumber string `json:"card_number" validate:"required"`
		}
		ctx = common.NewContext(c)
	)

	err = common.BindAndValidateReq(c, &req)
	if err != nil {
		return
	}

	balance, cardTypeID, redirectURL, parameters, err := travelcard.GetStrelkaBalance(ctx, req.CardNumber)
	if err != nil {
		return
	}

	if redirectURL != "" {
		return c.JSON(http.StatusOK, echo.Map{
			"redirect_url": redirectURL,
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"parameters":   parameters,
		"balance":      balance,
		"card_type_id": cardTypeID,
	})
}

func payStrelka(c echo.Context) (err error) {

	var (
		req struct {
			CardNumber string `json:"card_number" validate:"required"`
			CardTypeID string `json:"card_type_id" validate:"required"`
			Amount     int    `json:"amount" validate:"required"`
		}
		ctx = common.NewContext(c)
	)

	err = common.BindAndValidateReq(c, &req)
	if err != nil {
		return
	}

	redirectURL, err := device.PayStrelka(ctx, req.Amount, req.CardNumber, req.CardTypeID)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, echo.Map{
		"redirect_url": redirectURL,
	})
}

func getTravelCards(c echo.Context) (err error) {
	ctx := common.NewContext(c)

	cards, err := travelcard.GetTravelCards(ctx)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, echo.Map{
		"travel_cards": cards,
	})
}

func savePayment(c echo.Context) (err error) {
	var (
		req struct {
			PaymentId string `json:"payment_id" validate:"required"`
			Success   int    `json:"success" validate:"required"`
		}
		ctx = common.NewContext(c)
	)

	err = common.BindAndValidateReq(c, &req)
	if err != nil {
		return
	}

	processed := req.Success == 1

	err = device.UpdatePaymentProcessedStatus(ctx, req.PaymentId, processed)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, echo.Map{"result": "success"})
}

func getPayments(c echo.Context) (err error) {
	ctx := common.NewContext(c)

	payments, err := device.GetPayments(ctx)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": payments,
	})
}
