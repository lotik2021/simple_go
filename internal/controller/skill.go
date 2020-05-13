package controller

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"bitbucket.movista.ru/maas/maasapi/internal/skill"
	"github.com/labstack/echo/v4"
	"net/http"
)

func getSills(c echo.Context) (err error) {
	var (
		req struct {
			Version string `json:"version" query:"version"`
		}
		ctx = common.NewContext(c)
	)

	if err = common.BindAndValidateReq(c, &req); err != nil {
		return err
	}

	skills, err := skill.FindForVersion(ctx, req.Version)
	if err != nil {
		return err
	}

	skillObjects := make([]*models.DataObject, 0)

	for _, s := range skills {
		skillObjects = append(skillObjects, s.ConvertToDataObject())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"result": skillObjects,
	})
}
