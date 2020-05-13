package payment

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/parnurzeal/gorequest"
)

var (
	paymentClient *gorequest.SuperAgent
)

func init() {
	paymentClient = common.DefaultRequest.Clone().Timeout(config.C.Booking.RequestTimeout)
}

func Pay(ctx common.Context, paymentReq PayRequest) (response json.RawMessage, err error) {
	uapiPayReq := UapiPayRequest{
		OrderId:     paymentReq.OrderId,
		RequestId:   paymentReq.Uid,
		BackUrl:     config.C.Payment.Processing + fmt.Sprintf("?email=%s&orderId=%d", paymentReq.Email, paymentReq.OrderId),
		PaymentType: paymentReq.PaymentType,
	}

	rawResp, err := common.UapiAuthorizedRequest(ctx, paymentClient, uapiPayReq, http.MethodPost, config.C.Payment.Urls.Pay, common.GetInternalToken())
	if err != nil {
		return
	}

	return rawResp.Data, nil
}

func OrderState(ctx common.Context, orderStateReq OrderStateRequest) (response json.RawMessage, err error) {
	url := config.C.Payment.Urls.OrderState + fmt.Sprintf("?OrderId=%d&RequestId=%s", orderStateReq.OrderId, orderStateReq.Uid)

	rawResp, err := common.UapiAuthorizedRequest(ctx, paymentClient, nil, http.MethodGet, url, common.GetInternalToken())
	if err != nil {
		return
	}

	return rawResp.Data, nil
}

func GetOrderStateV1(c echo.Context) error {
	var (
		ctx = common.NewContext(c)
		req struct {
			Uid     string `json:"uid"  query:"uid" validate:"required"`
			OrderId string `json:"orderId"  query:"orderId" validate:"required"`
		}
	)

	if err := common.BindAndValidateReq(c, &req); err != nil {
		return err
	}

	u, _ := url.Parse(config.C.Payment.Urls.GetOrderStateV1)
	v := url.Values{}
	v.Set("OrderId", req.OrderId)
	v.Set("RequestId", req.Uid)
	u.RawQuery = v.Encode()

	resp, err := common.UapiAuthorizedGet(ctx, paymentClient, u.String())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func RefundOrderV1(c echo.Context) error {
	var (
		ctx = common.NewContext(c)
		in  RefundOrderRequest
	)

	if err := common.BindAndValidateReq(c, &in); err != nil {
		return err
	}

	resp, err := common.UapiAuthorizedPost(ctx, paymentClient, in, config.C.Payment.Urls.RefundOrderV1)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func AdditionalPay(ctx common.Context, req AdditionalPayRequest) (*UPayData, error) {
	upayReq := UAdditionalPayRequest{
		EventId:     req.EventId,
		BackUrl:     config.C.Payment.Processing,
		PaymentType: req.PaymentType,
	}

	resp, err := common.UapiAuthorizedPost(ctx, paymentClient, upayReq, config.C.Payment.Urls.AdditionalPay)
	if err != nil {
		return nil, err
	}

	payResp := &UPayData{}
	err = json.Unmarshal(resp.Data, payResp)
	if err != nil {
		return nil, err
	}

	return payResp, nil
}
