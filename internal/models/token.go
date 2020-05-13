package models

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"bitbucket.movista.ru/maas/maasapi/internal/logger"
	"github.com/coreos/go-oidc"
	"github.com/dgrijalva/jwt-go"
)

var (
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
)

func init() {
	var err error
	provider, err = oidc.NewProvider(context.Background(), config.C.Auth.Core.BaseURL)
	if err != nil {
		logger.Log.Fatal(err)
	}

	verifier = provider.Verifier(&oidc.Config{SkipClientIDCheck: true})
}

type Token struct {
	Value     string
	Content   map[string]interface{}
	ExpiresAt Time
}

type DeviceClaims struct {
	DeviceID string `json:"device_id"`
	Demo     bool   `json:"demo"`
	jwt.StandardClaims
}

func TryGetToken(auth, authMaas string) interface{} {
	var tkn interface{}
	token, err := GetHeaderToken(auth, authMaas)
	if err != nil {
		tkn = fmt.Sprintf("%v", err)
	} else {
		tkn = token
	}
	return tkn
}

func GetHeaderToken(auth, authMaas string) (*Token, error) {
	commonReq := auth
	webReq := authMaas // совместимость для web, так как frontapi ещё жив
	if commonReq == "" && webReq == "" {
		return nil, fmt.Errorf("no token")
	}

	var ss string
	if webReq != "" {
		ss = webReq
	} else {
		ss = commonReq
	}

	token, err := NewTokenFromString(ss)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func NewDeviceToken(claims DeviceClaims) *Token {
	eat := time.Now().AddDate(5, 0, 0).Unix()
	if claims.Demo {
		eat = time.Now().Add(config.C.DemoUser.TokenTTL).Unix()
	}

	claims.StandardClaims = jwt.StandardClaims{
		ExpiresAt: eat,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, _ := token.SignedString([]byte(config.C.Auth.TemporarySecret))

	tmp := make(map[string]interface{})
	mClaims, _ := json.Marshal(claims)
	_ = json.Unmarshal(mClaims, &tmp)

	return &Token{
		Value:     ss,
		Content:   tmp,
		ExpiresAt: Time{time.Unix(claims.ExpiresAt, 0)},
	}
}

func NewTokenFromString(ss string) (t *Token, err error) {

	ss = strings.ReplaceAll(ss, "Bearer ", "")

	token, _ := jwt.Parse(ss, nil)

	if token == nil || token.Claims == nil {
		err = NewUnauthorizedError()
		return
	}

	claims := token.Claims.(jwt.MapClaims)

	if _, ok := claims["iss"]; ok {
		_, err = verifier.Verify(context.Background(), ss)
		if err != nil {
			return
		}
	} else {
		_, err = jwt.Parse(ss, func(t *jwt.Token) (i interface{}, err error) {
			return []byte(config.C.Auth.TemporarySecret), nil
		})
		if err != nil {
			return
		}
	}

	t = &Token{
		Value:   ss,
		Content: claims,
	}

	if exp, ok := claims["exp"]; ok && exp != "" {
		t.ExpiresAt = Time{time.Unix(int64(exp.(float64)), 0)}
	}

	return
}

func (t *Token) String() string {
	if t.Value == "" {
		return ""
	}

	return fmt.Sprintf("Bearer %s", t.Value)
}

func (t *Token) GetDeviceID() string {
	if v, ok := t.Content["deviceId"]; ok {
		return v.(string)
	}

	if v, ok := t.Content["device_id"]; ok {
		return v.(string)
	}

	return ""
}

func (t *Token) GetUserID() int {
	if v, ok := t.Content["id"]; ok {
		sValue := v.(string)
		id, _ := strconv.Atoi(sValue)
		return id
	}

	return 0
}
