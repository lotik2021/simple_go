package dialog

import (
	"fmt"
	"net/http"

	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"bitbucket.movista.ru/maas/maasapi/internal/session"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo/v4"
)

func createSessionData(c echo.Context) (err error) {
	var (
		req struct {
			SessionID      string `json:"session_id" validate:"required"`
			ActionID       string `json:"action_id"`
			ActionName     string `json:"action_name"`
			UserResponse   string `json:"user_response"`
			DialogResponse string `json:"dialog_response"`
			UserEntryData  string `json:"user_entry_data,omitempty"`
			Actions        string `json:"actions,omitempty"`
			Objects        string `json:"objects,omitempty"`
		}
		ctx = common.NewContext(c)
	)

	err = common.BindAndValidateReq(c, &req)
	if err != nil {
		return models.NewInternalDialogError(err)
	}

	s := &DeviceSessionData{
		SessionID:      req.SessionID,
		ActionID:       req.ActionID,
		ActionName:     req.ActionName,
		UserResponse:   req.UserResponse,
		DialogResponse: req.DialogResponse,
		UserEntryData:  req.UserEntryData,
		Actions:        req.Actions,
		Objects:        req.Objects,
	}

	_, err = ctx.DB.Model(s).Insert()
	if err != nil {
		return models.NewInternalDialogError(err)
	}

	return c.NoContent(http.StatusOK)
}

func updateLastSessionData(c echo.Context) (err error) {
	var req struct {
		SessionID      string `json:"session_id" validate:"required"`
		DialogResponse string `json:"dialog_response"`
		UserEntryData  string `json:"user_entry_data,omitempty"`
		Objects        string `json:"objects,omitempty"`
		Actions        string `json:"actions,omitempty"`
	}

	ctx := common.NewContext(c)

	err = common.BindAndValidateReq(c, &req)
	if err != nil {
		return models.NewInternalDialogError(err)
	}

	var sID string
	err = ctx.DB.Model((*DeviceSessionData)(nil)).Column("id").Where("session_id = ?", req.SessionID).Order("created_at desc").Select(&sID)

	if err != nil {
		return
	}

	if sID == "" {
		return models.NewInternalDialogError(fmt.Errorf("cannot find last deviceSessionData by session_id %s", req.SessionID))
	}

	s := &DeviceSessionData{
		ID:             sID,
		DialogResponse: req.DialogResponse,
		UserEntryData:  req.UserEntryData,
		Objects:        req.Objects,
		Actions:        req.Actions,
	}

	_, err = ctx.DB.Model(s).WherePK().Returning("*").UpdateNotZero()
	if err != nil {
		return models.NewInternalDialogError(err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"result": s,
	})
}

func getSessionData(c echo.Context) (err error) {
	var req struct {
		SessionId string `json:"session_id"`
		Order     string `json:"order"`
		Limit     *int   `json:"limit"`
	}

	ctx := common.NewContext(c)

	err = common.BindAndValidateReq(c, &req)
	if err != nil {
		return models.NewInternalDialogError(err)
	}

	s := &DeviceSessionData{SessionID: req.SessionId}

	sessionData := make([]*DeviceSessionData, 0)

	q := ctx.DB.Model(s)

	if req.SessionId != "" {
		q = q.Where("session_id = ?", req.SessionId)
	}

	if req.Order != "" {
		q = q.Order(req.Order)
	} else {
		q = q.Order("created_at asc")
	}

	if req.Limit != nil {
		q = q.Limit(*req.Limit)
	}

	err = q.Select(&sessionData)
	if err != nil {
		return models.NewInternalDialogError(err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"result": sessionData,
	})
}

type getUserActionCountReq struct {
	StartTime string `json:"start_time" validate:"required" example:"14:09:03"` // 10:00:00
	EndTime   string `json:"end_time" validate:"required" example:"19:32:03"`
}

func getUserActionCount(c echo.Context) (err error) {
	var (
		req getUserActionCountReq
		ctx = common.NewContext(c)
	)

	if err = common.BindAndValidateReq(c, &req); err != nil {
		return
	}

	var result []struct {
		ActionId string
		Count    int
	}

	var deviceIds = []string{ctx.DeviceID}

	if ctx.IsUser() {
		deviceIds, err = session.FindActiveDeviceIDsByUserID(ctx)
		if err != nil {
			return
		}
		if len(deviceIds) == 0 {
			deviceIds = []string{ctx.DeviceID}
		}
	}

	sqlReq := `
	select action_id, count(*)
	from maasapi.device_session_data
	where 
	created_at::time between ?::time and ?::time
	and
	session_id in (
		select id from maasapi.device_session where device_id in (?)
	)
	group by action_id
	`

	_, err = ctx.DB.Query(&result, sqlReq, req.StartTime, req.EndTime, pg.In(deviceIds))
	if err != nil {
		return
	}

	resultMap := make(map[string]int, 0)

	for _, v := range result {
		resultMap[v.ActionId] = v.Count
	}

	return c.JSON(http.StatusOK, echo.Map{
		"result": resultMap,
	})
}
