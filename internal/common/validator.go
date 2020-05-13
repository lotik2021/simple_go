package common

import (
	"bitbucket.movista.ru/maas/maasapi/internal/logger"
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
	"time"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func RegisterValidator(e *echo.Echo) {
	v := validator.New()
	if err := v.RegisterValidation("optional-correct-rfc3339-date", ValidateOptionalCorrectDate); err != nil {
		logger.Log.WithError(err).Fatal("cannot register optional-correct-rfc3339-date")
	}
	if err := v.RegisterValidation("required-correct-rfc3339-date", ValidateRequiredCorrectDate); err != nil {
		logger.Log.WithError(err).Fatal("cannot register required-correct-rfc3339-date")
	}

	e.Validator = &CustomValidator{validator: v}
}

func ValidateOptionalCorrectDate(fl validator.FieldLevel) bool {
	v := fl.Field().String()

	if v == "" {
		return true
	}

	_, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return false
	}

	return true
}

func ValidateRequiredCorrectDate(fl validator.FieldLevel) bool {
	v := fl.Field().String()

	if v == "" {
		return false
	}

	_, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return false
	}

	return true
}
