package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"bitbucket.movista.ru/maas/maasapi/internal/constant"
)

const (
	SystemError = "system"
)

type Error struct {
	StatusCode int    `json:"-"`
	IsError    bool   `json:"isError"`
	Code       string `json:"code"`
	Message    string `json:"message,omitempty"`
	Type       string `json:"type"`
}

func (e Error) Error() string {
	if e.Message == "" {
		return e.Code
	}
	return fmt.Sprintf("%s - %s", e.Code, e.Message)
}

func (e Error) MarshallToJson() []byte {
	res, _ := json.Marshal(e)
	return res
}

func NewUnauthorizedError() error {
	return Error{
		StatusCode: http.StatusUnauthorized,
		Code:       "401",
		Message:    "Unauthorized request",
		Type:       SystemError,
		IsError:    true,
	}
}

func NewRegenerateAttemptWrongUIDError() error {
	return Error{
		Code:    "E_WRONG_REGENERATE_UID",
		Message: "Неверный или использованный uid",
		Type:    SystemError,
	}
}

func NewRegenerateAttemptDifferentDeviceError() error {
	return Error{
		Code:    "E_WRONG_REGENERATE_DEVICE",
		Message: "Устройство не совпадает с тем, с которого был первый запрос на отправку смс",
		Type:    SystemError,
	}
}

func NewRegenerateAttemptBeforeNBFError(nbf Time) error {
	return Error{
		Code:    "E_WRONG_REGENERATE_TIME",
		Message: fmt.Sprintf("Пока нельзя использовать uid, можно через %s", nbf.Time.Sub(time.Now())),
		Type:    SystemError,
	}
}

func NewBindOrValidateError(err error) error {
	return Error{
		StatusCode: http.StatusOK,
		Code:       constant.ErrorInvalidRequest,
		Message:    err.Error(),
		Type:       SystemError,
	}
}

func NewEmptyDeviceFieldInUserTokenError() error {
	return Error{
		StatusCode: http.StatusOK,
		Code:       "EMPTY_DEVICE_FIELD_IN_TOKEN",
		Message:    "Пустое поле device в токене",
		Type:       SystemError,
	}
}

func NewInvalidRecaptchaError() error {
	return Error{
		StatusCode: http.StatusOK,
		Code:       constant.ErrorInvalidCaptcha,
		Message:    "Ошибка валидации reCaptcha",
		Type:       SystemError,
	}
}

func NewFindPlaceError() error {
	return Error{
		IsError:    true,
		StatusCode: http.StatusOK,
		Code:       constant.ErrorInvalidIPRequest,
		Message:    "Место не найдено",
	}
}

func NewRateLimitError() error {
	return Error{
		IsError:    true,
		StatusCode: http.StatusTooManyRequests,
		Code:       constant.ErrorTooManyRequests,
		Message:    "Превышен лимит запросов",
		Type:       SystemError,
	}
}

func NewInternalDialogError(err error) error {
	e := new(Error)
	resultError := new(Error)
	if errors.As(err, e) {
		resultError = e

	} else {
		resultError.Code = "INTERNAL_ERROR"
		resultError.Message = err.Error()
	}
	resultError.StatusCode = http.StatusInternalServerError
	resultError.Type = SystemError
	return resultError
}
