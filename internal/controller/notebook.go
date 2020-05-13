package controller

import (
	"net/http"

	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"bitbucket.movista.ru/maas/maasapi/internal/notebook"
	"github.com/labstack/echo/v4"
)

func getDocumentsV2(c echo.Context) error {
	var (
		ctx = common.NewContext(c)
	)

	docs, err := notebook.GetDocuments(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, models.UapiResponse{
		Data: docs,
	})
}

func documentsDeactivateV2(c echo.Context) error {
	var (
		in struct {
			DocumentID int64 `json:"documentId"`
		}
		ctx = common.NewContext(c)
	)

	if err := common.BindAndValidateReq(c, &in); err != nil {
		return err
	}

	resp, err := notebook.DeactivateDocument(ctx, in.DocumentID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func documentsEditV2(c echo.Context) error {
	var (
		in  notebook.Customer
		ctx = common.NewContext(c)
	)

	if err := common.BindAndValidateReq(c, &in); err != nil {
		return err
	}

	for _, doc := range in.CustomerDocuments {
		_, err := notebook.DeactivateDocument(ctx, doc.ID)
		if err != nil {
			return err
		}
	}

	newCustomer, err := notebook.CreateDocuments(ctx, in)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, models.UapiResponse{
		Data: newCustomer,
	})
}

func documentsCreateV2(c echo.Context) error {
	var (
		in  notebook.Customer
		ctx = common.NewContext(c)
	)

	if err := common.BindAndValidateReq(c, &in); err != nil {
		return err
	}

	in.ID = 0
	if len(in.CustomerDocuments) > 0 {
		in.CustomerDocuments[0].ID = 0
	}

	newCustomer, err := notebook.CreateDocuments(ctx, in)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, models.UapiResponse{
		Data: newCustomer,
	})
}
