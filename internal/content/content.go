package content

import (
	"encoding/json"

	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"github.com/parnurzeal/gorequest"
)

var (
	contentClient *gorequest.SuperAgent
)

func init() {
	contentClient = common.DefaultRequest.Clone().Timeout(config.C.RequestTimeout)
}

type CmsContent struct {
	Query string `json:"query"`
}

func QueryCmsV1(ctx common.Context, in CmsContent) (*models.RawResponse, error) {
	req := contentClient.Clone().AppendHeader("Authorization", ctx.Token.String()).Get(config.C.Content.Urls.QueryCms + in.Query)
	resp, _, err := common.SendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	return &models.RawResponse{
		Data: json.RawMessage(resp),
	}, nil
}
