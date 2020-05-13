package cmd

import (
	"bitbucket.movista.ru/maas/maasapi/internal/ratelimit"
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"bitbucket.movista.ru/maas/maasapi/internal/middlewares"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4/middleware"

	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	apiController "bitbucket.movista.ru/maas/maasapi/internal/controller"
	dialogController "bitbucket.movista.ru/maas/maasapi/internal/dialog"
	"bitbucket.movista.ru/maas/maasapi/internal/logger"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"bitbucket.movista.ru/maas/maasapi/internal/transport/metro"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	echoSwagger "github.com/swaggo/echo-swagger"
)

var apiCmd = &cobra.Command{
	Use: "api",
	Run: func(cmd *cobra.Command, args []string) {
		if err := startServer(); err != nil {
			logger.Log.Fatal(err)
		}
	},
}

func startServer() (err error) {
	defer common.CloseDb()
	defer common.CloseRedis()

	m, err := migrate.New("file://./migrations", config.C.Database.URL)
	if err != nil {
		logger.Log.Fatalf("%v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		logger.Log.Fatalf("%v", err)
	}

	metro.Init()

	e := echo.New()

	common.RegisterValidator(e)
	common.RegisterPrometheus(e)

	// Разрешение доступа к сервису с веб-браузера
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:  config.C.AllowOrigins,
		AllowMethods:  []string{echo.OPTIONS, echo.HEAD, echo.GET, echo.POST, echo.PUT},
		ExposeHeaders: []string{"Access-Control-Allow-Origin"},
	}))

	e.Use(middleware.Recover())
	e.Use(middlewares.ResponseHijacker())
	e.Use(middlewares.JsonRequestLogger())
	e.Use(middlewares.Recover())
	e.Use(middleware.RequestID())

	api := e.Group("/api")

	dialogController.Add(api)
	apiController.Add(api)

	e.HTTPErrorHandler = customErrorEcho
	e.HideBanner = true

	if !config.IsProd() {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	// Запись заголовков и тела запроса в context
	e.Pre(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var bodyBytes []byte
			if c.Request().Body != nil {
				bodyBytes, _ = ioutil.ReadAll(c.Request().Body)
			}

			// Restore the io.ReadCloser to its original state
			c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			c.Set("reqBody", bodyBytes)
			c.Set("reqHeaders", c.Request().Header)
			c.Set("requestId", uuid.New().String())
			c.Set("reqStart", time.Now())
			return next(c)
		}
	})

	logger.Log.Debug("Service RUN on DEBUG mode")

	go ratelimit.UpdateBlackList()

	// Start server
	go func() {
		if err := e.Start(config.C.Server.URL); err != nil {
			logger.Log.Info("shutting down the server", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		logger.Log.Fatal(err)
	}

	return
}

func customErrorEcho(err error, c echo.Context) {

	var (
		uapiErr    models.UapiError
		e          models.Error
		statusCode = http.StatusOK
	)

	var echoError *echo.HTTPError
	if errors.As(err, &uapiErr) {
		err = c.JSON(statusCode, echo.Map{
			"errors": uapiErr,
		})
		return
	} else if errors.As(err, &e) {
		if e.StatusCode != 0 {
			statusCode = e.StatusCode
		}
	} else if errors.As(err, &echoError) {
		if echoError.Code != 0 {
			statusCode = echoError.Code
		}
		e.Message = echoError.Message.(string)
		e.Code = strconv.Itoa(echoError.Code)
	} else {
		e = models.Error{
			Code:    "E_INTERNAL_ERROR",
			Message: "Внутренняя ошибка",
			IsError: true,
		}
	}

	e.Type = models.SystemError

	err = c.JSON(statusCode, echo.Map{
		"error":  e,
		"errors": []interface{}{e},
	})
}
