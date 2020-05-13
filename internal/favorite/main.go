package favorite

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"bitbucket.movista.ru/maas/maasapi/internal/google"
)

const (
	WorkType = "work"
	HomeType = "home"
	NoneType = "none"
)

func Create(ctx common.Context, in *DeviceFavorite) (id int, created bool, err error) {

	if in.Name != "" {
		in.Type = NoneType
	}

	if ctx.IsUser() {
		created, err = ctx.DB.
			Model(&UserFavorite{UserID: ctx.UserID, DeviceFavorite: *in}).
			Where("type = ?type and place_id = ?place_id and device_id = ?device_id and user_id = ?user_id").
			OnConflict("DO NOTHING").SelectOrInsert()
	} else {
		created, err = ctx.DB.
			Model(in).
			Where("type = ?type and place_id = ?place_id and device_id = ?device_id").
			OnConflict("DO NOTHING").SelectOrInsert()
	}

	if err != nil {
		return
	}

	id = in.ID

	return
}

func UpdateFavorite(ctx common.Context, in *DeviceFavorite) (err error) {
	if in.Name != "" {
		in.Type = NoneType
	}

	if ctx.IsUser() {
		err = ctx.DB.Update(&UserFavorite{UserID: ctx.UserID, DeviceFavorite: *in})
	} else {
		err = ctx.DB.Update(&in)
	}

	if err != nil {
		return
	}

	return
}

func DeleteFavorite(ctx common.Context, id int) (err error) {
	if ctx.IsUser() {
		err = ctx.DB.Delete(&UserFavorite{UserID: ctx.UserID, DeviceFavorite: DeviceFavorite{ID: id}})
	} else {
		err = ctx.DB.Delete(&DeviceFavorite{ID: id, DeviceID: ctx.DeviceID})
	}

	if err != nil {
		return
	}

	return
}

func FindFavoriteByDeviceIDWithGooglePlaceInfo(ctx common.Context) (resultList []WithGooglePlaceInfo, err error) {
	list, err := FindFavorite(ctx)
	if err != nil {
		return
	}

	resultList = make([]WithGooglePlaceInfo, 0)

	if len(list) == 0 {
		return
	}

	var googlePlaceIDs []string
	for _, li := range list {
		googlePlaceIDs = append(googlePlaceIDs, li.PlaceID)
	}

	res, err := google.FindByIDsInDB(ctx, googlePlaceIDs)
	if err != nil {
		return
	}

	for _, item := range list {
		gp, ok := res[item.PlaceID]
		if !ok {
			continue
		}

		gp.GetIcons()

		li := WithGooglePlaceInfo{
			ID:        item.ID,
			DeviceID:  item.DeviceID,
			Type:      item.Type,
			Name:      item.Name,
			InActions: item.InActions,
			UpdatedAt: item.UpdatedAt,
			GooglePlace: google.Place{
				ID:                   gp.ID,
				Coordinate:           gp.Coordinate,
				MainText:             gp.MainText,
				SecondaryText:        gp.SecondaryText,
				PlaceTypes:           gp.PlaceTypes,
				IosIconURLLightTheme: gp.IosIconURLLightTheme,
				IosIconURLDarkTheme:  gp.IosIconURLDarkTheme,
				IosIconURL:           gp.IosIconURLDarkTheme,
				AndroidIconURL:       gp.AndroidIconURL,
			},
		}

		if item.Type == WorkType {
			li.GooglePlace.IosIconURLLightTheme = config.C.Icons.Places.Work.Ios
			li.GooglePlace.IosIconURLDarkTheme = config.C.Icons.Places.DarkTheme.Work.Ios
			li.GooglePlace.IosIconURL = li.GooglePlace.IosIconURLDarkTheme
			li.GooglePlace.AndroidIconURL = config.C.Icons.Places.Work.Android
		} else if item.Type == HomeType {
			li.GooglePlace.IosIconURLLightTheme = config.C.Icons.Places.Home.Ios
			li.GooglePlace.IosIconURLDarkTheme = config.C.Icons.Places.DarkTheme.Home.Ios
			li.GooglePlace.IosIconURL = li.GooglePlace.IosIconURLDarkTheme
			li.GooglePlace.AndroidIconURL = config.C.Icons.Places.Home.Android
		}

		resultList = append(resultList, li)
	}

	return
}

func FindFavoriteByDeviceIDAndType(ctx common.Context, fType string) (favs []DeviceFavorite, err error) {

	favs = make([]DeviceFavorite, 0)

	err = ctx.DB.Model(&favs).Where("device_id = ?", ctx.DeviceID).Where("type = ?", fType).Select()
	if err != nil {
		return
	}

	return
}

func FindFavorite(ctx common.Context) (favs []DeviceFavorite, err error) {

	favs = make([]DeviceFavorite, 0)

	if ctx.IsUser() {
		userFavs := make([]UserFavorite, 0)
		err = ctx.DB.Model(&UserFavorite{}).Where("user_id = ?", ctx.UserID).Select(&userFavs)
		if len(userFavs) > 0 {
			for _, v := range userFavs {
				favs = append(favs, v.DeviceFavorite)
			}
		}
	} else {
		err = ctx.DB.Model(&DeviceFavorite{}).Where("device_id = ?", ctx.DeviceID).Select(&favs)
	}

	if err != nil {
		return
	}

	return
}

func ToggleFavoriteInActions(ctx common.Context, id int) (err error) {
	if ctx.IsUser() {
		_, err = ctx.DB.Model((*UserFavorite)(nil)).Exec(`UPDATE ?TableName SET in_actions = NOT in_actions where user_id = ? and id = ? and type not in ('home','work')`, ctx.UserID, id)
	} else {
		_, err = ctx.DB.Model((*DeviceFavorite)(nil)).Exec(`UPDATE ?TableName SET in_actions = NOT in_actions where device_id = ? and id = ? and type not in ('home','work')`, ctx.DeviceID, id)
	}

	return
}
