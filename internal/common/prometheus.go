package common

import (
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"net"
	"net/http"
	"strings"
)

var ipNet *net.IPNet

func RegisterPrometheus(e *echo.Echo) {
	_, ipNet, _ = net.ParseCIDR("10.0.0.0/8")
	p := prometheus.NewPrometheus("echo", promSkipper)
	p.Use(e)
	e.Use(promAllowOnlyInternalNetworkMiddleware())
}

func promSkipper(c echo.Context) bool {
	if strings.Contains(c.Path(), "/user/ping/location") {
		return true
	}

	if strings.Contains(c.Path(), "/user/setosplayerid") {
		return true
	}
	return false
}

func promAllowOnlyInternalNetworkMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// если пришли не на /metrics - пропускаем дальше
			if c.Path() != "/metrics" {
				return next(c)
			}

			// только для сети внутри куба
			if !ipNet.Contains(net.ParseIP(c.RealIP())) {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "метод доступен только внутри сети"})
			}

			return next(c)

		}
	}
}
