package analytics

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"bitbucket.movista.ru/maas/maasapi/internal/logger"
	"github.com/parnurzeal/gorequest"
	"strconv"
	"time"
)

var (
	aclient *gorequest.SuperAgent
)

func init() {
	aclient = common.DefaultRequest.Clone().Timeout(config.C.OpenWeather.RequestTimeout)
}

func SendEventToCQ(ctx common.Context, event string, params *CarrotParams) {
	url := config.C.Push.Urls.CarrotQuestEvent

	request := CarrotEventRequest{
		ID:       strconv.Itoa(ctx.UserID),
		Event:    event,
		Params:   params,
		Created:  time.Now().Unix(),
		ByUserID: true,
	}

	req := aclient.Clone().Post(url).SendStruct(request)

	resp, _, err := common.SendRequest(ctx, req)
	if err != nil {
		logger.Log.Errorf("cannot send event to push-service, err: %w", err)
	}
	logger.Log.Infof("push-service response:", resp)

	return
}

func UpdateUserPropsCQ(ctx common.Context, request CarrotPropsRequest) {
	url := config.C.Push.Urls.CarrotQuestProps
	req := aclient.Clone().Post(url).SendStruct(request)

	resp, _, err := common.SendRequest(ctx, req)
	if err != nil {
		logger.Log.Errorf("cannot send user props to push-service, err: %w", err)
	}

	logger.Log.Info("push-service response: %s", string(resp))
}

func UpdatePhoneUserPropsCQ(phone string, userID int) (operations []map[string]string) {

	if len(phone) >= 0 {
		phoneCQ := make(map[string]string, 3)
		phoneCQ["op"] = "update_or_create"
		phoneCQ["key"] = "$phone"
		phoneCQ["value"] = phone
		operations = append(operations, phoneCQ)
	}

	if len(strconv.Itoa(userID)) >= 0 {
		userIDCQ := make(map[string]string, 3)
		userIDCQ["op"] = "update_or_create"
		userIDCQ["key"] = "$user_id"
		userIDCQ["value"] = strconv.Itoa(userID)
		operations = append(operations, userIDCQ)
	}

	return
}

func MakeNewUserPropsCQ(name, email string, userID int) (operations []map[string]string) {

	if len(name) >= 0 {
		nameCQ := make(map[string]string, 3)
		nameCQ["op"] = "update_or_create"
		nameCQ["key"] = "$name"
		nameCQ["value"] = name
		operations = append(operations, nameCQ)
	}

	if len(email) >= 0 {
		emailCQ := make(map[string]string, 3)
		emailCQ["op"] = "update_or_create"
		emailCQ["key"] = "$email"
		emailCQ["value"] = email
		operations = append(operations, emailCQ)
	}

	if len(strconv.Itoa(userID)) >= 0 {
		userIDCQ := make(map[string]string, 3)
		userIDCQ["op"] = "update_or_create"
		userIDCQ["key"] = "$user_id"
		userIDCQ["value"] = strconv.Itoa(userID)
		operations = append(operations, userIDCQ)
	}

	return
}
