package common

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"bitbucket.movista.ru/maas/maasapi/internal/logger"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
)

var (
	tlsConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	DefaultRequest *gorequest.SuperAgent
	InternalToken  *models.Token
)

func init() {
	DefaultRequest = gorequest.New().TLSClientConfig(&tls.Config{
		InsecureSkipVerify: true,
	})

	if err := fetchTokenForInternalRequests(); err != nil {
		panic(err)
	}
}

func GetInternalToken() string {
	if time.Now().Add(time.Second * 10).After(InternalToken.ExpiresAt.Time) {
		if err := fetchTokenForInternalRequests(); err != nil {
			panic(err)
		}
	}

	return InternalToken.String()
}

func GetLegacyTokenForWebDevice(ctx Context) (token string, expiresIn int64, err error) {
	var resp struct {
		Token     string `json:"access_token"`
		TokenType string `json:"token_type"`
		ExpiresIn int64  `json:"expires_in"`
	}

	req := fmt.Sprintf("grant_type=%s&scope=%s&client_id=%s&client_secret=%s",
		config.C.Auth.Credentials.Internal.GrantType,
		config.C.Auth.Credentials.Internal.Scope,
		config.C.Auth.Credentials.ClientID,
		config.C.Auth.Credentials.ClientSecret,
	)

	spanReq := DefaultRequest.Clone().
		Post(config.C.Auth.Core.Urls.ConnectToken).
		Type("urlencoded").
		Send(req)

	body, _, err := SendRequest(ctx, spanReq)
	err = json.Unmarshal(body, &resp)
	if err != nil {
		err = errors.Wrap(err, "cannot unmarshal auth response")
		return
	}

	tokenModel, _ := models.NewTokenFromString(resp.Token)

	token = tokenModel.String()
	expiresIn = resp.ExpiresIn
	return
}

func fetchTokenForInternalRequests() (err error) {
	var resp struct {
		Token     string `json:"access_token"`
		ExpiresIn int    `json:"expires_in"`
		TokenType string `json:"token_type"`
	}

	ctx := NewInternalContext()

	req := fmt.Sprintf("grant_type=%s&scope=%s&client_id=%s&client_secret=%s",
		config.C.Auth.Credentials.Internal.GrantType,
		config.C.Auth.Credentials.Internal.Scope,
		config.C.Auth.Credentials.ClientID,
		config.C.Auth.Credentials.ClientSecret,
	)

	spanReq := DefaultRequest.Clone().
		Post(config.C.Auth.Core.Urls.ConnectToken).
		Type("urlencoded").
		Send(req)

	body, _, err := SendRequest(ctx, spanReq)
	if err != nil {
		err = fmt.Errorf("cannot make request to auth - %w", err)
		return
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		err = fmt.Errorf("cannot unmarshal auth response - %s", string(body))
		return
	}

	InternalToken, _ = models.NewTokenFromString(resp.TokenType + " " + resp.Token)

	return
}

func BindAndValidateReq(c echo.Context, req interface{}) (e error) {
	err := c.Bind(req)
	if err != nil {
		return models.NewBindOrValidateError(err)
	}

	err = c.Validate(req)
	if err != nil {
		return models.NewBindOrValidateError(err)
	}

	return
}

func SendRequest(ctx Context, goreq *gorequest.SuperAgent) (body []byte, response *http.Response, err error) {

	var (
		spanReqID = uuid.New().String()
		jl        = logger.NewJsonLogger()
		req       *http.Request
	)

	goreq = goreq.AppendHeader("Request-Id", spanReqID)

	req, err = goreq.MakeRequest()
	if err != nil {
		return
	}

	token := models.TryGetToken(req.Header.Get("Authorization"), req.Header.Get("AuthorizationMaas"))

	err = jl.LogSubrequest(req, ctx.RequestID, spanReqID, token)
	if err != nil {
		return
	}

	startTime := time.Now()

	resp, body, errs := goreq.EndBytes()
	if len(errs) > 0 {
		jl.LogSubrequestError(req, ctx.RequestID, spanReqID, startTime, errs)
		err = errs[0]
		return
	}

	response = resp

	jl.LogSubrequestResponse(req, resp, ctx.RequestID, spanReqID, startTime, body)

	if resp.StatusCode == http.StatusUnauthorized {
		err = models.NewUnauthorizedError()
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("response status code is %d, response - %s", resp.StatusCode, string(body))
		return
	}

	return
}

func UapiRequest(ctx Context, client *gorequest.SuperAgent,
	in interface{}, method string, url string) (*models.RawResponse, error) {

	var (
		err error
		req *gorequest.SuperAgent
	)

	body := make([]byte, 0)

	if in != nil {
		switch in.(type) {
		case io.Reader:
			body, err = ioutil.ReadAll(in.(io.Reader))
			if err != nil {
				return nil, err
			}
		default:
			jsonBuf := new(bytes.Buffer)
			enc := json.NewEncoder(jsonBuf)
			enc.SetEscapeHTML(false)
			enc.Encode(in)
			body = jsonBuf.Bytes()
		}
	}

	req = client.Clone()
	if method == http.MethodPost {
		req = req.Post(url)
	} else if method == http.MethodGet {
		req = req.Get(url)
	}

	if len(body) > 0 {
		req = req.SendStruct(json.RawMessage(body))
	}

	rawResp := &models.RawResponse{}
	resp, _, err := SendRequest(ctx, req)
	// обработка ситуации, когда status code != 200 и есть массив errors
	// пример: ответ от personalarea на создание документа
	if err != nil && resp != nil && len(resp) > 0 {
		// status code != 200 и что-то пришло в body
		if err := json.Unmarshal(resp, rawResp); err == nil {
			// это что-то в формате UapiResponse
			if rawResp.Data == nil && rawResp.Error != nil && len(rawResp.Error) > 0 {
				// и в нем непустой массив Error с пустым Data
				return nil, rawResp.Error
			}
		}
		return nil, err
	} else if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, rawResp)
	if err != nil {
		return nil, err
	}

	if rawResp.Data != nil {
		return rawResp, nil
	}

	if len(rawResp.Error) != 0 {
		return nil, rawResp.Error
	}

	return rawResp, nil
}

func UapiAuthorizedRequest(ctx Context, client *gorequest.SuperAgent, in interface{},
	method string, url string, token interface{}) (*models.RawResponse, error) {

	var tkn string

	switch token.(type) {
	case string:
		tkn = token.(string)
	default:
		tkn = ctx.Token.String()
	}

	return UapiRequest(ctx, client.Clone().AppendHeader("Authorization", tkn), in, method, url)
}

func UapiAuthorizedPost(ctx Context, client *gorequest.SuperAgent, in interface{}, url string) (*models.RawResponse, error) {

	return UapiAuthorizedRequest(ctx, client, in, http.MethodPost, url, ctx.Token.String())
}
func UapiAuthorizedGet(ctx Context, client *gorequest.SuperAgent, url string) (*models.RawResponse, error) {

	return UapiAuthorizedRequest(ctx, client, nil, http.MethodGet, url, ctx.Token.String())
}

func UapiGet(ctx Context, client *gorequest.SuperAgent, url string) (*models.RawResponse, error) {
	return UapiRequest(ctx, client, nil, http.MethodGet, url)
}

func UapiPost(ctx Context, client *gorequest.SuperAgent, in interface{}, url string) (*models.RawResponse, error) {
	return UapiRequest(ctx, client, in, http.MethodPost, url)
}
