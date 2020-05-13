package middlewares

import (
	"net"
	"net/http"
	"strings"

	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"bitbucket.movista.ru/maas/maasapi/internal/constant"
	"bitbucket.movista.ru/maas/maasapi/internal/device"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"github.com/labstack/echo/v4"
	ua "github.com/mileusna/useragent"
)

func Token(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		if c.Request().Method == http.MethodOptions {
			return c.NoContent(http.StatusOK)
		}

		if strings.Contains(c.Path(), "/metrics") {
			// для dev, test, stage, local
			if !config.IsProd() {
				return next(c)
			}

			// для prod
			ip := net.ParseIP(c.RealIP())
			_, cidr, _ := net.ParseCIDR("10.0.0.0/8")
			if cidr.Contains(ip) {
				return next(c)
			}
		}

		token, err := models.GetHeaderToken(
			c.Request().Header.Get("Authorization"),
			c.Request().Header.Get("AuthorizationMaas"),
		)
		if err != nil {
			return models.NewUnauthorizedError()
		}

		c.Set("token", token)

		ctx := common.NewContext(c)

		// проверяем если старый токен авторизованного пользователя (тогда в токене будет поле id и device, и не будет одного из [deviceId, device_id])
		// => в контексте будет UserID и не будет DeviceID
		if ctx.DeviceID == "" && ctx.IsUser() {

			oldAuthDeviceID := ctx.Token.Content["device"]
			if oldAuthDeviceID == nil {
				return models.NewUnauthorizedError()
			}

			tmpContext := common.NewContext(c)
			tmpContext.DeviceID = oldAuthDeviceID.(string)
			oldDevice, err := device.GetOldOneByAuthDeviceID(tmpContext)
			if err != nil {
				return models.NewUnauthorizedError()
			}

			ctx.DeviceID = oldDevice.ID

			c.Set("deviceID", oldDevice.ID) // для дальнейшего использования этого device_id
		}

		d, _ := device.GetOne(ctx)
		if d == nil {
			// если тут, тогда в токене старый authDeviceId
			d, _ = device.GetOldOneByAuthDeviceID(ctx)
			if d == nil {
				return models.NewUnauthorizedError()
			}

			c.Set("deviceID", d.ID)
		}

		userAgent := ua.Parse(c.Request().UserAgent())
		if userAgent.Mobile || userAgent.Desktop || userAgent.Tablet {
			d.DeviceType = constant.WEB
		} else if strings.Contains(c.Request().UserAgent(), "okhttp") {
			d.DeviceType = constant.Android
		} else {
			d.DeviceType = constant.IOS
		}

		c.Set("deviceType", d.DeviceType)

		return next(c)
	}
}
