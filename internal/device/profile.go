package device

import (
	"net/http"

	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/favorite"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"bitbucket.movista.ru/maas/maasapi/internal/uapi/auth"
	"github.com/go-pg/pg/v9"
)

func GetOne(ctx common.Context) (out *Device, err error) {
	out = &Device{ID: ctx.DeviceID}
	err = ctx.DB.Model(out).Where("id = ?id").Limit(1).Select()
	if err == pg.ErrNoRows {
		out = nil
	}
	return
}

func GetOldOneByAuthDeviceID(ctx common.Context) (out *Device, err error) {
	var correctDeviceID string
	_, err = ctx.DB.QueryOne(pg.Scan(&correctDeviceID), `SELECT id FROM maasapi.user_profile where auth_device_id = ? limit 1`, ctx.DeviceID)
	if err != nil {
		return
	}

	if correctDeviceID == "" {
		err = models.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       "500",
			Message:    "WRONG_DEVICE_ID",
		}
	}

	out = &Device{ID: correctDeviceID}

	return
}

func GetOneWithSettingsAndFavorites(ctx common.Context) (interface{}, error) {
	// TODO: handle errors
	r, err := GetGoogleTransitModeSettingsWithInfo(ctx)
	if err != nil {
		return nil, err
	}

	var firstName, lastName, phone, email string

	if ctx.IsUser() {
		userInfo, err := auth.GetUserInfo(ctx)
		if err != nil {
			return nil, err
		}

		firstName = userInfo.FirstName
		lastName = userInfo.LastName
		phone = userInfo.Phone
		email = userInfo.Email
	} else {
		up, err := GetOne(ctx)
		if err != nil {
			return nil, err
		}

		firstName = up.Name
	}

	fv, err := favorite.FindFavoriteByDeviceIDWithGooglePlaceInfo(ctx)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"name":              firstName,
		"last_name":         lastName,
		"email":             email,
		"phone":             phone,
		"google_transports": r,
		"favorite_places":   fv,
	}, nil
}

func Create(ctx common.Context, in *Device) (err error) {
	_, err = ctx.DB.Model(in).OnConflict("DO NOTHING").Insert()
	return
}

func Update(ctx common.Context, in *Device) (err error) {
	_, err = ctx.DB.Model(in).WherePK().Returning("*").ExcludeColumn("id").UpdateNotZero()
	return
}
