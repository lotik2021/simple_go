package controller

import (
	"net/http"
	"strings"

	"bitbucket.movista.ru/maas/maasapi/internal/analytics"
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"bitbucket.movista.ru/maas/maasapi/internal/constant"
	"bitbucket.movista.ru/maas/maasapi/internal/device"
	"bitbucket.movista.ru/maas/maasapi/internal/logger"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"bitbucket.movista.ru/maas/maasapi/internal/ratelimit"
	"bitbucket.movista.ru/maas/maasapi/internal/uapi/auth"
	"bitbucket.movista.ru/maas/maasapi/internal/user"
	"github.com/labstack/echo/v4"
)

type CreateTemporaryUserReq struct {
	DeviceID       string `json:"device_id" validate:"required" example:"something-unique"`
	OS             string `json:"os" validate:"required" example:"ios" enums:"ios,android,vk"`
	DeviceInfo     string `json:"device_info" validate:"required" example:"iPhone XS"`
	DeviceCategory string `json:"device_category"`
	DeviceType     string `json:"device_type"`
}

type CreateTemporaryUserRes struct {
	Token       string `json:"token"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// authCreateTemporaryUser godoc
// @Tags auth
// @Summary authCreateTemporaryUser
// @Description create temporary user
// @ID create-temporary-user
// @Accept json
// @Produce json
// @Param req body CreateTemporaryUserReq false "Device"
// @Success 200 {object} CreateTemporaryUserRes
// @Router /api/auth/create [post]
func authCreate(c echo.Context) (err error) {

	var (
		req CreateTemporaryUserReq
		ctx = common.NewContext(c)
	)

	err = common.BindAndValidateReq(c, &req)
	if err != nil {
		return err
	}

	err = ratelimit.Apply(ctx, ratelimit.MethodAuthCreate)
	if err != nil {
		return
	}

	d := &device.Device{
		ID:         req.DeviceID,
		DeviceOS:   req.OS,
		DeviceInfo: req.DeviceInfo,
	}

	if req.DeviceType == "" && req.DeviceCategory == "" {
		d.DeviceCategory, d.DeviceType = common.FetchDeviceCategoryAndType(req.DeviceID, req.DeviceInfo, req.OS)
	} else {
		d.DeviceCategory, d.DeviceType = req.DeviceCategory, req.DeviceType
	}

	err = device.Create(ctx, d)
	if err != nil {
		return
	}

	token := models.NewDeviceToken(models.DeviceClaims{DeviceID: d.ID})

	accessToken, expiresIn, err := common.GetLegacyTokenForWebDevice(ctx)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, CreateTemporaryUserRes{Token: token.String(), AccessToken: accessToken, ExpiresIn: expiresIn})
}

type SendSmsCodeReq struct {
	Phone string `json:"phone" validate:"required" example:"+7999887766"`
}

type SendSmsCodeRes struct {
	RefreshAfter int `json:"refresh_after"`
}

// authSendSmsCode godoc
// @Tags auth
// @Summary authSendSmsCode
// @Description send sms code
// @ID send-sms-code
// @Accept json
// @Produce json
// @Param req body SendSmsCodeReq false "Device"
// @Success 200 {object} SendSmsCodeRes
// @Router /api/auth/sendSms [post]
// @description Demo user phone is +70001112233
func authSendSmsCode(c echo.Context) (err error) {
	var (
		req SendSmsCodeReq
		ctx = common.NewContext(c)
	)

	err = common.BindAndValidateReq(c, &req)
	if err != nil {
		return
	}

	if req.Phone == config.C.DemoUser.Phone {
		return c.JSON(http.StatusOK, SendSmsCodeRes{RefreshAfter: 30})
	}

	timeToWait, err := auth.SendSms(ctx, req.Phone)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, SendSmsCodeRes{RefreshAfter: timeToWait})
}

func authSendSmsSecure(c echo.Context) (err error) {
	var (
		req struct {
			Token string `json:"recaptcha_token" validate:"required"`
			Phone string `json:"phone" validate:"required" example:"+7999887766"`
		}
		ctx = common.NewContext(c)
	)

	err = common.BindAndValidateReq(c, &req)
	if err != nil {
		return
	}

	err = ratelimit.Apply(ctx, ratelimit.MethodSendSms)
	if err != nil {
		return
	}

	res, err := device.SendSms(ctx, req.Token, req.Phone)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, res)
}

func authRegenerateCodeSecure(c echo.Context) (err error) {
	var (
		req struct {
			UID string `json:"uid" validate:"required"`
		}
		ctx = common.NewContext(c)
	)

	err = common.BindAndValidateReq(c, &req)
	if err != nil {
		return
	}

	res, err := device.RegenerateCode(ctx, req.UID)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, res)
}

type AuthorizeReq struct {
	Vuid string `json:"vuid" validate:"required"`
	Code string `json:"code" validate:"required"`
}

// authAuthorize godoc
// @Tags auth
// @Summary authAuthorize
// @Description authorize
// @ID authorize
// @Accept json
// @Produce json
// @Param req body AuthorizeReq false "Device"
// @Success 200 {object} user.AuthorizeAndCompleteRegistrationUsecaseResponse
// @Router /api/auth/authorize [post]
func authAuthorize(c echo.Context) (err error) {
	var (
		req AuthorizeReq
		ctx = common.NewContext(c)
	)

	err = common.BindAndValidateReq(c, &req)
	if err != nil {
		return
	}

	res, err := device.Authorize(ctx, req.Vuid, req.Code)
	if err != nil {
		return
	}

	//carrot
	client := c.Request().Header.Get("device-type")

	if client == constant.WEB || client == constant.WEBMOBILE {
		cquid := c.Request().Header.Get("carrotquest-uid")
		logger.Log.Infof("carrot quest uid: %s", cquid)

		operations := analytics.UpdatePhoneUserPropsCQ(req.Vuid, ctx.UserID)

		//отправка аналитики на carrot-quest
		if cquid != "" {
			cqPropsRequest := analytics.CarrotPropsRequest{
				ID:         cquid,
				Operations: operations,
				ByUserID:   false,
			}
			go func() {
				analytics.UpdateUserPropsCQ(ctx, cqPropsRequest)
				analytics.SendEventToCQ(ctx, "Авторизация", nil)
			}()
		} else {
			logger.Log.Infof("carrot quest uid empty, no info to send")
		}
	}

	return c.JSON(http.StatusOK, res)
}

type CompleteRegistrationReq struct {
	RequestID string `json:"request_id" validate:"required"`
	Name      string `json:"name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" validate:"required"`
}

func authCompleteRegistration(c echo.Context) (err error) {
	var (
		req CompleteRegistrationReq
		ctx = common.NewContext(c)
	)
	err = common.BindAndValidateReq(c, &req)
	if err != nil {
		return
	}

	res, err := device.CompleteRegistration(ctx, req.RequestID, req.Email, req.Name, req.LastName)
	if err != nil {
		return
	}

	client := c.Request().Header.Get("device-type")

	if client == constant.WEB || client == constant.WEBMOBILE {
		cquid := c.Request().Header.Get("carrotquest-uid")
		logger.Log.Infof("carrot quest uid: %s", cquid)

		operations := analytics.MakeNewUserPropsCQ(req.Name, req.Email, ctx.UserID)

		//отправка аналитики на carrot-quest
		if cquid != "" {
			cqPropsRequest := analytics.CarrotPropsRequest{
				ID:         cquid,
				Operations: operations,
				ByUserID:   false,
			}
			go func() {
				analytics.UpdateUserPropsCQ(ctx, cqPropsRequest)
				analytics.SendEventToCQ(ctx, "Регистрация завершена", nil)
			}()
		} else {
			logger.Log.Infof("carrot quest uid empty, no info to send")
		}
	}

	return c.JSON(http.StatusOK, res)
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token"`
	DeviceID     string `json:"user_id"`
}

