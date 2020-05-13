package place

import (
	"encoding/json"
	"fmt"

	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"github.com/parnurzeal/gorequest"
)

var (
	placesClient *gorequest.SuperAgent
)

func init() {
	placesClient = common.DefaultRequest.Clone().Timeout(config.C.Places.RequestTimeout)
}

type PlaceIdsRequest struct {
	PlaceIds []int `json:"placeids" query:"placeids"`
}

type Place struct {
	ID                   int     `json:"id"`
	Name                 string  `json:"name"`
	Lat                  float64 `json:"lat"`
	Lon                  float64 `json:"lon"`
	TimeZone             string  `json:"timeZone"`
	CountryName          string  `json:"countryName"`
	StateName            string  `json:"stateName"`
	CityName             string  `json:"cityName"`
	StationName          string  `json:"stationName"`
	PlatformName         string  `json:"platformName"`
	Description          string  `json:"description"`
	FullName             string  `json:"fullName"`
	PlaceClassId         int     `json:"placeclassId"`
	NearestBiggerPlaceId int     `json:"nearestBiggerPlaceId"`
	TypePlace            int     `json:"typePlace"`
}

func FindByName(ctx common.Context, name string, count int, placeTypes []string) (places []Place, err error) {
	if placeTypes == nil || len(placeTypes) == 0 {
		placeTypes = []string{"city", "station"}
	}

	var (
		reqBody = struct {
			Text         string   `json:"text"`
			Count        int      `json:"count"`
			PlaceClasses []string `json:"placeClasses"`
		}{
			Text:         name,
			Count:        count,
			PlaceClasses: placeTypes,
		}
		resp = struct {
			ResultsCount int     `json:"resultsCount"`
			Places       []Place `json:"places"`
			TimeSpent    string  `json:"timeSpent"`
		}{}
	)

	rawResp, err := common.UapiPost(ctx, placesClient, reqBody,
		config.C.Places.Urls.SearchPlaces)
	if err != nil {
		err = fmt.Errorf("places err %v", err)
		return
	}

	err = json.Unmarshal(rawResp.Data, &resp)
	if err != nil {
		err = fmt.Errorf("cannot unmarshal places response - %s", string(rawResp.Data))
		return
	}

	places = resp.Places

	return
}

func FindOneByLocation(ctx common.Context, location *models.GeoPoint) (placeId int, placeName string, err error) {
	var (
		reqBody = struct {
			Longitude    float64  `json:"longitude"`
			Latitude     float64  `json:"latitude"`
			Radius       int      `json:"radius"`
			Count        int      `json:"count"`
			PlaceClasses []string `json:"placeClasses"`
		}{
			Longitude:    location.Longitude,
			Latitude:     location.Latitude,
			Radius:       100,
			Count:        1,
			PlaceClasses: []string{"city"},
		}
		resp = struct {
			Places []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"places"`
			TimeSpent string `json:"timeSpent"`
		}{}
	)

	rawResp, err := common.UapiPost(ctx, placesClient, reqBody,
		config.C.Places.Urls.SearchPlacesByGeo)
	if err != nil {
		err = fmt.Errorf("places err %w", err)
		return
	}

	err = json.Unmarshal(rawResp.Data, &resp)
	if err != nil {
		err = fmt.Errorf("cannot unmarshal places response - %s", string(rawResp.Data))
		return
	}

	if len(resp.Places) == 0 {
		err = fmt.Errorf("empty response from places")
		return
	}

	placeId = resp.Places[0].ID
	placeName = resp.Places[0].Name

	return
}

func FindByID(ctx common.Context, id int) (*Place, error) {
	places, err := FindByIds(ctx, []int{id})
	if err != nil {
		return nil, err
	}
	if place, ok := places[id]; ok {
		return &place, nil
	}
	return nil, fmt.Errorf("Не удалось найти placeID в ответе от FindByIds")
}

func FindByIds(ctx common.Context, ids []int) (places map[int]Place, err error) {
	var (
		req = PlaceIdsRequest{
			PlaceIds: ids,
		}

		resp = map[string]Place{}
	)

	rawResp, err := PlacesByIds(ctx, req)
	if err != nil {
		return
	}

	err = json.Unmarshal(rawResp.Data, &resp)
	if err != nil {
		return
	}

	places = make(map[int]Place, 0)

	for _, v := range resp {
		places[v.ID] = v
	}

	return
}

func PlacesByIds(ctx common.Context, req PlaceIdsRequest) (*models.RawResponse, error) {
	return common.UapiAuthorizedPost(ctx, placesClient, req, config.C.Places.Urls.SearchPlacesByIds)
}
