package dialog

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/google"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

type getGooglePlaceInfoReq struct {
	PlaceID string `json:"place_id" validate:"required" example:"ChIJBcroiVxKtUYRGTTp5LLFaN8"`
}

type getGooglePlaceInfoRes struct {
	Coordinate *models.GeoPoint `json:"coordinate"`
	ID         string           `json:"google_place_id"`
	google.Place
}

func getGooglePlaceInfo(c echo.Context) (err error) {
	var (
		req getGooglePlaceInfoReq
		ctx = common.NewContext(c)
	)

	err = common.BindAndValidateReq(c, &req)
	if err != nil {
		return models.NewInternalDialogError(err)
	}

	gp, err := google.FindByIDInDB(ctx, req.PlaceID)
	if err != nil {
		return models.NewInternalDialogError(err)
	}

	res := getGooglePlaceInfoRes{
		Coordinate: gp.Coordinate,
		ID:         gp.ID,
		Place:      gp,
	}

	return c.JSON(http.StatusOK, echo.Map{
		"result": res,
	})
}
