package google

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"bitbucket.movista.ru/maas/maasapi/internal/logger"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"github.com/go-pg/pg/v9"
	"github.com/google/uuid"
	"googlemaps.github.io/maps"
)

var (
	cityAddressTypes   = []string{"locality", "political"}
	streetAddressTypes = []string{"street_address"}
	mapsClient         *maps.Client
)

func init() {
	var err error
	mapsClient, err = maps.NewClient(maps.WithAPIKey(config.C.Google.ApiKey))
	if err != nil {
		logger.Log.Fatal(err)
	}
}

func Save(ctx common.Context, in *Place) (placeID string, err error) {
	// TODO: если координаты нулевые и placeId != "", тогда запишется в базу place с пустыми координатами
	if in.ID != "" && in.Coordinate.IsZero() {
		var pr ShortPlace
		pr, err = FindByID(ctx, uuid.Nil, in.ID, "ru")
		if err != nil {
			return
		}

		in.Coordinate = &models.GeoPoint{
			Longitude: pr.Location.Longitude,
			Latitude:  pr.Location.Latitude,
		}
	}

	if in.ID == "" {
		var pr Place
		pr, err = FindByLocation(ctx, in.Coordinate, "ru")
		if err != nil {
			return
		}

		in.ID = pr.ID
	}

	placeID = in.ID

	updateStatement := `
		ON CONSTRAINT google_place_pkey DO UPDATE 
		SET place_types = ?, updated_at = now()
		where gp.id = ?
	`

	_, err = ctx.DB.Model(in).Where("id = ?", in.ID).OnConflict(updateStatement, pg.Array(in.PlaceTypes), placeID).Insert()

	return
}
