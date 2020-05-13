package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"github.com/labstack/echo/v4"
)

type Request struct {
	//Client         string      `json:"client,omitempty"`
	Meta           Meta        `json:"meta"`
	ActionId       string      `json:"action_id" example:"yes"`
	UserResponse   string      `json:"user_response" example:"Москва"`
	Session        Session     `json:"session"`
	ClientEntities interface{} `json:"client_entities"`
	Version        string      `json:"version"`
}

type Meta struct {
	ClientID string           `json:"client_id"`
	Timezone string           `json:"time_zone"`
	Location *models.GeoPoint `json:"location"`
	Locale   string           `json:"locale"`
	UserId   string           `json:"user_id" example:"009dew4enew1711123"`
}

type Session struct {
	MessageId       *int64       `json:"message_id,omitempty"`
	SessionId       string       `json:"session_id" example:"a6e6f1c8-5f61-43c8-b9ee-f1dd176108cf"`
	MessageDateTime *models.Time `json:"message_date_time,omitempty" example:"2019-08-08T14:15:22+03:00"`
}
type Response struct {
	Actions *[]Action            `json:"actions,omitempty"`
	Objects *[]models.DataObject `json:"objects,omitempty"`
	Session Session              `json:"session"`
	Hint    *string              `json:"hint,omitempty"`
	Extras  interface{}          `json:"extras,omitempty"`
	Error   *Error               `json:"error,omitempty"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Action struct {
	Handable bool   `json:"handable"`
	Title    string `json:"title"`
	ActionId string `json:"action_id"`
	Color    string `json:"color"`
}

// Dialog godoc
// @Tags dialog
// @Summary Dialog
// @Description Dialog
// @ID dialog
// @Accept json
// @Produce json
// @Param req body dialog.Request false "Dialog"
// @Success 200 {object} dialog.Response
// @Router /api/dialogs/dialog [post]
func dialog(c echo.Context) (err error) {
	var (
		ctx = common.NewContext(c)
		req = Request{}
		res = Response{}
	)

	err = common.BindAndValidateReq(c, &req)
	if err != nil {
		return
	}

	spanReq := common.DefaultRequest.Clone().Post(config.C.Dialog.Urls.Dialog).AppendHeader("Authorization", ctx.Token.String()).SendStruct(req)

	body, _, err := common.SendRequest(ctx, spanReq)
	if err != nil {
		return fmt.Errorf("cannot make request to dialog - %w", err)
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		return fmt.Errorf("cannot unmarshal response from dialog - %s", string(body))
	}

	return c.JSON(http.StatusOK, res)
}
