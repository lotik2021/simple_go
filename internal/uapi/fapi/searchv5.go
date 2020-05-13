package fapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"bitbucket.movista.ru/maas/maasapi/internal/ratelimit"
	"bitbucket.movista.ru/maas/maasapi/internal/uapi/place"
	"github.com/labstack/echo/v4"
)

func SearchAsyncV5(ctx common.Context, req AsyncRequestV5) (*AsyncResponseV5, error) {

	resp, err := common.UapiPost(ctx, faClient, req, config.C.FapiAdapter.Urls.SearchAsyncV5)
	if err != nil {
		return nil, err
	}

	asyncResp := &AsyncResponseV5{}

	err = json.Unmarshal(resp.Data, asyncResp)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal UAPI response - %v", err)
	}

	return asyncResp, nil
}

func GetSearchStatusV5(ctx common.Context, uid string) (interface{}, error) {
	url := config.C.FapiAdapter.Urls.GetSearchStatusV5 + fmt.Sprintf("?uid=%s", uid)
	req := faClient.Clone().Get(url)

	resp, _, err := common.SendRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	return json.RawMessage(resp), nil
}

func GetSearchResultsV5(ctx common.Context, uid string) (interface{}, error) {
	url := config.C.FapiAdapter.Urls.GetSearchResultsV5 + fmt.Sprintf("?uid=%s", uid)

	req := faClient.Clone().Get(url)

	resp, _, err := common.SendRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	var b map[string]interface{}

	err = json.Unmarshal(resp, &b)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal UAPI response - %w", err)
	}

	if _, ok := b["data"]; !ok {
		return nil, fmt.Errorf("no data")
	}

	searchObject, ok := b["data"].(map[string]interface{})["pathGroups"].([]interface{})
	if ok && len(searchObject) > 0 {
		for _, object := range searchObject {
			delete(object.(map[string]interface{}), "services")
			delete(object.(map[string]interface{}), "segments")
			delete(object.(map[string]interface{}), "trips")
		}
	}

	return b, nil
}

func GetPathGroupV5(ctx common.Context, uid, pathGroupId string) (interface{}, error) {
	url := config.C.FapiAdapter.Urls.GetPathGroupV5 + fmt.Sprintf("?uid=%s", uid) + fmt.Sprintf("&pathGroupId=%s", pathGroupId)

	req := faClient.Clone().Get(url)

	resp, _, err := common.SendRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	return json.RawMessage(resp), nil
}

func GetSegmentRoutesV5(ctx common.Context, body []byte) (interface{}, error) {
	url := config.C.FapiAdapter.Urls.GetSegmentRoutesV5

	req := faClient.Clone().Post(url).SendStruct(json.RawMessage(body))

	resp, _, err := common.SendRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	return json.RawMessage(resp), nil

}

func SaveSelectedRoutesV5(c echo.Context) error {

	ctx := common.NewContext(c)

	err := ratelimit.Apply(ctx, ratelimit.MethodSaveSelectedRoutes)
	if err != nil {
		return err
	}

	rawResp, err := common.UapiPost(ctx, faClient, c.Request().Body, config.C.FapiAdapter.Urls.SaveSelectedRoutesV5)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, rawResp)
}

func SearchSyncV5(ctx common.Context, origin, destination models.GeoPoint, departure string,
	tripTypes map[string]bool) (*SearchResult, string, error) {

	placeIdFrom, _, err := place.FindOneByLocation(ctx, &origin)
	if err != nil {
		return nil, "", err
	}
	placeIdTo, _, err := place.FindOneByLocation(ctx, &destination)
	if err != nil {
		return nil, "", err
	}

	const RefreshTime = 1 * time.Second

	req := AsyncRequestV5{
		SearchParams: SearchParamsV5{
			CurrencyCode:    "RUB",
			CultureCode:     "ru",
			DepartureBegin:  departure,
			From:            placeIdFrom,
			To:              placeIdTo,
			Customers:       []Customers{{Age: 36, SeatRequired: true}},
			TripTypes:       unmapTripTypes(tripTypes),
			AirServiceClass: "Economy",
		},
	}

	asyncData, err := SearchAsyncV5(ctx, req)
	if err != nil {
		return nil, "", err
	}

	// getStatus Loop until IsComplete = true

	getStatusUrl := config.C.FapiAdapter.Urls.GetSearchStatusV5 + fmt.Sprintf("?uid=%s", asyncData.UID)

	for {
		time.Sleep(RefreshTime)

		resp, err := common.UapiGet(ctx, faClient, getStatusUrl)
		if err != nil {
			return nil, "", err
		}

		err = json.Unmarshal(resp.Data, &asyncData)
		if err != nil {
			return nil, "", fmt.Errorf("cannot unmarshal UAPI response - %w", err)
		}

		if asyncData.IsComplete {
			break
		}
	}

	searchResultUrl := config.C.FapiAdapter.Urls.GetSearchResultsV5 + fmt.Sprintf("?uid=%s", asyncData.UID)

	resp, err := common.UapiGet(ctx, faClient, searchResultUrl)
	if err != nil {
		return nil, "", err
	}

	searchData := &SearchResult{}

	err = json.Unmarshal(resp.Data, searchData)
	if err != nil {
		return nil, "", fmt.Errorf("cannot unmarshal UAPI response - %w", err)
	}

	return searchData, asyncData.UID, nil
}

func GetTrainSchedule(ctx common.Context, origin, destination models.GeoPoint, departure string, count int) (interface{}, error) {

	searchData, uid, err := SearchSyncV5(ctx, origin, destination, departure, map[string]bool{MOVISTA_TRAIN: true})
	if err != nil {
		return nil, err
	}

	var pathGroupID string

	for _, pgids := range searchData.PathGroups {
		pathGroupID = pgids.ID
	}

	// getPathGroup
	pathGroupUrl := config.C.FapiAdapter.Urls.GetPathGroupV5 + fmt.Sprintf("?uid=%s", uid) + fmt.Sprintf("&pathGroupId=%s", pathGroupID)

	resp, err := common.UapiGet(ctx, faClient, pathGroupUrl)
	if err != nil {
		return nil, err
	}

	var schedule PathGroupResponse
	err = json.Unmarshal(resp.Data, &schedule)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal UAPI response - %w", err)
	}

	scheduleSlice := make([]*TrainSchedule, 0)

	for _, v := range schedule.Trips {
		if v.TripTrain != nil && v.TripTrain.Departure > departure {
			scheduleSlice = append(scheduleSlice, v.TripTrain)
		}
		if v.TripTrainSuburban != nil && v.TripTrainSuburban.Departure > departure {
			scheduleSlice = append(scheduleSlice, v.TripTrainSuburban)
		}
	}
	sort.Slice(scheduleSlice, func(i int, j int) bool {
		return scheduleSlice[i].Departure < scheduleSlice[j].Departure
	})

	if count == 0 {
		count = 10 // default value
	}
	scheduleSlice = scheduleSlice[:count]
	return scheduleSlice, nil
}

func GetSelectedRoutesV5(c echo.Context) error {
	var (
		req struct {
			UID string `json:"uid"`
		}
		ctx = common.NewContext(c)
	)

	err := common.BindAndValidateReq(c, &req)
	if err != nil {
		return err
	}

	url := config.C.FapiAdapter.Urls.GetSelectedRoutesV5 + fmt.Sprintf("?uid=%s", req.UID)
	rawResp, err := common.UapiGet(ctx, faClient, url)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, rawResp)
}
