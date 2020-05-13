package travelcard

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
)

type TravelCard struct {
	ID                int      `json:"id"`
	Image             string   `json:"image"`
	IosIconUrl        string   `json:"ios_icon_url"`
	AndroidIconUrl    string   `json:"android_icon_url"`
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	Warn              string   `json:"warn"`
	MaximumAmount     *float64 `json:"maximum_amount,omitempty"`
	MinimumAmount     *float64 `json:"minimum_amount,omitempty"`
	MinimumCardLength float64  `json:"minimum_card_length"`
	MaximumCardLength float64  `json:"maximum_card_length"`
	Balance           *int     `json:"balance,omitempty"`
}

func GetTravelCards(ctx common.Context) (travelCards []TravelCard, err error) {
	travelCards = make([]TravelCard, 2)

	troikaCardInfo := GetTroikaCardInfo()
	troikaParameters, err := GetTroikaParameters(ctx)
	if err != nil {
		return
	}
	troikaCardInfo.MinimumAmount = troikaParameters.MinimumAmount
	troikaCardInfo.MaximumAmount = troikaParameters.MaximumAmount
	troikaCardInfo.MinimumCardLength = troikaParameters.CardMinimumLength
	troikaCardInfo.MaximumCardLength = troikaParameters.CardMaximumLength
	travelCards[0] = troikaCardInfo

	strelkaCardInfo := GetStrelkaCardInfo()
	strelkaParameters, err := GetStrelkaParameters(ctx)
	if err != nil {
		return
	}
	strelkaCardInfo.MinimumCardLength = strelkaParameters.CardMinimumLength
	strelkaCardInfo.MaximumCardLength = strelkaParameters.CardMaximumLength
	travelCards[1] = strelkaCardInfo

	return
}
