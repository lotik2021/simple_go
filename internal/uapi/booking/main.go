package booking

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"bitbucket.movista.ru/maas/maasapi/internal/uapi/payment"
	"github.com/labstack/echo/v4"
	"github.com/parnurzeal/gorequest"
)

var (
	bookingClient *gorequest.SuperAgent
)

func init() {
	bookingClient = common.DefaultRequest.Clone().Timeout(config.C.Booking.RequestTimeout)
}

func GetFareRulesV2(c echo.Context) (err error) {

	var (
		body []byte
		ctx  = common.NewContext(c)
	)

	body, err = ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	req := bookingClient.Clone().Post(config.C.Booking.Urls.GetFareRulesV2).SendStruct(json.RawMessage(body))

	resp, _, err := common.SendRequest(ctx, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, json.RawMessage(resp))
}

func CancelOrder(c echo.Context) (err error) {

	var (
		ctx = common.NewContext(c)
	)

	resp, err := common.UapiPost(ctx, bookingClient, c.Request().Body, config.C.Booking.Urls.CancelOrder)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func CheckRefundV1(c echo.Context) error {
	var (
		ctx = common.NewContext(c)
		in  payment.RefundOrderRequest
	)

	if err := common.BindAndValidateReq(c, &in); err != nil {
		return err
	}

	rawResp, err := common.UapiAuthorizedPost(ctx, bookingClient, in, config.C.Booking.Urls.CheckRefundV1)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, rawResp)
}

func GetFareRulesV1(c echo.Context) error {
	var (
		ctx = common.NewContext(c)
		in  struct {
			OrderId          int64       `json:"orderId,omitempty"`
			Service          interface{} `json:"service"`
			AlternativeIndex int         `json:"alternativeIndex"`
			Uid              string      `json:"uid" validate:"required"`
			Currencycode     string      `json:"currencyCode" validate:"required"`
		}
	)

	if err := common.BindAndValidateReq(c, &in); err != nil {
		return err
	}

	rawResp, err := common.UapiAuthorizedPost(ctx, bookingClient, in, config.C.Booking.Urls.GetFareRulesV1)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, rawResp)
}

func ChangeERegistrationV1(c echo.Context) error {
	var (
		ctx = common.NewContext(c)
		req payment.RefundOrderRequest
	)

	if err := common.BindAndValidateReq(c, &req); err != nil {
		return err
	}

	resp, err := common.UapiAuthorizedPost(ctx, bookingClient, req, config.C.Booking.Urls.ChangeERegistrationV1)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}
