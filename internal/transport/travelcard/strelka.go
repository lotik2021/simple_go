package travelcard

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"encoding/json"
	"fmt"
)

const (
	ErrorStrelkaApiRequest                 = "6000"
	ErrorStrelkaApiResponse                = "6001"
	ErrorStrelkaApiResponseCardBlocked     = "6002"
	ErrorStrelkaApiResponseEmptyPaymentURL = "6003"
)

type StrelkaParameters struct {
	MinimumAmount     *float64 `json:"minimum_amount"`
	MaximumAmount     *float64 `json:"maximum_amount"`
	CardMinimumLength float64  `json:"minimum_card_length"`
	CardMaximumLength float64  `json:"maximum_card_length"`
}

func completeStrelkaError(code string, err error) error {
	resultCode := ErrorStrelkaApiRequest
	if code != "" {
		resultCode = code
	}

	return models.Error{
		Code:    resultCode,
		Message: err.Error(),
	}
}

func GetStrelkaParameters(ctx common.Context) (parameters StrelkaParameters, err error) {
	parametersRequest := common.DefaultRequest.Clone().Get(config.C.YandexMoney.Urls.StrelkaParameters)

	body, _, _ := common.SendRequest(ctx, parametersRequest)

	var resp struct {
		Form []struct {
			Items []struct {
				MinLength float64 `json:"minlength"`
				MaxLength float64 `json:"maxlength"`
			} `json:"items"`
		} `json:"form"`
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		err = completeStrelkaError(ErrorStrelkaApiResponse, fmt.Errorf("cannot unmarshall strelka response - %s", string(body)))
		return
	}

	//m := objx.MustFromJSON(string(body))
	//cardMinimumLength := m.Get("form[0].items[0].minlength").Float64()
	//cardMaximumLength := m.Get("form[0].items[0].maxlength").Float64()

	cardMinimumLength := resp.Form[0].Items[0].MinLength
	cardMaximumLength := resp.Form[0].Items[0].MaxLength

	if cardMinimumLength == 0 || cardMaximumLength == 0 {
		err = completeStrelkaError(ErrorStrelkaApiResponse, fmt.Errorf("cannot get minlength or maxlength from response - %s", string(body)))
		return
	}

	parameters = StrelkaParameters{
		CardMinimumLength: cardMinimumLength,
		CardMaximumLength: cardMaximumLength,
	}

	return
}

