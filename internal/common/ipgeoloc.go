package common

import (
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"encoding/json"
	"net/http"
)

func FindPlaceTitleByIP(ctx Context, url string) (placeTitle string, err error) {
	serviceResp := struct {
		City string `json:"city"`
	}{}

	req := DefaultRequest.Clone().Get(url)

	body, httpInternals, err := SendRequest(ctx, req)
	if err != nil || httpInternals.StatusCode == http.StatusTooManyRequests || body == nil {
		err = models.NewFindPlaceError()
		return
	}

	err = json.Unmarshal(body, &serviceResp)
	if err != nil || len(serviceResp.City) == 0 {
		err = models.NewFindPlaceError()
		return
	}

	placeTitle = serviceResp.City

	return
}
