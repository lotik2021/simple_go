package favorite

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
)

func CopyDeviceFavoritesToUser(ctx common.Context) (err error) {
	favs, err := FindFavorite(ctx)
	if err != nil {
		return
	}

	if len(favs) == 0 {
		return
	}

	userFavs := make([]UserFavorite, 0)
	for _, v := range favs {
		userFavs = append(userFavs, UserFavorite{
			DeviceFavorite: DeviceFavorite{
				DeviceID:  v.DeviceID,
				PlaceID:   v.PlaceID,
				Type:      v.Type,
				Name:      v.Name,
				InActions: v.InActions,
			},
			UserID: ctx.Token.GetUserID(),
		})
	}

	_, err = ctx.DB.Model(&userFavs).Insert()

	return
}