func GetStrelkaBalance(ctx common.Context, cardNumber string) (balance float64, cardTypeID string, redirectURL string, parameters StrelkaParameters, err error) {
	if !config.C.Strelka.ApiAvailable {
		redirectURL = config.C.Strelka.Urls.DefaultRedirect + cardNumber
		return
	}

	type strelkaType struct {
		MinimumPayment float64 `json:"cardtypepaymin"`
		MaximumBalance float64 `json:"cardtypemaxbalance"`
		ID             string  `json:"cardtypeid"`
		Name           string  `json:"cardtypename"`
		NFCID          string  `json:"cardtypenfcid"`
		Description    string  `json:"cardtypedesc"`
		MaximumPayment float64 `json:"cardtypepaymax"`
		IsTicket       bool    `json:"cardtypeisticket"`
		BaseRate       int64   `json:"cardtypebaserate"`
		Sectors        string  `json:"cardtypesectors"`
		IsNFC          bool    `json:"cardtypeisnfc"`
		IsWallet       bool    `json:"cardtypeiswallet"`
		Code           string  `json:"cardtypecode"`
	}
	response := make([]strelkaType, 0)

	spanReq := common.DefaultRequest.Clone().Get(config.C.Strelka.Urls.Types)

	body, _, err := common.SendRequest(ctx, spanReq)
	if err != nil {
		err = completeStrelkaError(ErrorStrelkaApiRequest, err)
		return
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		err = fmt.Errorf("cannot unmarshal strelka response - %s", string(body))
		return
	}

	if len(response) == 0 {
		err = completeStrelkaError(ErrorStrelkaApiResponse, fmt.Errorf("zero response length - %s", string(body)))
		return
	}

	cardTypeID = response[0].ID

	if cardTypeID == "" {
		err = completeStrelkaError(ErrorStrelkaApiResponse, fmt.Errorf("empty cardtypeid in response - %s", string(body)))
		return
	}

	var balanceResponse struct {
		Balance      float64  `json:"balance"`
		CardBlocked  bool     `json:"cardblocked"`
		ErrorMessage []string `json:"__all__"`
	}

	balanceReq := common.DefaultRequest.Clone().Get(config.C.Strelka.Urls.Balance).Param("cardnum", cardNumber).Param("cardtypeid", response[0].ID)
	body, _, err = common.SendRequest(ctx, balanceReq)
	if err != nil {
		err = completeStrelkaError(ErrorStrelkaApiRequest, err)
		return
	}

	err = json.Unmarshal(body, &balanceResponse)
	if err != nil {
		err = completeStrelkaError(ErrorStrelkaApiResponse, fmt.Errorf("cannot unmarshal strelka balance response - %s", string(body)))
		return
	}

	if len(balanceResponse.ErrorMessage) > 0 {
		err = completeStrelkaError(ErrorStrelkaApiResponse, fmt.Errorf("%+v", balanceResponse.ErrorMessage))
		return
	}

	if balanceResponse.CardBlocked {
		err = completeStrelkaError(ErrorStrelkaApiResponseCardBlocked, fmt.Errorf("card %s blocked", cardNumber))
		return
	}

	balance = balanceResponse.Balance / 100
	calculatedMaximumBalance := response[0].MaximumBalance/100 - balance
	minAmount := response[0].MinimumPayment / 100

	parameters = StrelkaParameters{
		MinimumAmount: &minAmount,
		MaximumAmount: &calculatedMaximumBalance,
	}

	return
}

func PayStrelka(ctx common.Context, amount int, paymentId, cardNumber, cardTypeID string) (redirectURL string, err error) {
	var payRequest struct {
		CardNumber string `json:"cardnum"`
		CardTypeID string `json:"cardtypeid"`
		Amount     int    `json:"paysum"`
		Redirect   string `json:"returnurl"`
	}

	payRequest.CardTypeID = cardTypeID
	payRequest.CardNumber = cardNumber
	payRequest.Amount = amount
	payRequest.Redirect = config.C.Strelka.Urls.Redirect + "?paymentId=" + paymentId

	spanReq := common.DefaultRequest.Clone().Post(config.C.Strelka.Urls.Pay).Type("json").SendStruct(&payRequest)

	var paymentResponse struct {
		PaymentURL string `json:"payurl"`
		ID         string `json:"payid"`
	}

	body, _, err := common.SendRequest(ctx, spanReq)
	if err != nil {
		err = completeStrelkaError(ErrorStrelkaApiRequest, err)
		return
	}

	err = json.Unmarshal(body, &paymentResponse)
	if err != nil {
		err = completeStrelkaError(ErrorStrelkaApiResponse, fmt.Errorf("cannot unmarshal strelka pay response - %s", string(body)))
		return
	}

	if paymentResponse.PaymentURL == "" {
		err = completeStrelkaError(ErrorStrelkaApiResponseEmptyPaymentURL, fmt.Errorf("empty paymentURL in response - %s", string(body)))
		return
	}

	redirectURL = paymentResponse.PaymentURL

	return
}

func GetStrelkaCardInfo() TravelCard {
	return TravelCard{
		ID:             1,
		Image:          "icon_travel_card_strelka",
		Name:           "Стрелка",
		Description:    "Для наземного транспорта и электричек.",
		Warn:           "Комиссии нет. Картой можно оплатить через 15 минут после пополнения.",
		IosIconUrl:     config.C.Icons.TravelCards.Strelka.Ios,
		AndroidIconUrl: config.C.Icons.TravelCards.Strelka.Android,
	}
}
