package google

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"context"
	"github.com/go-pg/pg/v9"
	"github.com/google/uuid"
	"googlemaps.github.io/maps"
)

func FindByID(ctx common.Context, sessionToken uuid.UUID, placeID string, region string) (result ShortPlace, err error) {

	if sessionToken == uuid.Nil {
		sessionToken = uuid.UUID(maps.NewPlaceAutocompleteSessionToken())
	}

	placeRequest := &maps.PlaceDetailsRequest{
		PlaceID:      placeID,
		Language:     region,
		SessionToken: maps.PlaceAutocompleteSessionToken(sessionToken),
	}

	placeResult, err := mapsClient.PlaceDetails(context.Background(), placeRequest)
	if err != nil {
		return
	}

	result = ShortPlace{
		PlaceID:     placeResult.PlaceID,
		Description: placeResult.FormattedAddress,
		Location:    &models.GeoPoint{Longitude: placeResult.Geometry.Location.Lng, Latitude: placeResult.Geometry.Location.Lat},
	}

	return
}

func FindByIDInDB(ctx common.Context, id string) (result Place, err error) {
	err = ctx.DB.Model(&result).Where("id = ?", id).First()
	if err != nil {
		return
	}

	return
}

func FindByIDsInDB(ctx common.Context, ids []string) (result map[string]Place, err error) {
	var (
		places []Place
	)
	result = make(map[string]Place)

	err = ctx.DB.Model(&places).Where("id in (?)", pg.In(ids)).Select()
	if err != nil {
		return
	}

	if len(places) == 0 {
		return
	}

	for _, p := range places {
		result[p.ID] = p
	}

	return
}
