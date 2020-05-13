package device

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"bitbucket.movista.ru/maas/maasapi/internal/favorite"
	"bitbucket.movista.ru/maas/maasapi/internal/logger"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"bitbucket.movista.ru/maas/maasapi/internal/notebook"
	"bitbucket.movista.ru/maas/maasapi/internal/session"
	"bitbucket.movista.ru/maas/maasapi/internal/uapi/auth"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

type CreateTemporaryUserUsecaseResponse struct {
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIn   int64  `json:"expires_in,omitempty"`
}

type AuthorizeAndCompleteRegistrationResponse struct {
	AccessToken   string      `json:"access_token,omitempty"`
	RefreshToken  string      `json:"refresh_token,omitempty"`
	AuthUserID    string      `json:"user_id,omitempty"`
	HasEmail      bool        `json:"has_email"`
	RequestID     string      `json:"request_id,omitempty"`
	ExpiresIn     int64       `json:"expires_in,omitempty"`
	UserProfile   interface{} `json:"user_profile,omitempty"`
	UserDocuments interface{} `json:"user_documents,omitempty"`
	CarrotKey     string      `json:"carrot_key,omitempty"`
}

type SendSmsResponse struct {
	RefreshAfter int    `json:"refresh_after"`
	UID          string `json:"uid"`
}

func Authorize(ctx common.Context, phone, code string) (res *AuthorizeAndCompleteRegistrationResponse, err error) {
	// DEMO user flow
	if phone == config.C.DemoUser.Phone && code == config.C.DemoUser.Code {
		token := models.NewDeviceToken(models.DeviceClaims{DeviceID: ctx.DeviceID, Demo: true})
		res = &AuthorizeAndCompleteRegistrationResponse{
			AccessToken:  token.String(),
			RefreshToken: "emptyDemoUserRefreshToken",
			AuthUserID:   "0",
			HasEmail:     true,
		}
		return
	}

	aResp, err := auth.Authorize(ctx, phone, code, map[string]interface{}{"deviceId": ctx.DeviceID})
	if err != nil {
		return
	}

	if aResp.AuthenticationRequestId != "" {
		res = &AuthorizeAndCompleteRegistrationResponse{
			RequestID: aResp.AuthenticationRequestId,
			HasEmail:  false,
		}
		return
	}

	t, _ := models.NewTokenFromString(aResp.AccessToken)

	ctx.SetToken(t)

	err = session.Create(ctx, ctx.Token.GetUserID())
	if err != nil {
		return
	}

	userProfile, err := GetOneWithSettingsAndFavorites(ctx)
	if err != nil {
		return
	}

	userDocs, err := notebook.GetDocuments(ctx)
	if err != nil {
		logger.Log.Error(err)
	}

	// генерация HMAC для carrotQuest
	cqsha := getCarrotHMACsha(ctx.Token.GetUserID())

	res = &AuthorizeAndCompleteRegistrationResponse{
		AccessToken:   aResp.AccessToken,
		RefreshToken:  aResp.RefreshToken,
		AuthUserID:    fmt.Sprintf("%d", ctx.Token.GetUserID()),
		HasEmail:      true,
		ExpiresIn:     aResp.ExpiresIn,
		UserProfile:   userProfile,
		UserDocuments: userDocs,
		CarrotKey:     cqsha,
	}

	return
}

func CompleteRegistration(ctx common.Context, requestId, email, firstName, lastName string) (res *AuthorizeAndCompleteRegistrationResponse, err error) {

	crResp, err := auth.CompleteRegistration(ctx, requestId, email, firstName, lastName, map[string]interface{}{"deviceId": ctx.DeviceID})
	if err != nil {
		return
	}

	if firstName == "" {
		firstName = strings.Split(email, "@")[0]
	}

	t, _ := models.NewTokenFromString(crResp.AccessToken)

	ctx.SetToken(t)

	err = session.Create(ctx, ctx.Token.GetUserID())
	if err != nil {
		return
	}

	err = favorite.CopyDeviceFavoritesToUser(ctx)
	if err != nil {
		return
	}

	userProfile, err := GetOneWithSettingsAndFavorites(ctx)
	if err != nil {
		return
	}

	userDocs, err := notebook.GetDocuments(ctx)
	if err != nil {
		logger.Log.Error(err)
	}

	// генерация HMAC для carrotQuest
	cqsha := getCarrotHMACsha(ctx.Token.GetUserID())

	res = &AuthorizeAndCompleteRegistrationResponse{
		AccessToken:   crResp.AccessToken,
		RefreshToken:  crResp.RefreshToken,
		AuthUserID:    string(ctx.Token.GetUserID()),
		HasEmail:      true,
		ExpiresIn:     crResp.ExpiresIn,
		UserProfile:   userProfile,
		UserDocuments: userDocs,
		CarrotKey:     cqsha,
	}

	return
}

