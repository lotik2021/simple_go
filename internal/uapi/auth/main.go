package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"bitbucket.movista.ru/maas/maasapi/internal/constant"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"github.com/parnurzeal/gorequest"
)

var (
	authClient *gorequest.SuperAgent
)

func init() {
	authClient = common.DefaultRequest.Clone().Timeout(config.C.Auth.RequestTimeout)
}

func SendSms(ctx common.Context, phone string) (timeWait int, err error) {
	sendSmsCodeReq := struct {
		Phone          string `json:"phone"`
		PasswordLength int    `json:"passwordLength"`
	}{
		Phone:          phone,
		PasswordLength: 6,
	}

	if ctx.DeviceType == constant.WEB {
		sendSmsCodeReq.PasswordLength = 4
	}

	var sendSmsCodeRes struct {
		TimeWait *int `json:"timeWait"`
	}

	resp, err := common.UapiPost(ctx, authClient, sendSmsCodeReq,
		config.C.Auth.Account.Urls.Smscode)
	if err != nil {
		return
	}

	err = json.Unmarshal(resp.Data, &sendSmsCodeRes)
	if err != nil {
		err = fmt.Errorf("cannot unmarshal auth response - %s", string(resp.Data))
		return
	}

	if sendSmsCodeRes.TimeWait == nil {
		err = fmt.Errorf("error response from auth - %s", string(resp.Data))
		return
	}

	timeWait = *sendSmsCodeRes.TimeWait

	return
}

func Authorize(ctx common.Context, phone, code string, extraParams map[string]interface{}) (resp *AuthorizeAndCompleteRegistrationResponse, err error) {
	authorizeReq := struct {
		ClientID    string                 `json:"clientId"`
		Secret      string                 `json:"secret"`
		Scope       string                 `json:"scope"`
		SmsCode     string                 `json:"smsCode"`
		Phone       string                 `json:"phone"`
		ExtraParams map[string]interface{} `json:"extraParams"`
	}{
		ClientID:    config.C.Auth.Credentials.ClientID,
		Secret:      config.C.Auth.Credentials.ClientSecret,
		Scope:       config.C.Auth.Credentials.External.Scope,
		SmsCode:     code,
		Phone:       phone,
		ExtraParams: extraParams,
	}

	authorizeRes := authorizeAndCompleteRegistrationApiResponse{}

	rawResp, err := common.UapiPost(ctx, authClient, authorizeReq,
		config.C.Auth.Account.Urls.Authorize)
	if err != nil {
		return
	}

	err = json.Unmarshal(rawResp.Data, &authorizeRes)
	if err != nil {
		err = fmt.Errorf("cannot unmarshal auth response - %s", string(rawResp.Data))
		return
	}

	if authorizeRes.AuthenticationRequestExpired {
		err = fmt.Errorf("AuthenticationRequestExpired - %s", string(rawResp.Data))
		return
	}

	if authorizeRes.AuthenticationRequestId != "" {
		resp = &AuthorizeAndCompleteRegistrationResponse{
			AuthenticationRequestId: authorizeRes.AuthenticationRequestId,
		}
		return
	}

	if authorizeRes.AccessToken != "" && authorizeRes.RefreshToken != "" {
		resp = &AuthorizeAndCompleteRegistrationResponse{
			AccessToken:  fmt.Sprintf("Bearer %s", authorizeRes.AccessToken),
			RefreshToken: authorizeRes.RefreshToken,
			ExpiresIn:    authorizeRes.ExpiresIn,
		}
		return
	}

	err = fmt.Errorf("invalid response(%s), request(%v)", string(rawResp.Data), authorizeReq)

	return
}

