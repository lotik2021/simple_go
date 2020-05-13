package history

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/google"
	"bitbucket.movista.ru/maas/maasapi/internal/logger"
	"bitbucket.movista.ru/maas/maasapi/internal/session"
	"github.com/go-pg/pg/v9"
	"sync"
)

var (
	sqlRecentFromPlacesByUser = "select from_id as place_id, count(from_id) from ?TableName where device_id = ? group by from_id order by count desc limit ?"
	sqlRecentToPlacesByUser   = "select to_id as place_id, count(to_id) from ?TableName where device_id = ? group by to_id order by count desc limit ?"
	sqlPopularFromPlaces      = "select from_id as place_id, count(from_id) from ?TableName group by from_id order by count desc limit ?"
	sqlPopularToPlaces        = "select to_id as place_id, count(to_id) from ?TableName group by to_id order by count desc limit ?"
)

type PlaceIdCount struct {
	PlaceID int
	Count   int
}

func FindPopularAndRecentMovistaPlaceIDsByDirection(ctx common.Context, direction string, rCount, pCount int) (recentPlaceIdsCount, popularPlaceIdsCount []PlaceIdCount, err error) {
	recentPlaceIdsCount, popularPlaceIdsCount = make([]PlaceIdCount, 0), make([]PlaceIdCount, 0)
	wg := sync.WaitGroup{}

	if ctx.DeviceID != "" {
		wg.Add(1)
		go func(places *[]PlaceIdCount) {
			defer wg.Done()
			var (
				err error
			)
			if direction == "from" {
				_, err = ctx.DB.Model((*MovistaSearch)(nil)).Query(places, sqlRecentFromPlacesByUser, ctx.DeviceID, rCount)
			} else {
				_, err = ctx.DB.Model((*MovistaSearch)(nil)).Query(places, sqlRecentToPlacesByUser, ctx.DeviceID, rCount)
			}
			if err != nil {
				logger.Log.Error(err)
			}
		}(&recentPlaceIdsCount)
	}

	wg.Add(1)
	go func(places *[]PlaceIdCount) {
		defer wg.Done()
		var err error
		if direction == "from" {
			_, err = ctx.DB.Model((*MovistaSearch)(nil)).Query(places, sqlPopularFromPlaces, pCount)
		} else {
			_, err = ctx.DB.Model((*MovistaSearch)(nil)).Query(places, sqlPopularToPlaces, pCount)
		}
		if err != nil {
			logger.Log.Error(err)
		}
	}(&popularPlaceIdsCount)

	wg.Wait()

	return
}

func SaveGooglePlaceMention(ctx common.Context, placeID string) (err error) {
	gpm := &GooglePlaceHistory{
		DeviceID:         ctx.DeviceID,
		PlaceID:          placeID,
		NumberOfSearches: 1,
	}

	updateStatement := `
		ON CONSTRAINT device_google_place_history_device_id_place_id DO UPDATE 
		SET number_of_searches = dgph.number_of_searches + 1, updated_at = now() 
		where dgph.device_id = ? and dgph.place_id = ?
	`

	_, err = ctx.DB.Model(gpm).OnConflict(updateStatement, ctx.DeviceID, placeID).Insert()
	return
}

func GetGooglePlaceMentions(ctx common.Context) (locations []*google.Place, err error) {

	deviceIds := []string{ctx.DeviceID}

	if ctx.IsUser() {
		deviceIds, err = session.FindActiveDeviceIDsByUserID(ctx)
		if err != nil {
			return
		}
	}

	sql := `
		select a.* from (
		  select distinct on (gp.id) gp.id, gp.*, gph.updated_at as last_use_at
		  from maasapi.google_place gp
		  inner join ?TableName gph on gph.place_id = gp.id and gph.device_id in (?)
		) as a
		order by a.last_use_at desc
	`

	_, err = ctx.DB.Model((*GooglePlaceHistory)(nil)).Query(&locations, sql, pg.In(deviceIds))

	for _, v := range locations {
		v.GetIcons()
	}

	return
}
