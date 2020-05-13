package common

import (
	"bitbucket.movista.ru/maas/maasapi/internal/logger"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"context"
	"github.com/go-pg/pg/v9"
	"github.com/go-redis/redis/v7"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type Context struct {
	Context     context.Context
	DB          *pg.DB
	RedisClient *redis.Client
	RequestID   string
	Token       *models.Token
	UserID      int
	DeviceID    string
	DeviceType  string
	RealIP      string
	Logger      *logrus.Entry
}

func (c *Context) SetToken(t *models.Token) {
	c.Token = t
	c.UserID = c.Token.GetUserID()
}

func (c *Context) IsUser() bool {
	if c.UserID > 0 {
		return true
	}

	return false
}

func NewInternalContext() (ctx Context) {
	ctx.RequestID = "internal"
	ctx.DeviceID = "internal"
	ctx.DB = GetDatabaseConnection()
	ctx.RedisClient = GetRedisConnection()
	ctx.Logger = logger.Log.WithFields(logger.Fields{
		"internal-request": true,
		"request-id":       ctx.RequestID,
		"device-id":        ctx.DeviceID,
	})
	return
}

func NewContext(c echo.Context) (ctx Context) {
	defer func() {
		fields := logger.Fields{
			"request-id": ctx.RequestID,
			"real-ip":    ctx.RealIP,
		}

		if ctx.DeviceID != "" {
			fields["device-id"] = ctx.DeviceID
		}

		if ctx.UserID != 0 {
			fields["user-id"] = ctx.UserID
		}

		if ctx.DeviceType != "" {
			fields["device-type"] = ctx.DeviceType
		}

		ctx.Logger = logger.Log.WithFields(fields)
	}()

	ctx.DB = GetDatabaseConnection()
	ctx.RedisClient = GetRedisConnection()
	ctx.RequestID = c.Get("requestId").(string)
	ctx.RealIP = c.RealIP()

	token := c.Get("token")
	if token == nil {
		return
	}

	ctx.Token = token.(*models.Token)
	if deviceID := ctx.Token.GetDeviceID(); deviceID != "" {
		ctx.DeviceID = deviceID
	}

	// если в middleware подменили на верный deviceID
	if deviceID := c.Get("deviceID"); deviceID != nil {
		ctx.DeviceID = deviceID.(string)
	}

	if userID := ctx.Token.GetUserID(); userID != 0 {
		ctx.UserID = userID
	}

	dType := c.Get("deviceType")
	if dType != nil {
		ctx.DeviceType = dType.(string)
	}

	return
}
