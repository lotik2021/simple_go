package dialog

import (
	"net/http"
	"time"

	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/device"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"github.com/labstack/echo/v4"
)

func createSession(c echo.Context) (err error) {

	var req struct {
		SessionID string `json:"session_id" validate:"required"`
		State     string `json:"state" validate:"required"`
	}

	ctx := common.NewContext(c)

	err = common.BindAndValidateReq(c, &req)
	if err != nil {
		return models.NewInternalDialogError(err)
	}

	_, err = ctx.DB.Model(&DeviceSession{ID: req.SessionID, DeviceID: ctx.DeviceID, State: req.State}).OnConflict("DO NOTHING").Insert()
	if err != nil {
		return models.NewInternalDialogError(err)
	}

	// get session count
	dev := &device.Device{}
	err = ctx.DB.Model(dev).Where("id = ?", ctx.DeviceID).Select()
	if err != nil {
		return models.NewInternalDialogError(err)
	}

	ctx.DB.Model(dev).Set("session_count = ?", dev.SessionCount+1).Where("id = ?", ctx.DeviceID).Update()

	return c.JSON(http.StatusOK, echo.Map{
		"result": echo.Map{
			"session_count":     dev.SessionCount,
			"last_session_time": &dev.UpdatedAt,
		},
	})
}

func getSession(c echo.Context) (err error) {
	var req struct {
		ID string `json:"session_id" validate:"required"`
	}

	ctx := common.NewContext(c)

	err = common.BindAndValidateReq(c, &req)
	if err != nil {
		return models.NewInternalDialogError(err)
	}

	s := &DeviceSession{ID: req.ID}

	err = ctx.DB.Select(s)
	if err != nil {
		return models.NewInternalDialogError(err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"result": s,
	})
}

func updateSession(c echo.Context) (err error) {
	var (
		req struct {
			ID                 string           `json:"session_id" validate:"required"`
			Origin             *models.GeoPoint `json:"origin,omitempty"`
			Destination        *models.GeoPoint `json:"destination,omitempty"`
			OriginPlaceId      string           `json:"origin_place_id,omitempty"`
			DestinationPlaceId string           `json:"destination_place_id,omitempty"`
			DepartureTime      string           `json:"departure_time,omitempty"`
			ArrivalTime        string           `json:"arrival_time,omitempty"`
			State              string           `json:"state,omitempty"`
			TripTypes          []string         `json:"trip_types,omitempty"`
		}
		ctx = common.NewContext(c)
		s   = &DeviceSession{}
	)

	err = common.BindAndValidateReq(c, &req)
	if err != nil {
		return models.NewInternalDialogError(err)
	}

	s.ID = req.ID
	s.DeviceID = ctx.DeviceID

	if req.DepartureTime != "" {
		t, _ := time.Parse(time.RFC3339, req.DepartureTime)
		s.DepartureTime = &models.Time{Time: t}
	}

	if req.ArrivalTime != "" {
		t, _ := time.Parse(time.RFC3339, req.ArrivalTime)
		s.ArrivalTime = &models.Time{Time: t}
	}

	if req.State != "" {
		s.State = req.State
	} else {
		s.State = "init"
	}

	if !req.Origin.IsZero() {
		s.Origin = req.Origin
	}

	if req.OriginPlaceId != "" {
		s.OriginPlaceId = req.OriginPlaceId
	}

	if !req.Destination.IsZero() {
		s.Destination = req.Destination
	}

	if req.DestinationPlaceId != "" {
		s.DestinationPlaceId = req.DestinationPlaceId
	}

	if req.TripTypes != nil && len(req.TripTypes) > 0 {
		s.TripTypes = req.TripTypes
	}

	_, err = ctx.DB.Model(s).WherePK().ExcludeColumn("device_id", "created_at").Update()
	if err != nil {
		return models.NewInternalDialogError(err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"result": "ok",
	})
}
