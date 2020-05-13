package search

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"bitbucket.movista.ru/maas/maasapi/internal/transport/taxi"
	"github.com/google/uuid"
)

func convertTaxiRoutesToDataObject(trips []*taxi.Trip) (dataObject models.DataObject) {

	if len(trips) == 0 {
		return
	}

	routes := make([]models.DataObject, 0)

	for _, v := range trips {
		routes = append(routes, models.DataObject{ObjectId: v.ObjectId, Id: uuid.New().String(), Data: v})
	}

	dataObject = models.DataObject{
		ObjectId: GROUP_OF_TAXI,
		Data: models.RouteObject{
			Routes: routes,
		},
	}
	return
}

func findTaxiRoutes(ctx common.Context, origin, destination models.GeoPoint, departureTime, arrivalTime models.Time, region string) (trips []*taxi.Trip, deeplinks map[string]models.DeepLink) {
	return taxi.CalculateRoute(ctx, origin, destination, departureTime, arrivalTime, region)
}
