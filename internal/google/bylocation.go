package google

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"context"
	"fmt"
	"github.com/go-pg/pg/v9"
	"googlemaps.github.io/maps"
)

func FindByLocationInDB(ctx common.Context, point *models.GeoPoint) (placeID string, err error) {
	var p Place

	err = ctx.DB.Model(&p).
		Where("ST_X(coordinate::geometry)::numeric = ?", point.Longitude).
		Where("ST_Y(coordinate::geometry)::numeric = ?", point.Latitude).
		Limit(1).
		Select()
	if err != nil && err != pg.ErrNoRows {
		return
	}

	placeID = p.ID

	if placeID == "" {
		placeID = fmt.Sprint(point.Latitude, point.Longitude)
	} else {
		placeID = "place_id:" + placeID
	}

	return
}

func FindByLocation(ctx common.Context, point *models.GeoPoint, region string) (place Place, err error) {
	locationRequest := &maps.GeocodingRequest{
		LatLng:   &maps.LatLng{Lat: point.Latitude, Lng: point.Longitude},
		Language: region,
	}

	result, err := mapsClient.ReverseGeocode(context.Background(), locationRequest)
	if err != nil {
		return
	}

	mainText, secondaryText, placeId, placeTypes := getAddressTexts(result)
	place = Place{
		ID:            placeId,
		MainText:      mainText,
		SecondaryText: secondaryText,
		Coordinate:    point,
		PlaceTypes:    placeTypes,
	}

	return
}

func getAddressTexts(response []maps.GeocodingResult) (mainText string, secondaryText string, placeId string, placeTypes []string) {
	for _, geo := range response {
		if typeContains(geo.Types, streetAddressTypes) {
			mainText = geo.FormattedAddress
			placeId = geo.PlaceID
			for _, ac := range geo.AddressComponents {
				if typeContains(ac.Types, cityAddressTypes) {
					secondaryText = ac.LongName
					placeTypes = ac.Types
					return
				}
			}
		}
	}
	return
}

func typeContains(types []string, referenceTypes []string) bool {
	cnt := 0
	for _, tp := range types {
		for _, atp := range referenceTypes {
			if tp == atp {
				cnt++
			}
		}
	}

	return cnt == len(types)
}