func CompleteRegistration(ctx common.Context, requestId, email, firstName, lastName string, extraParams map[string]interface{}) (resp *AuthorizeAndCompleteRegistrationResponse, err error) {
	confirmRegistrationReq := struct {
		AuthenticationRequestId string                 `json:"authenticationRequestId"`
		Email                   string                 `json:"email"`
		FirstName               string                 `json:"firstName"`
		LastName                string                 `json:"lastName"`
		ExtraParams             map[string]interface{} `json:"extraParams"`
	}{
		AuthenticationRequestId: requestId,
		Email:                   email,
		FirstName:               firstName,
		LastName:                lastName,
		ExtraParams:             extraParams,
	}

	confirmRegistrationRes := authorizeAndCompleteRegistrationApiResponse{}

	rawResp, err := common.UapiPost(ctx, authClient, confirmRegistrationReq,
		config.C.Auth.Account.Urls.CompleteRegistration)
	if err != nil {
		return
	}

	err = json.Unmarshal(rawResp.Data, &confirmRegistrationRes)
	if err != nil {
		err = fmt.Errorf("cannot unmarshal auth response - %s", string(rawResp.Data))
		return
	}

	if confirmRegistrationRes.AuthenticationRequestExpired {
		err = fmt.Errorf("AuthenticationRequestExpired - %s", string(rawResp.Data))
		return
	}

	if confirmRegistrationRes.AccessToken != "" && confirmRegistrationRes.RefreshToken != "" {
		resp = &AuthorizeAndCompleteRegistrationResponse{
			AccessToken:  fmt.Sprintf("Bearer %s", confirmRegistrationRes.AccessToken),
			RefreshToken: confirmRegistrationRes.RefreshToken,
			ExpiresIn:    confirmRegistrationRes.ExpiresIn,
		}
		return
	}

	err = fmt.Errorf("invalid response(%s), request(%v)", string(rawResp.Data), confirmRegistrationReq)

	return
}

func RefreshToken(ctx common.Context, refreshToken string) (res *AuthorizeAndCompleteRegistrationResponse, err error) {
	req := fmt.Sprintf("grant_type=%s&scope=%s&client_id=%s&client_secret=%s&refresh_token=%s",
		config.C.Auth.Credentials.RefreshToken.GrantType,
		config.C.Auth.Credentials.RefreshToken.Scope,
		config.C.Auth.Credentials.ClientID,
		config.C.Auth.Credentials.ClientSecret,
		refreshToken,
	)

	var resp struct {
		IdToken      string `json:"id_token"`
		Token        string `json:"access_token"`
		TokenType    string `json:"token_type"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int64  `json:"expires_in"`
	}

	spanReq := authClient.Clone().
		Post(config.C.Auth.Core.Urls.ConnectToken).
		Type("urlencoded").
		Send(req)

	body, _, err := common.SendRequest(ctx, spanReq)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		err = fmt.Errorf("cannot unmarshal auth response - %s", string(body))
		return
	}

	if resp.Token == "" || resp.RefreshToken == "" {
		err = fmt.Errorf("empty accessToken or refreshToken - response %+v", resp)
		return
	}

	res = &AuthorizeAndCompleteRegistrationResponse{
		AccessToken:  fmt.Sprintf("Bearer %s", resp.Token),
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    resp.ExpiresIn,
	}

	return
}

func GetMainUserInfo(ctx common.Context) (resp *UserInfo, err error) {
	spanReq := authClient.Clone().
		AppendHeader("Authorization", ctx.Token.String()).
		Get(config.C.Auth.Core.Urls.ConnectUserInfo)

	body, _, err := common.SendRequest(ctx, spanReq)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		err = fmt.Errorf("cannot unmarshal auth response - %s", string(body))
		return
	}

	return
}

func GetUserInfo(ctx common.Context) (u *User, err error) {
	resp, err := common.UapiAuthorizedRequest(
		ctx, authClient, nil, http.MethodGet,
		fmt.Sprintf("%s/%d", config.C.Auth.Account.Urls.GetUser, ctx.UserID),
		common.GetInternalToken())
	if err != nil {
		return
	}

	var res User

	err = json.Unmarshal(resp.Data, &res)
	if err != nil {
		err = fmt.Errorf("cannot unmarshal auth response - %s", string(resp.Data))
		return
	}

	if res.ID == 0 {
		err = fmt.Errorf("cannot get info about userID %d", ctx.UserID)
		return
	}

	u = &User{
		ID:        res.ID,
		FirstName: res.FirstName,
		LastName:  res.LastName,
		Email:     res.Email,
		Phone:     res.Phone,
	}

	return
}

func ConfirmEmail(ctx common.Context, req ConfirmEmailRequest) (*models.RawResponse, error) {
	u, err := url.Parse(config.C.Auth.Account.Urls.ConfirmEmail)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("code", req.Code)
	u.RawQuery = q.Encode()

	return common.UapiAuthorizedGet(ctx, authClient, u.String())
}
