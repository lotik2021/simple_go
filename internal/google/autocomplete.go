package google

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"context"
	"github.com/google/uuid"
	"googlemaps.github.io/maps"
	"sync"
)

func FindAutocomplete(ctx common.Context, sessionToken *uuid.UUID, input string, region string, location *models.GeoPoint) (result []Place, err error) {

	result = make([]Place, 0)

	if *sessionToken == uuid.Nil {
		*sessionToken = uuid.UUID(maps.NewPlaceAutocompleteSessionToken())
	}

	autocompleteReqWithoutSpace := maps.PlaceAutocompleteRequest{
		Input:        input,
		Language:     region,
		Radius:       20000,
		SessionToken: maps.PlaceAutocompleteSessionToken(*sessionToken),
	}

	autocompleteReqWithSpace := autocompleteReqWithoutSpace
	if location != nil {
		autocompleteReqWithSpace.Location = &maps.LatLng{Lat: location.Latitude, Lng: location.Longitude}
	}

	autocompleteReqWithSpace.Input = autocompleteReqWithSpace.Input + " "

	wg := &sync.WaitGroup{}
	wg.Add(2)

	var searchResultsWithSpace, searchResultsWithoutSpace maps.AutocompleteResponse

	go func(result *maps.AutocompleteResponse, req maps.PlaceAutocompleteRequest) {
		defer wg.Done()
		response, err := mapsClient.PlaceAutocomplete(context.Background(), &req)
		if err != nil {
			return
		}
		*result = response
	}(&searchResultsWithSpace, autocompleteReqWithoutSpace)

	go func(result *maps.AutocompleteResponse, req maps.PlaceAutocompleteRequest) {
		defer wg.Done()
		response, err := mapsClient.PlaceAutocomplete(context.Background(), &req)
		if err != nil {
			return
		}
		*result = response
	}(&searchResultsWithSpace, autocompleteReqWithSpace)

	wg.Wait()

	if len(searchResultsWithoutSpace.Predictions) == 0 && len(searchResultsWithSpace.Predictions) == 0 {
		return
	}

	result = appendAutocompleteResult(searchResultsWithoutSpace, searchResultsWithSpace)

	return
}

func appendAutocompleteResult(a maps.AutocompleteResponse, b maps.AutocompleteResponse) (result []Place) {
	exists := make(map[string]bool)
	filteredPredictions := make([]maps.AutocompletePrediction, 0)
	result = make([]Place, 0)

	for _, val := range a.Predictions {
		filteredPredictions = append(filteredPredictions, val)
		exists[val.PlaceID] = true
	}

	for _, val := range b.Predictions {
		if _, ok := exists[val.PlaceID]; !ok {
			filteredPredictions = append(filteredPredictions, val)
		}
	}

	for _, prediction := range filteredPredictions {

		item := Place{
			MainText:      prediction.StructuredFormatting.MainText,
			SecondaryText: prediction.StructuredFormatting.SecondaryText,
			ID:            prediction.PlaceID,
			PlaceTypes:    prediction.Types,
		}

		item.GetIcons()

		result = append(result, item)
	}

	return result
}
