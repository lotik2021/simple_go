package fapi

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"github.com/parnurzeal/gorequest"
)

var (
	faClient *gorequest.SuperAgent
)

func init() {
	faClient = common.DefaultRequest.Clone().Timeout(config.C.FapiAdapter.RequestTimeout)
}
