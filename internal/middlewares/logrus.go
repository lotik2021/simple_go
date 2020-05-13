package middlewares

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// LogrusLogger is echo middleware that logs requests using logger "logrus"
func LogrusLogger(l *logrus.Entry) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			reqID := c.Get("requestId")

			l.WithField("method", c.Request().Method).
				WithField("uri", c.Request().RequestURI).
				WithField("request-id", reqID).
				Infof("Started processing request")

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			req := c.Request()
			res := c.Response()

			fields := logrus.Fields{
				"status":     res.Status,
				"latency":    time.Since(start),
				"request-id": reqID,
				"method":     req.Method,
				"uri":        req.RequestURI,
			}

			n := res.Status
			switch {
			case n >= 500:
				l.WithFields(fields).Error("server_error")
			case n >= 400:
				l.WithFields(fields).Error("client_error")
			default:
				l.WithFields(fields).Info("Ended processing request")
			}

			return nil
		}
	}
}
