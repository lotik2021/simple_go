package device

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/google"
	"bitbucket.movista.ru/maas/maasapi/internal/logger"
	"github.com/go-pg/pg/v9"
	"googlemaps.github.io/maps"
)

func GetGoogleTransitModeSettings(ctx common.Context) (modes []maps.TransitMode) {
	err := ctx.DB.Model((*Device)(nil)).Where("id = ?", ctx.DeviceID).Column("google_transit_modes").Select(pg.Array(&modes))
	if err != nil {
		logger.Log.Error(err)
	}

	if len(modes) == 0 {
		modes = []maps.TransitMode{"bus", "subway", "tram", "train"}
	}

	return
}

type GetGoogleTransportResp struct {
	GoogleTransport string `pg:"google_transport" json:"google_transport"`
	TransportName   string `pg:"transport_name" json:"transport_name"`
	FilterStatus    bool   `pg:"filter_status" json:"filter_status"`
	IconName        string `pg:"icon_name" json:"icon_name"`
	IosIconUrl      string `pg:"ios_icon_url" json:"ios_icon_url"`
	AndroidIconUrl  string `pg:"android_icon_url" json:"android_icon_url"`
}

func GetGoogleTransitModeSettingsWithInfo(ctx common.Context) (resp []GetGoogleTransportResp, err error) {
	var res []struct {
		FilterStatus bool `pg:"fs"`
		google.TransitModeName
	}

	sql := `
		select t.*, t.name=ANY(d.google_transit_modes) as fs
		from maasapi.device d, maasapi.google_transit_mode_name t
		where d.id = ?
	`

	_, err = ctx.DB.Model().Query(&res, sql, ctx.DeviceID)

	resp = make([]GetGoogleTransportResp, 0)

	for _, v := range res {
		resp = append(resp, GetGoogleTransportResp{
			GoogleTransport: v.Name,
			TransportName:   v.DisplayName,
			FilterStatus:    v.FilterStatus,
			IconName:        v.IconName,
			IosIconUrl:      v.IosIconUrl,
			AndroidIconUrl:  v.AndroidIconUrl,
		})
	}

	return
}