func RefreshDeviceToken(ctx common.Context, correctDeviceID string) (aur *AuthorizeAndCompleteRegistrationResponse, err error) {
	ctx.DeviceID = correctDeviceID
	_, err = GetOne(ctx)
	if err != nil {
		return
	}

	token := models.NewDeviceToken(models.DeviceClaims{DeviceID: correctDeviceID})

	aur = &AuthorizeAndCompleteRegistrationResponse{
		AccessToken:  token.String(),
		RefreshToken: "device_token",
	}

	return
}

func RefreshToken(ctx common.Context, refreshToken string) (aur *AuthorizeAndCompleteRegistrationResponse, err error) {
	res, err := auth.RefreshToken(ctx, refreshToken)
	if err != nil {
		return
	}

	token, _ := models.NewTokenFromString(res.AccessToken)
	oldDeviceField := token.Content["device"]
	if token.GetDeviceID() == "" && oldDeviceField == nil {
		err = models.NewEmptyDeviceFieldInUserTokenError()
		return
	}

	aur = &AuthorizeAndCompleteRegistrationResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		HasEmail:     true,
		ExpiresIn:    res.ExpiresIn,
	}

	return
}

func SendSms(ctx common.Context, recaptchaToken, phone string) (resp SendSmsResponse, err error) {
	if phone == config.C.DemoUser.Phone {
		var uid string
		uid, err = createRegenerateCodeUID(ctx, phone, 30)
		if err != nil {
			return
		}

		resp = SendSmsResponse{
			RefreshAfter: 30,
			UID:          uid,
		}

		return
	}

	if err = common.CheckRecaptcha(ctx, recaptchaToken); err != nil {
		return
	}

	timeToWait, err := auth.SendSms(ctx, phone)
	if err != nil {
		return
	}

	uid, err := createRegenerateCodeUID(ctx, phone, timeToWait)
	if err != nil {
		return
	}

	resp = SendSmsResponse{
		RefreshAfter: timeToWait,
		UID:          uid,
	}

	return
}

func RegenerateCode(ctx common.Context, uid string) (resp SendSmsResponse, err error) {
	rga := RegenerateCodeAttempt{ID: uid}

	err = ctx.DB.Model(&rga).WherePK().Select()
	if err != nil {
		err = models.NewRegenerateAttemptWrongUIDError()
		return
	}

	if rga.DeviceID != ctx.DeviceID {
		err = models.NewRegenerateAttemptDifferentDeviceError()
		return
	}

	if time.Now().Before(rga.NotBefore.Time) {
		err = models.NewRegenerateAttemptBeforeNBFError(rga.NotBefore)
		return
	}

	timeToWait, err := auth.SendSms(ctx, rga.Phone)
	if err != nil {
		return
	}

	nextUID, err := createRegenerateCodeUID(ctx, rga.Phone, timeToWait)
	if err != nil {
		return
	}

	resp = SendSmsResponse{
		RefreshAfter: timeToWait,
		UID:          nextUID,
	}

	_, err = ctx.DB.Model(&rga).WherePK().Delete()

	return
}

func createRegenerateCodeUID(ctx common.Context, phone string, timeToWait int) (uid string, err error) {
	rga := RegenerateCodeAttempt{
		DeviceID:  ctx.DeviceID,
		Phone:     phone,
		NotBefore: models.Time{Time: time.Now().Add(time.Duration(timeToWait) * time.Second)},
	}
	_, err = ctx.DB.Model(&rga).Insert()
	if err != nil {
		return
	}

	uid = rga.ID

	return
}

func getCarrotHMACsha(userID int) (cqsha string) {
	h := hmac.New(sha256.New, []byte(string(userID)))

	h.Write([]byte(config.C.CarrotQuestAuthKey))

	cqsha = hex.EncodeToString(h.Sum(nil))

	return
}
