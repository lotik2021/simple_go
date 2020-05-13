package travelcard

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"encoding/json"
	"fmt"
	"strings"
)

const (
	ErrorTroikaApiRequest  = "5000"
	ErrorTroikaApiResponse = "5001"
)

type TroikaParameters struct {
	MinimumAmount     *float64 `json:"minimum_amount"`
	MaximumAmount     *float64 `json:"maximum_amount"`
	CardMinimumLength float64  `json:"minimum_card_length"`
	CardMaximumLength float64  `json:"maximum_card_length"`
}

func completeTroikaError(code string, err error) error {
	resultCode := ErrorTroikaApiRequest
	if code != "" {
		resultCode = code
	}

	return models.Error{
		Code:    resultCode,
		Message: err.Error(),
	}
}

func GetTroikaParameters(ctx common.Context) (parameters TroikaParameters, err error) {

	var resp struct {
		Form []struct {
			Items []struct {
				MinLength float64 `json:"minlength"`
				MaxLength float64 `json:"maxlength"`
				Min       float64 `json:"min"`
				Max       float64 `json:"max"`
			} `json:"items"`
		} `json:"form"`
	}

	spanReq := common.DefaultRequest.Clone().Get(config.C.YandexMoney.Urls.TroikaParameters)

	body, _, _ := common.SendRequest(ctx, spanReq)

	err = json.Unmarshal(body, &resp)
	if err != nil {
		err = completeTroikaError(ErrorTroikaApiResponse, fmt.Errorf("cannot unmarshall troika response - %s", string(body)))
		return
	}

	//m := objx.MustFromJSON(string(body))
	//cardMinimumLength := m.Get("form[0].items[0].minlength").Float64()
	//cardMaximumLength := m.Get("form[0].items[0].maxlength").Float64()
	//minimumAmount := m.Get("form[0].items[1].min").Float64()
	//maximumAmount := m.Get("form[0].items[1].max").Float64()

	cardMinimumLength := resp.Form[0].Items[0].MinLength
	cardMaximumLength := resp.Form[0].Items[0].MaxLength
	minimumAmount := resp.Form[0].Items[1].Min
	maximumAmount := resp.Form[0].Items[1].Max

	if cardMinimumLength == 0 || cardMaximumLength == 0 || minimumAmount == 0 || maximumAmount == 0 {
		err = completeStrelkaError(ErrorTroikaApiResponse, fmt.Errorf("cannot get minlength or maxlength or min or max from response - %s", string(body)))
		return
	}

	parameters = TroikaParameters{
		MinimumAmount:     &minimumAmount,
		MaximumAmount:     &maximumAmount,
		CardMinimumLength: cardMinimumLength,
		CardMaximumLength: cardMaximumLength,
	}

	return
}

func PayTroika(ctx common.Context, amount int, cardNumber, paymentId string) (redirectURL string, err error) {

	instanceIDRequest := common.DefaultRequest.Clone().Post(config.C.Troika.Urls.RequestInstanceID).
		Type("urlencoded").
		Send(fmt.Sprintf(`client_id=%s`, config.C.YandexMoney.ClientID))

	var instanceIDResponse struct {
		Status     string `json:"status"`
		InstanceID string `json:"instance_id"`
		Error      string `json:"error"`
	}

	body, _, err := common.SendRequest(ctx, instanceIDRequest)
	if err != nil {
		err = completeTroikaError(ErrorTroikaApiRequest, err)
		return
	}

	err = json.Unmarshal(body, &instanceIDResponse)
	if err != nil {
		err = completeTroikaError(ErrorTroikaApiResponse, fmt.Errorf("cannot unmarshall troika response - %s", string(body)))
		return
	}

	if instanceIDResponse.Status != "success" && instanceIDResponse.Status != "ext_auth_required" {
		err = completeTroikaError(ErrorTroikaApiResponse, fmt.Errorf("error response from troika - %s", string(body)))
		return
	}

	var externalPaymentResponse struct {
		Status    string `json:"status"`
		RequestID string `json:"request_id"`
		Error     string `json:"error"`
	}

	sum := fmt.Sprintf("%.2f", float64(amount/100))
	externalPaymentRequestBody := fmt.Sprintf(`pattern_id=%s&instance_id=%s&sum=%s&customerNumber=%s`,
		config.C.Troika.ClientID, strings.Replace(instanceIDResponse.InstanceID, "+", "%2B", -1), sum, cardNumber)

	externalPaymentRequest := common.DefaultRequest.Clone().Post(config.C.Troika.Urls.ExternalPayment).
		Type("urlencoded").
		Send(externalPaymentRequestBody)

	body, _, err = common.SendRequest(ctx, externalPaymentRequest)
	if err != nil {
		err = completeTroikaError(ErrorTroikaApiRequest, err)
		return
	}

	err = json.Unmarshal(body, &externalPaymentResponse)
	if err != nil {
		err = completeTroikaError(ErrorTroikaApiResponse, fmt.Errorf("cannot unmarshall troika response - %s", string(body)))
		return
	}

	if externalPaymentResponse.Status != "success" && externalPaymentResponse.Status != "ext_auth_required" {
		err = completeTroikaError(ErrorTroikaApiResponse, fmt.Errorf("error response from troika - %s", string(body)))
		return
	}

	var processExternalPaymentResponse struct {
		Status     string            `json:"status"`
		URL        string            `json:"acs_uri"`
		Parameters map[string]string `json:"acs_params"`
		Error      string            `json:"error"`
	}

	processExternalPaymentRequest := common.DefaultRequest.Clone().Post(config.C.Troika.Urls.ProcessExternalPayment).
		Type("urlencoded").
		Send(fmt.Sprintf(`request_id=%s&instance_id=%s&ext_auth_success_uri=%s&ext_auth_fail_uri=%s`,
			externalPaymentResponse.RequestID,
			strings.Replace(instanceIDResponse.InstanceID, "+", "%2B", -1),
			config.C.Troika.Urls.SuccessRedirect+"?paymentId="+paymentId,
			config.C.Troika.Urls.FailRedirect+"?paymentId="+paymentId))

	body, _, err = common.SendRequest(ctx, processExternalPaymentRequest)
	if err != nil {
		err = completeTroikaError(ErrorTroikaApiRequest, err)
		return
	}

	err = json.Unmarshal(body, &processExternalPaymentResponse)
	if err != nil {
		err = completeTroikaError(ErrorTroikaApiResponse, fmt.Errorf("cannot unmarshall troika response - %s", string(body)))
		return
	}

	if processExternalPaymentResponse.Status != "success" && processExternalPaymentResponse.Status != "ext_auth_required" {
		err = completeTroikaError(ErrorTroikaApiResponse, fmt.Errorf("error response from troika - %s", string(body)))
		return
	}

	redirectURL = processExternalPaymentResponse.URL + "?"

	for key, value := range processExternalPaymentResponse.Parameters {
		redirectURL = redirectURL + key + "=" + value + "&"
	}

	return
}

func GetTroikaCardInfo() TravelCard {
	return TravelCard{
		ID:             0,
		Image:          "icon_travel_card_troika",
		Name:           "Тройка",
		Description:    "Для наземного транспорта и метро",
		Warn:           "После пополнения «Тройки» в приложении, нужно подойти к желтому терминалу в метро, нажать на нем «Удаленное пополнение» и поднести Тройку к ридеру. Не убирайте Тройку, пока не появится сообщение, что поездки записаны.",
		IosIconUrl:     config.C.Icons.TravelCards.Troika.Ios,
		AndroidIconUrl: config.C.Icons.TravelCards.Troika.Android,
	}
}
