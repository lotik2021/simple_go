package dialog

import (
	"net/http"
	"time"

	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"bitbucket.movista.ru/maas/maasapi/internal/session"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo/v4"
)

var (
	googleLocSql = `
		select c.place_id, c.search_count, c.direction, gp.coordinate, gp.main_text from (
			select a.* from (
				select origin_place_id as place_id, count(origin_place_id) as search_count, 'from' as direction 
				from maasapi.device_google_search_history
				where device_id in (?)
				and created_at::time between ?::time and ?::time
				and created_at::date between ?::date and ?::date
				and origin_place_id is not null
				group by origin_place_id
	  		) as a
		  	union all
		  	select b.* from (
				select destination_place_id as place_id, count(destination_place_id) as search_count, 'to' as direction 
				from maasapi.device_google_search_history
				where device_id in (?)
				and created_at::time between ?::time and ?::time
				and created_at::date between ?::date and ?::date
				and destination_place_id is not null
				group by destination_place_id
		  	) as b
		) as c
		join maasapi.google_place gp on gp.id = c.place_id
	`

	movistaLocSql = `
		select c.place_id, c.search_count, c.direction, gp.coordinate, gp.main_text from (
			select a.* from (
				select from_google_place_id as place_id, count(from_google_place_id) as search_count, 'from' as direction 
				from maasapi.device_movista_search_history
				where device_id in (?)
				and created_at::time between ?::time and ?::time
				and created_at::date between ?::date and ?::date
				and from_google_place_id is not null
				group by from_google_place_id
			) as a
			union all
			select b.* from (
				select to_google_place_id as place_id, count(to_google_place_id) as search_count, 'to' as direction 
				from maasapi.device_movista_search_history
				where device_id in (?)
				and created_at::time between ?::time and ?::time
				and created_at::date between ?::date and ?::date
				and to_google_place_id is not null
				group by to_google_place_id
			) as b
		) as c
		join maasapi.google_place gp on gp.id = c.place_id

	`
)

type getUserPlaceCountReq struct {
	StartTime string `json:"start_time" example:"14:09:03"` // 10:00:00
	EndTime   string `json:"end_time" example:"19:32:03"`
	StartDate string `json:"start_date" example:"01.10.2019"`
	EndDate   string `json:"end_date" example:"30.10.2019"`
}

func getUserPlaceCount(c echo.Context) (err error) {
	var (
		req getUserPlaceCountReq
		ctx = common.NewContext(c)
	)

	if err = common.BindAndValidateReq(c, &req); err != nil {
		return
	}

	if req.StartTime == "" {
		req.StartTime = "00:00:00"
	}

	if req.EndTime == "" {
		req.EndTime = "23:59:59"
	}

	if req.StartDate == "" {
		req.StartDate = time.Unix(0, 0).Format("2006-01-02")
	} else {
		pt, err := time.Parse("02.01.2006", req.StartDate)
		if err != nil {
			return models.NewInternalDialogError(err)
		}
		req.StartDate = pt.Format("2006-01-02")
	}

	if req.EndDate == "" {
		req.EndDate = time.Unix(99999999999999, 0).Format("2006-01-02")
	} else {
		pt, err := time.Parse("02.01.2006", req.EndDate)
		if err != nil {
			return models.NewInternalDialogError(err)
		}
		req.EndDate = pt.Format("2006-01-02")
	}

	var deviceIds = []string{ctx.DeviceID}

	if ctx.IsUser() {
		deviceIds, err = session.FindActiveDeviceIDsByUserID(ctx)
		if err != nil {
			return
		}
	}

	type placeResult struct {
		Count    int             `json:"count"`
		MainText string          `json:"main_text"`
		Location models.GeoPoint `json:"location"`
	}

	var resp struct {
		FromLoc map[string]placeResult `json:"from_loc"`
		ToLoc   map[string]placeResult `json:"to_loc"`
	}

	var searchType string

	// смотрим последний поиск тип пользователя (вернёт "google" или "movista"), если не было поисков до этого - searchType будет ""
	sqlReq := `
		select c.search_type
		from (
			select a.* from ( select 'google' as search_type, created_at from maasapi.device_google_search_history where device_id in (?) order by created_at asc limit 1) as a
		  	union
		  	select b.* from ( select 'movista' as search_type, created_at from maasapi.device_movista_search_history where device_id in (?) order by created_at asc limit 1) as b
		) as c
		order by c.created_at asc
		limit 1
	`

	_, _ = ctx.DB.Query(&searchType, sqlReq, pg.In(deviceIds), pg.In(deviceIds))

	if searchType == "" {
		return c.JSON(http.StatusOK, echo.Map{
			"result": resp,
		})
	}

	type result struct {
		PlaceID    string          `pg:"place_id"`
		Count      int             `pg:"search_count"`
		Direction  string          `pg:"direction"`
		Coordinate models.GeoPoint `pg:"coordinate"`
		Name       string          `pg:"main_text"`
	}

	var res []result

	if searchType == "google" {
		_, err = ctx.DB.Query(&res, googleLocSql, pg.In(deviceIds), req.StartTime, req.EndTime, req.StartDate, req.EndDate, pg.In(deviceIds), req.StartTime, req.EndTime, req.StartDate, req.EndDate)
		if err != nil {
			return models.NewInternalDialogError(err)
		}
	} else {
		_, err = ctx.DB.Query(&res, movistaLocSql, pg.In(deviceIds), req.StartTime, req.EndTime, req.StartDate, req.EndDate, pg.In(deviceIds), req.StartTime, req.EndTime, req.StartDate, req.EndDate)
		if err != nil {
			return models.NewInternalDialogError(err)
		}
	}

	if len(res) == 0 {
		return c.JSON(http.StatusOK, echo.Map{
			"result": resp,
		})
	}

	for _, v := range res {
		if v.Direction == "from" {
			if resp.FromLoc == nil {
				resp.FromLoc = make(map[string]placeResult)
			}
			resp.FromLoc[v.PlaceID] = placeResult{
				Count:    v.Count,
				MainText: v.Name,
				Location: v.Coordinate,
			}
		} else if v.Direction == "to" {
			if resp.ToLoc == nil {
				resp.ToLoc = make(map[string]placeResult)
			}
			resp.ToLoc[v.PlaceID] = placeResult{
				Count:    v.Count,
				MainText: v.Name,
				Location: v.Coordinate,
			}
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"result": resp,
	})
}