func authRefreshToken(c echo.Context) (err error) {
	var (
		req RefreshTokenReq
		ctx = common.NewContext(c)
	)

	err = common.BindAndValidateReq(c, &req)
	if err != nil {
		return
	}

	var res *device.AuthorizeAndCompleteRegistrationResponse

	if req.DeviceID != "" && (req.RefreshToken == "" || req.RefreshToken == "device_token") {
		res, err = device.RefreshDeviceToken(ctx, req.DeviceID)
	} else {
		res, err = device.RefreshToken(ctx, req.RefreshToken)
	}

	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, res)
}

func authLogout(c echo.Context) (err error) {
	ctx := common.NewContext(c)

	token, err := user.Logout(ctx)

	return c.JSON(http.StatusOK, CreateTemporaryUserRes{Token: token})
}

func authProxyUserInfo(c echo.Context) (err error) {
	ctx := common.NewContext(c)

	userInfo, err := auth.GetMainUserInfo(ctx)
	if err != nil {
		return
	}

	var resp struct {
		Email     string `json:"email,omitempty"`
		Id        int    `json:"id,omitempty"`
		Phone     string `json:"phone,omitempty"`
		Username  string `json:"username,omitempty"`
		FirstName string `json:"first_name,omitempty"`
		LastName  string `json:"last_name,omitempty"`
	}

	if userInfo.Email != "" {
		resp.Email = userInfo.Email
		resp.Username = strings.Split(userInfo.Email, "@")[0]
		resp.FirstName = resp.Username
	}

	if userInfo.Phone != "" {
		resp.Phone = userInfo.Phone
	}

	resp.Id = userInfo.Id

	return c.JSON(http.StatusOK, resp)
}
