package controller

import (
	"net/http"

	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/notification"
	"github.com/labstack/echo/v4"
)

func sendPushNotification(c echo.Context) (err error) {
	var (
		req notification.PushNotificationsRequest
		ctx = common.NewContext(c)
	)

	err = common.BindAndValidateReq(c, &req)
	if err != nil {
		return
	}

	err = notification.SendPushNotification(ctx, req)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": "success",
	})
}

func sendFeedback(c echo.Context) error {
	var (
		req notification.Notification
		ctx = common.NewContext(c)
	)

	err := common.BindAndValidateReq(c, &req)
	if err != nil {
		return err
	}

	resp, err := notification.SendFeedback(ctx, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}
