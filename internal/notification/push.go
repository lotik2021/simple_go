package notification

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"fmt"
	"github.com/parnurzeal/gorequest"
)

var (
	pushClient *gorequest.SuperAgent
)

func init() {
	pushClient = common.DefaultRequest.Clone().Timeout(config.C.Push.RequestTimeout)
}

type PushNotification struct {
	PlayerID string `json:"player_id" validate:"required"`
	Message  string `json:"message" validate:"required"`
	Ttl      int    `json:"ttl,omitempty"`
}

type PushNotificationsRequest struct {
	Data []PushNotification `json:"data"`
}

func SendPushNotification(ctx common.Context, req PushNotificationsRequest) (err error) {

	spanReq := pushClient.Clone().Post(config.C.Push.Urls.Send).SendStruct(req)

	_, _, err = common.SendRequest(ctx, spanReq)
	if err != nil {
		err = fmt.Errorf("cannot make request to push-service - %w", err)
		return
	}

	return
}
