package controller

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/history"
	"bitbucket.movista.ru/maas/maasapi/internal/uapi/place"
	"github.com/labstack/echo/v4"
	"net/http"
)

type WebPlacesAutoSuggestionsReq struct {
	Direction    string `query:"direction" validate:"required,oneof=from to"`
	RecentCount  *int   `query:"recentCount"`
	PopularCount *int   `query:"popularCount"`
}

type WebPlacesAutoSuggestionsRes struct {
	Data struct {
		Recent  []place.Place `json:"recent"`
		Popular []place.Place `json:"popular"`
	} `json:"data"`
}

// webPlacesAutoSuggestions godoc
// @Tags web
// @Summary webPlacesAutoSuggestions
// @Description return recent and popular places
// @ID web-places-auto-suggestions
// @Accept  json
// @Produce json
// @Param AS-CID header string false "user_id to get recent, if empty - recent is []"
// @Param direction query string true "direction" Enums(from, to)
// @Param recentCount query int false "recent count in response" minimum(1) default(5)
// @Param popularCount query int false "popular count in response" minimum(1) default(5)
// @Success 200 {object} router.WebPlacesAutoSuggestionsRes
// @Router /api/places/autosuggestions [get]
func webPlacesAutoSuggestions(c echo.Context) (err error) {
	var (
		req WebPlacesAutoSuggestionsReq
		ctx = common.NewContext(c)
	)
	err = common.BindAndValidateReq(c, &req)
	if err != nil {
		return err
	}

	if req.RecentCount == nil {
		*req.RecentCount = 5
	}

	if req.PopularCount == nil {
		*req.PopularCount = 5
	}

	if ascid := c.Request().Header.Get("AS-CID"); ascid != "" {
		ctx.DeviceID = ascid
	}

	var res WebPlacesAutoSuggestionsRes

	recentPlaceIdsCount, popularPlaceIdsCount, err := history.FindPopularAndRecentMovistaPlaceIDsByDirection(ctx, req.Direction, *req.RecentCount, *req.PopularCount)
	if err != nil {
		return err
	}

	if len(recentPlaceIdsCount) == 0 && len(popularPlaceIdsCount) == 0 {
		return
	}

	allPlaceIds := make([]int, 0, len(recentPlaceIdsCount)+len(popularPlaceIdsCount))
	for _, v := range recentPlaceIdsCount {
		allPlaceIds = append(allPlaceIds, v.PlaceID)
	}
	for _, v := range popularPlaceIdsCount {
		allPlaceIds = append(allPlaceIds, v.PlaceID)
	}

	// можно отправлять массив с дупликатами
	placesFromUAPI, err := place.FindByIds(ctx, allPlaceIds)
	if err != nil {
		return
	}

	// убрать дубликаты из recent и popular
	uniquePlaces := make(map[int]bool)

	for _, v := range recentPlaceIdsCount {
		if uniquePlaces[v.PlaceID] {
			continue
		}

		place, ok := placesFromUAPI[v.PlaceID]
		if !ok {
			continue
		}
		uniquePlaces[v.PlaceID] = true
		res.Data.Recent = append(res.Data.Recent, place)
	}

	for _, v := range popularPlaceIdsCount {
		if uniquePlaces[v.PlaceID] {
			continue
		}

		place, ok := placesFromUAPI[v.PlaceID]
		if !ok {
			continue
		}
		uniquePlaces[v.PlaceID] = true
		res.Data.Popular = append(res.Data.Popular, place)
	}

	if len(res.Data.Popular) == 0 {
		res.Data.Popular = make([]place.Place, 0)
	}

	if len(res.Data.Recent) == 0 {
		res.Data.Recent = make([]place.Place, 0)
	}

	return c.JSON(http.StatusOK, res)
}
