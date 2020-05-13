package fapi

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"bitbucket.movista.ru/maas/maasapi/internal/logger"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/url"
)

func GetBlankBase64(c echo.Context) (err error) {
	var (
		in struct {
			FileStorageIDs []string `json:"file_storage_ids" validate:"required"`
		}
		ctx = common.NewContext(c)
	)

	if err = common.BindAndValidateReq(c, &in); err != nil {
		return err
	}

	uapires := struct {
		Data string `json:"data"`
	}{}

	files := make(map[string]string, 0)
	for _, id := range in.FileStorageIDs {
		url := config.C.FapiAdapter.Urls.Blank + fmt.Sprintf("/%s", id) + fmt.Sprintf("/base64")
		req := faClient.Clone().Get(url)

		resp, _, err := common.SendRequest(ctx, req)
		if err != nil {
			logger.Log.Errorf("uapi blanks error: %w", err)
			break
		}

		err = json.Unmarshal(resp, &uapires)
		if err != nil {
			logger.Log.Errorf("cannot unmarshal uapi response: %w", err)
		}

		files[id] = uapires.Data
	}

	return c.JSON(http.StatusOK, echo.Map{"data": files})
}

func GetOpenedBlankFile(c echo.Context) (err error) {
	var (
		in struct {
			FileStorageID string `json:"file_storage_id" validate:"required"`
		}
		ctx = common.NewContext(c)
	)

	if err = common.BindAndValidateReq(c, &in); err != nil {
		return err
	}

	url := config.C.FapiAdapter.Urls.Blank + fmt.Sprintf("/%s", in.FileStorageID) + fmt.Sprintf("/file")

	req := faClient.Clone().Get(url)

	resp, _, err := common.SendRequest(ctx, req)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, json.RawMessage(resp))
}

func GetBlankPDF(c echo.Context) (err error) {
	var (
		in struct {
			FileStorageID string `json:"file_storage_id" validate:"required"`
		}
		ctx = common.NewContext(c)
	)

	if err = common.BindAndValidateReq(c, &in); err != nil {
		return err
	}

	url := config.C.FapiAdapter.Urls.Blank + fmt.Sprintf("/%s", in.FileStorageID)

	req := faClient.Clone().Get(url)

	resp, _, err := common.SendRequest(ctx, req)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, json.RawMessage(resp))
}

func GetBlankPDFs(c echo.Context) (err error) {
	var (
		in struct {
			FileStorageID []string `json:"file_storage_ids" validate:"required"`
		}
		ctx = common.NewContext(c)
	)

	if err = common.BindAndValidateReq(c, &in); err != nil {
		return err
	}

	qparams := url.Values{}

	for _, id := range in.FileStorageID {
		qparams.Add("fileStorageIds", id)
	}

	url := config.C.FapiAdapter.Urls.BlanksPdf + fmt.Sprintf("?%s", qparams.Encode())

	req := faClient.Clone().Get(url)

	resp, httpInternals, err := common.SendRequest(ctx, req)
	if err != nil {
		return
	}

	c.Response().Before(func() {
		c.Response().Header().Set(echo.HeaderContentDisposition, httpInternals.Header.Get("Content-Disposition"))
	})

	return c.Blob(http.StatusOK, "application/pdf", resp)
}

func GetBlankZIP(c echo.Context) (err error) {
	var (
		in struct {
			FileStorageIDs []string `json:"file_storage_id" validate:"required"`
		}
		ctx = common.NewContext(c)
	)

	if err = common.BindAndValidateReq(c, &in); err != nil {
		return err
	}

	v := url.Values{}

	for _, id := range in.FileStorageIDs {
		v.Add("fileStorageIds", id)
	}

	// get-запрос с массивом из ids в query
	url := config.C.FapiAdapter.Urls.BlanksZip + fmt.Sprintf("?%s", v.Encode())

	req := faClient.Clone().Get(url)

	resp, _, err := common.SendRequest(ctx, req)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, json.RawMessage(resp))
}

func WebGetBlankBase64(c echo.Context) (err error) {
	var (
		ctx = common.NewContext(c)
	)

	fileStorageID := c.Param("fileStorageId")
	if fileStorageID == "" {
		return fmt.Errorf("empty fileStorageId param")
	}

	url := config.C.FapiAdapter.Urls.Blank + fmt.Sprintf("/%s", fileStorageID) + fmt.Sprintf("/base64")

	req := faClient.Clone().Get(url)

	resp, _, err := common.SendRequest(ctx, req)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, json.RawMessage(resp))
}

func WebGetOpenedBlankFile(c echo.Context) (err error) {
	var (
		ctx = common.NewContext(c)
	)

	fileStorageID := c.Param("fileStorageId")
	if fileStorageID == "" {
		return fmt.Errorf("empty fileStorageId param")
	}

	url := config.C.FapiAdapter.Urls.Blank + fmt.Sprintf("/%s", fileStorageID) + fmt.Sprintf("/file")

	req := faClient.Clone().Get(url)

	resp, _, err := common.SendRequest(ctx, req)
	if err != nil {
		return
	}

	return c.Blob(http.StatusOK, "application/pdf", resp)
}

func WebGetBlankPDF(c echo.Context) (err error) {
	var (
		ctx = common.NewContext(c)
	)

	fileStorageID := c.Param("fileStorageId")
	if fileStorageID == "" {
		return fmt.Errorf("empty fileStorageId param")
	}

	url := config.C.FapiAdapter.Urls.Blank + fmt.Sprintf("/%s", fileStorageID)

	req := faClient.Clone().Get(url)

	resp, httpInternals, err := common.SendRequest(ctx, req)
	if err != nil {
		return
	}

	c.Response().Before(func() {
		c.Response().Header().Set(echo.HeaderContentDisposition, httpInternals.Header.Get("Content-Disposition"))
	})

	return c.Blob(http.StatusOK, "application/pdf", resp)
}

func WebGetBlankZIP(c echo.Context) (err error) {
	var (
		ctx = common.NewContext(c)
	)

	queryString := c.QueryString()
	if queryString == "" {
		return fmt.Errorf("empty fileStorageIds query")
	}

	url := config.C.FapiAdapter.Urls.BlanksZip + fmt.Sprintf("?%s", queryString)

	req := faClient.Clone().Get(url)

	resp, httpInternals, err := common.SendRequest(ctx, req)
	if err != nil {
		return
	}

	c.Response().Before(func() {
		c.Response().Header().Set(echo.HeaderContentDisposition, httpInternals.Header.Get("Content-Disposition"))
	})

	return c.Stream(http.StatusOK, "application/octet-stream", bytes.NewReader(resp))
}

func WebGetBlankAllPDF(c echo.Context) (err error) {
	var (
		ctx = common.NewContext(c)
	)

	queryString := c.QueryString()
	if queryString == "" {
		return fmt.Errorf("empty fileStorageIds query")
	}

	url := config.C.FapiAdapter.Urls.BlanksPdf + fmt.Sprintf("?%s", queryString)

	req := faClient.Clone().Get(url)

	resp, httpInternals, err := common.SendRequest(ctx, req)
	if err != nil {
		return
	}

	c.Response().Before(func() {
		c.Response().Header().Set(echo.HeaderContentDisposition, httpInternals.Header.Get("Content-Disposition"))
	})

	return c.Blob(http.StatusOK, "application/pdf", resp)
}
