package device

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
)

func PingLocation(ctx common.Context, location *models.GeoPoint) (err error) {
	p := &Location{
		DeviceID:   ctx.DeviceID,
		Coordinate: location,
	}

	_, err = ctx.DB.Model(p).Insert()
	if err != nil {
		return
	}

	return
}

func GetLastLocation(ctx common.Context) (loc *models.GeoPoint) {

	var p Location

	_ = ctx.DB.Model(&p).Where("device_id = ?", ctx.DeviceID).
		Order("created_at desc").Limit(1).Select()

	if !p.Coordinate.IsZero() {
		loc = p.Coordinate
	}

	return
}
