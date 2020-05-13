package middlewares

import (
	"fmt"
	"runtime"

	"bitbucket.movista.ru/maas/maasapi/internal/logger"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	rwriter "bitbucket.movista.ru/maas/maasapi/internal/response_writer"
	"github.com/labstack/echo/v4"
)

func ResponseHijacker() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Writer = rwriter.NewWriter(c.Response().Writer)
			return next(c)
		}
	}
}

func Recover() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				const DefaultStackSize = 4 << 10 // 4 KB
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					stack := make([]byte, DefaultStackSize)
					length := runtime.Stack(stack, true)

					js := logger.NewJsonLogger()
					js.LogPanic(c, err, stack[:length])

					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}

func JsonRequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := models.TryGetToken(
				c.Request().Header.Get("Authorization"),
				c.Request().Header.Get("AuthorizationMaas"),
			)
			js := logger.NewJsonLogger()
			js.LogRequest(c, token)

			err := next(c)
			if err != nil {
				c.Error(err)
			}
			js.LogResponse(c, err)
			return nil
		}
	}
}
