package analytics

import (
	"net/url"

	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/config"
)

func SendGoogleAnalytics(ctx common.Context, query string) {
	// TODO: set content-type
	req := aclient.Clone().AppendHeader("User-Agent", "Mozilla/5.0")
	u, _ := url.Parse(config.C.Analytics.Urls.Google)
	u.RawQuery = query
	common.SendRequest(ctx, req.Get(u.String()))
	u, _ = url.Parse(config.C.Analytics.Urls.Google2)
	u.RawQuery = query
	common.SendRequest(ctx, req.Get(u.String()))
}
