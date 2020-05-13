package search

import (
	"fmt"
	"strings"

	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"bitbucket.movista.ru/maas/maasapi/internal/uapi/fapi"
	"github.com/goodsign/monday"
	"github.com/google/uuid"
)

func convertMovistaFastestLowestNearestRoutesToDataObjects(res *fapi.FastestLowestNearestRes) (dataObjects []models.DataObject) {

	dataObjects = make([]models.DataObject, 0)

	if res == nil {
		return
	}

	dataObjects = append(dataObjects, models.DataObject{
		ObjectId: TEXT_DATA,
		Data: models.MessageObject{
			Text: fmt.Sprintf("%s - %s, %s",
				res.OriginName,
				res.DestinationName,
				strings.ToLower(monday.Format(res.DepartureTime.Time, "Monday 2 January 2006", monday.LocaleRuRU))),
		},
	})

	dataObjects = append(dataObjects, models.DataObject{ObjectId: TEXT_DATA, Data: models.MessageObject{Text: "Самый дешевый маршрут"}})
	dataObjects = append(dataObjects, models.DataObject{
		ObjectId: GROUP_OF_MOVISTA,
		Data:     models.RouteObject{Routes: []models.DataObject{{ObjectId: MOVISTA_ROUTE, Id: res.Routes["lowest"].ID, Data: res.Routes["lowest"]}}},
	},
	)

	dataObjects = append(dataObjects, models.DataObject{ObjectId: TEXT_DATA, Data: models.MessageObject{Text: "Самый быстрый маршрут"}})
	dataObjects = append(dataObjects, models.DataObject{
		ObjectId: GROUP_OF_MOVISTA,
		Data:     models.RouteObject{Routes: []models.DataObject{{ObjectId: MOVISTA_ROUTE, Id: res.Routes["fastest"].ID, Data: res.Routes["fastest"]}}},
	},
	)

	dataObjects = append(dataObjects, models.DataObject{ObjectId: TEXT_DATA, Data: models.MessageObject{Text: "Ближайший маршрут"}})
	dataObjects = append(dataObjects, models.DataObject{
		ObjectId: GROUP_OF_MOVISTA,
		Data:     models.RouteObject{Routes: []models.DataObject{{ObjectId: MOVISTA_ROUTE, Id: res.Routes["nearest"].ID, Data: res.Routes["nearest"]}}},
	},
	)

	ttl := make(map[string]interface{})
	ttl["total_count"] = res.Count
	ttl["redirect_url"] = res.Routes["fastest"].RedirectUrl
	ttl["id"] = uuid.New().String()

	dataObjects = append(dataObjects, models.DataObject{ObjectId: MOVISTA_PLACEHOLDER, Data: ttl})

	return
}

func findMovistaFastestLowestNearestRoutes(ctx common.Context, departureTime models.Time, from, to models.GeoPoint, tripTypes map[string]bool) *fapi.FastestLowestNearestRes {
	res, err := fapi.FindFastestLowestNearestRoutesV4(ctx, departureTime, from, to, tripTypes)
	if err != nil {
		return res
	}

	return res
}
