package common

import (
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"bitbucket.movista.ru/maas/maasapi/internal/constant"
	"bitbucket.movista.ru/maas/maasapi/internal/logger"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"encoding/json"
	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
	"net/url"
	"strings"
)

func CheckRecaptcha(ctx Context, token string) (err error) {
	if !config.C.Recaptcha.Enabled {
		return nil
	}

	var secretKey string
	if ctx.DeviceType == constant.IOS {
		secretKey = config.C.Recaptcha.SecretKeyIos
	} else if ctx.DeviceType == constant.Android {
		secretKey = config.C.Recaptcha.SecretKeyAndroid
	} else if ctx.DeviceType == constant.WEB {
		secretKey = config.C.Recaptcha.SecretKeyWeb
	}

	var (
		captchaReq = url.Values{
			"secret":   []string{secretKey},
			"response": []string{token},
			"remoteip": []string{ctx.RealIP},
		}
		captchaResp struct {
			Success     bool     `json:"success"`
			ChallengeTs string   `json:"challenge_ts"`
			Hostname    string   `json:"hostname"`
			ErrorCodes  []string `json:"error-codes"`
		}
	)

	spanReq := DefaultRequest.Clone().Type(gorequest.TypeUrlencoded).Post(config.C.Recaptcha.BaseURL).Send(captchaReq.Encode())

	body, _, err := SendRequest(ctx, spanReq)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &captchaResp)
	if err != nil {
		err = errors.Wrap(err, "cannot unmarshal google recaptcha response")
		return
	}

	if captchaResp.ErrorCodes != nil {
		logger.Log.Warn("google recaptcha errors: %s", strings.Join(captchaResp.ErrorCodes, ","))
	}

	if !captchaResp.Success {
		return models.NewInvalidRecaptchaError()
	}
	return
}
