package fapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"bitbucket.movista.ru/maas/maasapi/internal/constant"
	"bitbucket.movista.ru/maas/maasapi/internal/dictionary"
	"bitbucket.movista.ru/maas/maasapi/internal/logger"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"bitbucket.movista.ru/maas/maasapi/internal/notebook"
	"bitbucket.movista.ru/maas/maasapi/internal/session"
	"bitbucket.movista.ru/maas/maasapi/internal/uapi/auth"

	"time"

	"github.com/labstack/echo/v4"
)

func createDeviceMapping(resp *models.RawResponse, ctx common.Context) error {

	bookResp := &CreateBookingResponse{}
	err := json.Unmarshal(resp.Data, bookResp)
	if err != nil {
		return err
	}

	// сохраняем соответствие device_id -> order_id
	err = session.Create(ctx, bookResp.Order.UserID)
	if err != nil {
		return err
	}

	return nil
}

func checkDeviceMapping(resp *models.RawResponse, ctx common.Context) error {
	bookResp := &CreateBookingResponse{}
	err := json.Unmarshal(resp.Data, bookResp)
	if err != nil {
		return err
	}

	// находим связку device_id -> user_id
	_, err = session.Find(ctx, bookResp.Order.UserID)
	return err
}

func CreateBookingV2(c echo.Context) error {
	var (
		ctx = common.NewContext(c)
	)

	rawResp, err := common.UapiAuthorizedPost(ctx, faClient, c.Request().Body, config.C.Booking.Urls.CreateBookingV2)
	if err != nil {
		return err
	}

	err = createDeviceMapping(rawResp, ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, rawResp)
}

func CreateBookingV4(c echo.Context) (err error) {
	var (
		ctx = common.NewContext(c)
	)

	resp, err := common.UapiAuthorizedPost(ctx, faClient, c.Request().Body, config.C.FapiAdapter.Urls.CreateBookingV4)
	if err != nil {
		return err
	}

	err = createDeviceMapping(resp, ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func GetOrderV2(c echo.Context) (err error) {
	var (
		ctx = common.NewContext(c)
	)

	resp, err := common.UapiAuthorizedPost(ctx, faClient, c.Request().Body, config.C.FapiAdapter.Urls.GetOrderV2)
	if err != nil {
		return err
	}

	err = checkDeviceMapping(resp, ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func GetOrderMobileV2(c echo.Context) (err error) {
	var (
		ctx = common.NewContext(c)
	)

	rawResp, err := common.UapiPost(ctx, faClient, c.Request().Body, config.C.FapiAdapter.Urls.GetOrderByEmail)
	if err != nil {
		return
	}

	err = checkDeviceMapping(rawResp, ctx)
	if err != nil {
		return err
	}

	maasOrderResp, err := orderToMaasResp(ctx, rawResp)
	if err != nil {
		return
	}

	if c.Request().Header.Get("mode") == "debug" {
		return c.JSON(http.StatusOK, echo.Map{"data": maasOrderResp, "debug": rawResp})
	}

	return c.JSON(http.StatusOK, echo.Map{"data": maasOrderResp})
}

func CreateOrderV4(c echo.Context) (err error) {
	// этот метод объединяет CreateBooking и GetOrder в один метод

	// логика работы:
	// if CreateBooking() returned ERROR --> return ERROR
	// return GetOrder()

	var (
		in struct {
			Booking interface{} `json:"booking"`
			Order   interface{} `json:"order"`
		}
		ctx = common.NewContext(c)
	)

	if err = common.BindAndValidateReq(c, &in); err != nil {
		return err
	}

	// посылаем запрос на CreateBooking
	resp, err := common.UapiAuthorizedPost(ctx, faClient, in.Booking, config.C.FapiAdapter.Urls.CreateBookingV4)
	if err != nil {
		return err
	}

	bookErrs := resp.Error

	// сохраняем соответствие device_id -> user_id
	err = createDeviceMapping(resp, ctx)
	if err != nil {
		return err
	}

	var order struct {
		Order struct {
			ID int `json:"id"`
		} `json:"order"`
	}

	err = json.Unmarshal(resp.Data, &order)
	if err != nil {
		return err
	}

	orderPayload, ok := in.Order.(map[string]interface{})
	if !ok {
		return fmt.Errorf("cannot cast GetOrder payload to map[string]interface{}\n")
	}

	orderPayload["orderId"] = order.Order.ID

	// посылаем запрос на GetOrder
	resp, err = common.UapiAuthorizedPost(ctx, faClient, orderPayload, config.C.FapiAdapter.Urls.GetOrderV2)
	if err != nil {
		return err
	}

	// проверяем, что user_id совпадают
	err = checkDeviceMapping(resp, ctx)
	if err != nil {
		return err
	}

	maasOrderResp, err := orderToMaasResp(ctx, resp)
	if err != nil {
		return
	}

	// если были ошибки при вызове createBooking, добавляем их в ответ
	maasOrderResp.Error = append(maasOrderResp.Error, bookErrs...)

	if c.Request().Header.Get("mode") == "debug" {
		return c.JSON(http.StatusOK, echo.Map{"data": maasOrderResp, "debug": resp})
	}

	return c.JSON(http.StatusOK, echo.Map{"data": maasOrderResp})
}

func AddServicesToOrder(c echo.Context) (err error) {
	var (
		body []byte
		ctx  = common.NewContext(c)
	)

	body, err = ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	req := faClient.Clone().Post(config.C.FapiAdapter.Urls.AddServicesToOrder).SendStruct(json.RawMessage(body))

	resp, _, err := common.SendRequest(ctx, req)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, json.RawMessage(resp))
}

func RefundOrderAsync(ctx common.Context, jobID string) (interface{}, error) {
	url := config.C.FapiAdapter.Urls.RefundOrderAsyncV2 + fmt.Sprintf("/%s", jobID)

	req := faClient.Clone().Get(url)

	resp, _, err := common.SendRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	return json.RawMessage(resp), nil
}

func CheckRefund(c echo.Context) (err error) {
	var (
		body []byte
		ctx  = common.NewContext(c)
	)

	body, err = ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return
	}

	req := faClient.Clone().Post(config.C.FapiAdapter.Urls.CheckRefund).SendStruct(json.RawMessage(body))

	resp, _, err := common.SendRequest(ctx, req)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, json.RawMessage(resp))
}

func CheckBookingV4(ctx common.Context, uid string) (*models.RawResponse, error) {
	url := config.C.FapiAdapter.Urls.CheckBookingV4 + fmt.Sprintf("/%s", uid)

	return common.UapiGet(ctx, faClient, url)
}

func CheckBookingMobileV4(c echo.Context) (err error) {
	var (
		in struct {
			UID string `json:"uid" validate:"required"`
		}
		ctx = common.NewContext(c)
	)

	if err = common.BindAndValidateReq(c, &in); err != nil {
		return err
	}

	url := config.C.FapiAdapter.Urls.CheckBookingV4 + fmt.Sprintf("/%s", in.UID)

	rawResp, err := common.UapiGet(ctx, faClient, url)
	if err != nil {
		return
	}

	var (
		checkBookingResp  CheckBookingResponse
		tripIDToServiceID TripIDToServiceID
		route             CheckBookingMobile
		docs              []notebook.Customer
	)

	err = json.Unmarshal(rawResp.Data, &checkBookingResp)
	if err != nil {
		ctx.Logger.WithError(err).WithField("response", string(rawResp.Data)).Error("cannot unmarshal UAPI response")
	}

	tripsmp := make(map[string]MaasTrip)
	servicesmp := make(map[string]interface{})

	for _, s := range checkBookingResp.Order.Services {
		for _, tid := range s.TripIds {
			tripIDToServiceID.ServiceID = s.ID
			tripIDToServiceID.TripID = tid
			for _, trp := range checkBookingResp.Order.Trips {
				optionsmp := make(map[string]string)
				if trp.ID == tid && !trp.IsReturnTrip {
					route.Forward = append(route.Forward, tripIDToServiceID)
					for _, contData := range s.ContainerData.TripsData {
						for _, opt := range contData.GeneralData.Options {
							optionsmp[opt.Code] = opt.Name
						}
					}
					trp.Options = optionsmp
					tripsmp[tid] = trp
					servicesmp[s.ID] = s
				}
				if trp.ID == tid && trp.IsReturnTrip {
					route.Backward = append(route.Backward, tripIDToServiceID)
					for _, contData := range s.ContainerData.TripsData {
						for _, opt := range contData.GeneralData.Options {
							optionsmp[opt.Code] = opt.Name
						}
					}
					trp.Options = optionsmp
					tripsmp[tid] = trp
					servicesmp[s.ID] = s
				}
			}
		}
	}
	// чтобы не пробрасывать поля с null'ами
	if len(tripsmp) != 0 {
		route.Trips = tripsmp
	}
	if len(servicesmp) != 0 {
		route.Services = servicesmp
		route.Customers = checkBookingResp.Order.Customers
		route.PriceInfo = checkBookingResp.Order.PriceInfo
		route.OrderStatus = checkBookingResp.Order.OrderStatus
		route.Owner = checkBookingResp.Order.Owner
		route.BeginDate = checkBookingResp.Order.BeginDate
		route.Customers = checkBookingResp.Order.Customers
		route.PriceInfo = checkBookingResp.Order.PriceInfo
		route.OrderStatus = checkBookingResp.Order.OrderStatus
		route.Owner = checkBookingResp.Order.Owner
		route.BeginDate = checkBookingResp.Order.BeginDate
		route.Places = checkBookingResp.Places
	}

	citShips, err := dictionary.GetCitizenships(ctx, "ru")
	if err != nil {
		logger.Log.Errorf("cannot get user citizenship - %w", err)
	}
	route.Citizenship = citShips

	if ctx.IsUser() {
		userInfo, _ := auth.GetMainUserInfo(ctx)

		if userInfo != nil {
			route.UserInfo = userInfo
		}

		resp, err := common.UapiAuthorizedRequest(ctx, faClient, nil, http.MethodGet, config.C.PersonalArea.Urls.DocumentsV2, nil)
		if err != nil {
			logger.Log.Errorf("cannot get user documents - %w", err)
		}

		if resp.Data != nil {
			err = json.Unmarshal(resp.Data, &docs)
			if err != nil {
				logger.Log.Errorf("cannot unmarshal document - %w", err)
			}
			route.CustomerDocuments = docs
		}
	}

	route.Errors = rawResp.Error

	return c.JSON(http.StatusOK, echo.Map{"data": route})
}

func GetOrderByEmailV2(c echo.Context) (err error) {
	var (
		body []byte
		ctx  = common.NewContext(c)
	)

	body, err = ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return
	}

	req := faClient.Clone().Post(config.C.FapiAdapter.Urls.GetOrderByEmail).SendStruct(json.RawMessage(body))

	resp, _, err := common.SendRequest(ctx, req)
	if err != nil {
		return
	}

	client := c.Request().Header.Get("device-type")

	if client == constant.WEB || client == constant.WEBMOBILE {
		cquid := c.Request().Header.Get("carrotquest-uid")
		logger.Log.Infof("carrot quest uid: %s", cquid)

		if cquid != "" {
			reqToPush := faClient.Clone().Post(config.C.Push.Urls.CarrotQuestOrderEvent).
				SendMap(echo.Map{"data": json.RawMessage(resp), "user_id": ctx.UserID, "carrot_quest_id": cquid, "created": time.Now().Unix()})

			go common.SendRequest(ctx, reqToPush)
		} else {
			logger.Log.Infof("carrot quest uid empty, no info to send")
		}
	}

	return c.JSON(http.StatusOK, json.RawMessage(resp))
}

func GetOrderByEmailMobileV2(c echo.Context) (err error) {
	var (
		ctx = common.NewContext(c)
	)

	rawResp, err := common.UapiPost(ctx, faClient, c.Request().Body, config.C.FapiAdapter.Urls.GetOrderByEmail)
	if err != nil {
		return
	}

	maasOrderResp, err := orderToMaasResp(ctx, rawResp)
	if err != nil {
		return
	}

	if c.Request().Header.Get("mode") == "debug" {
		return c.JSON(http.StatusOK, echo.Map{"data": maasOrderResp, "debug": rawResp})
	}

	return c.JSON(http.StatusOK, echo.Map{"data": maasOrderResp})
}

func CheckBookingAviaBus(c echo.Context) (err error) {
	var (
		ctx = common.NewContext(c)
		in  struct {
			UID string `json:"uid" validate:"required"`
		}
	)

	if err = common.BindAndValidateReq(c, &in); err != nil {
		return err
	}

	url := config.C.FapiAdapter.Urls.CheckBookingV4 + fmt.Sprintf("/%s", in.UID)

	rawResp, err := common.UapiGet(ctx, faClient, url)
	if err != nil {
		return
	}

	var (
		checkBookingResp CheckBookingAviaBusResponse
		customer         CheckBookingABCustomer
		route            CheckBookingAviaBusMDL
		busbytripID      BusByTripID
		bookingService   CBServicesAviaBus
		sellerPrice      SellerPrice
		cBAlternative    CBAlternative
		cBSegment        CBSegment
		avalDocs         []string
		customerIDs      []int32
	)

	err = json.Unmarshal(rawResp.Data, &checkBookingResp)
	if err != nil {
		ctx.Logger.WithError(err).WithField("response", string(rawResp.Data)).Error("cannot unmarshal UAPI response")
	}

	tripsmp := make(map[string]MaasTripAviaBus)
	servicesmp := make(map[string]CBServicesAviaBus)
	tripIDToServiceIDmp := make(map[string]string)
	cbaltermp := make(map[int]CBAlternative)
	customerIDToInfo := make(map[int32]CheckBookingABCustomer)

	for _, s := range checkBookingResp.Order.Services {
		priceByCustomerID := make(map[int]SellerPrice)
		for _, cp := range s.AlternativePriceInfos[0].CustomerPrices {
			sellerPrice.SellerPrice = cp.SellerPrice
			priceByCustomerID[cp.CustomerID] = sellerPrice
		}
		for _, tid := range s.TripIds {
			tripIDToServiceIDmp[tid] = s.ID
			for _, trp := range checkBookingResp.Order.Trips {
				busbytripIDmp := make(map[string]BusByTripID)
				if trp.ID == tid && !trp.IsReturnTrip {
					if s.ServiceType == "bus" && s.SellingType == "own" {
						bookingService.ObjectType = "serviceBusOwn"
						for _, contData := range s.ContainerData.TripsData {
							busbytripID.NeedPrintDocument = contData.GeneralData.NeedPrintDocument
							busbytripID.ID = tid
							busbytripID.FreeSeats = contData.SeatsData.FreeSeats
							busbytripID.Options = contData.GeneralData.Options
							trp.Options = contData.GeneralData.Options
							busbytripID.FreeSeatsCount = contData.SeatsData.FreeSeatsCount
							busbytripIDmp[tid] = busbytripID
						}
					} else if s.ServiceType == "flight" && s.SellingType == "own" {
						bookingService.ObjectType = "serviceFlightOwn"
						for i, aldata := range s.ContainerData.AlternativesData {
							faremp := make(map[int64]interface{})
							alterPriceByCustomerID := make(map[int]interface{})
							for _, cp := range s.AlternativePriceInfos[i].CustomerPrices {
								alterPriceByCustomerID[cp.CustomerID] = s.AlternativePriceInfos[i]
								cBAlternative.Fee = s.AlternativePriceInfos[i].Fee
								cBAlternative.SellerPrice = s.AlternativePriceInfos[i].SellerPrice
								cBAlternative.PriceStatus = s.AlternativePriceInfos[i].PriceStatus
								cBAlternative.Fee = s.AlternativePriceInfos[i].Fee
								cBAlternative.Reward = s.AlternativePriceInfos[i].Reward
								cBAlternative.Commission = s.AlternativePriceInfos[i].Commission
								cBAlternative.CurrencyCode = s.AlternativePriceInfos[i].CurrencyCode
								cBAlternative.Commission = s.AlternativePriceInfos[i].Commission
								cBAlternative.AgencyVAT = s.AlternativePriceInfos[i].AgencyVAT
								cBAlternative.ID = i
							}
							segmentbytripIDmp := make(map[string]CBSegment)
							for _, generalSegment := range s.ContainerData.GeneralData.Segments {
								cBSegment.FreeSeatsCount = generalSegment.FreeSeatsCount
								cBSegment.MarkCarrierName = generalSegment.MarkCarrierName
								cBSegment.MarkCarrierCode = generalSegment.MarkCarrierCode
								cBSegment.OpCarrierCode = generalSegment.OpCarrierCode
								cBSegment.ValCarrierCode = generalSegment.ValCarrierCode
								cBSegment.ValCarrierName = generalSegment.ValCarrierName
							}
							cBSegment.FareByCustomerID = faremp
							for iseg, segment := range aldata.Segments {
								cBSegment.Options = segment.Fares[0].Options
								cBSegment.FareFamily = segment.Fares[0].FareFamily
								for _, fare := range segment.Fares {
									cBSegment.FareFamilyName = fare.FareFamilyName
									cBSegment.ComfortType = fare.ComfortType
									faremp[fare.CustomerID] = fare
								}
								for itid, id := range s.TripIds {
									if itid == iseg {
										segmentbytripIDmp[id] = cBSegment
									}
								}
							}
							cBAlternative.PriceByCustomerID = alterPriceByCustomerID
							cBAlternative.SegmentByTripID = segmentbytripIDmp
							cbaltermp[i] = cBAlternative
						}
						bookingService.Alternatives = cbaltermp
						bookingService.BusByTripID = nil
						bookingService.SellerPrice = 0
						bookingService.Fee = 0
						bookingService.Price = 0
						bookingService.PriceByCustomer = nil
					} else {
						bookingService.ObjectType = ""
						bookingService.Alternatives = nil
						bookingService.ServiceType = s.ServiceType
					}
					// оставляем в availableDocuments только те документы, которые пересекаются у всех сервисов
					if s.SellingType == "own" && (s.ServiceType == "flight" || s.ServiceType == "train" || s.ServiceType == "bus") {
						if len(avalDocs) == 0 {
							avalDocs = append(avalDocs, s.AvailableDocumentTypes...)
						}
						commonElements := intersection(avalDocs, s.AvailableDocumentTypes)
						var avalNew []string
						for _, el := range commonElements {
							avalNew = append(avalNew, fmt.Sprintf("%v", el))
						}
						avalDocs = avalNew
					}
					tripsmp[tid] = trp
					bookingService.BusByTripID = busbytripIDmp
					bookingService.ID = s.ID
					if bookingService.ObjectType != "serviceFlightOwn" {
						bookingService.Fee = s.AlternativePriceInfos[0].Fee
						bookingService.Price = s.AlternativePriceInfos[0].Price
						bookingService.SellerPrice = s.AlternativePriceInfos[0].SellerPrice
						bookingService.PriceByCustomer = priceByCustomerID
					}
					bookingService.TripIds = s.TripIds
					bookingService.ProviderServiceCode = s.ProviderServiceCode
					bookingService.SellingType = s.SellingType
					bookingService.ServiceType = s.ServiceType
					route.Routes.Forward = append(route.Routes.Forward, tid)
					servicesmp[s.ID] = bookingService
				}
				if trp.ID == tid && trp.IsReturnTrip {
					if s.ServiceType == "bus" && s.SellingType == "own" {
						bookingService.ObjectType = "serviceBusOwn"
						for _, contData := range s.ContainerData.TripsData {
							busbytripID.NeedPrintDocument = contData.GeneralData.NeedPrintDocument
							busbytripID.ID = tid
							busbytripID.FreeSeats = contData.SeatsData.FreeSeats
							busbytripID.Options = contData.GeneralData.Options
							trp.Options = contData.GeneralData.Options
							busbytripID.FreeSeatsCount = contData.SeatsData.FreeSeatsCount
							busbytripIDmp[tid] = busbytripID
						}
					} else if s.ServiceType == "flight" && s.SellingType == "own" {
						bookingService.ObjectType = "serviceFlightOwn"
						for i, aldata := range s.ContainerData.AlternativesData {
							faremp := make(map[int64]interface{})
							alterPriceByCustomerID := make(map[int]interface{})
							for _, cp := range s.AlternativePriceInfos[i].CustomerPrices {
								alterPriceByCustomerID[cp.CustomerID] = s.AlternativePriceInfos[i]
								cBAlternative.Fee = s.AlternativePriceInfos[i].Fee
								cBAlternative.SellerPrice = s.AlternativePriceInfos[i].SellerPrice
								cBAlternative.PriceStatus = s.AlternativePriceInfos[i].PriceStatus
								cBAlternative.Fee = s.AlternativePriceInfos[i].Fee
								cBAlternative.Reward = s.AlternativePriceInfos[i].Reward
								cBAlternative.Commission = s.AlternativePriceInfos[i].Commission
								cBAlternative.CurrencyCode = s.AlternativePriceInfos[i].CurrencyCode
								cBAlternative.Commission = s.AlternativePriceInfos[i].Commission
								cBAlternative.AgencyVAT = s.AlternativePriceInfos[i].AgencyVAT
								cBAlternative.ID = i
							}
							segmentbytripIDmp := make(map[string]CBSegment)
							for _, generalSegment := range s.ContainerData.GeneralData.Segments {
								cBSegment.FreeSeatsCount = generalSegment.FreeSeatsCount
								cBSegment.MarkCarrierName = generalSegment.MarkCarrierName
								cBSegment.MarkCarrierCode = generalSegment.MarkCarrierCode
								cBSegment.OpCarrierCode = generalSegment.OpCarrierCode
								cBSegment.ValCarrierCode = generalSegment.ValCarrierCode
								cBSegment.ValCarrierName = generalSegment.ValCarrierName
							}
							cBSegment.FareByCustomerID = faremp
							for iseg, segment := range aldata.Segments {
								cBSegment.Options = segment.Fares[0].Options
								cBSegment.FareFamily = segment.Fares[0].FareFamily
								for _, fare := range segment.Fares {
									cBSegment.FareFamilyName = fare.FareFamilyName
									cBSegment.ComfortType = fare.ComfortType
									faremp[fare.CustomerID] = fare
								}
								for itid, id := range s.TripIds {
									if itid == iseg {
										segmentbytripIDmp[id] = cBSegment
									}
								}
							}
							cBAlternative.PriceByCustomerID = alterPriceByCustomerID
							cBAlternative.SegmentByTripID = segmentbytripIDmp
							cbaltermp[i] = cBAlternative
						}
						bookingService.Alternatives = cbaltermp
						bookingService.BusByTripID = nil
						bookingService.SellerPrice = 0
						bookingService.Fee = 0
						bookingService.Price = 0
						bookingService.PriceByCustomer = nil
					} else {
						bookingService.ObjectType = ""
						bookingService.Alternatives = nil
						bookingService.ServiceType = s.ServiceType
					}
					// оставляем в availableDocuments только те документы, которые пересекаются у всех сервисов
					if s.SellingType == "own" && (s.ServiceType == "flight" || s.ServiceType == "train" || s.ServiceType == "bus") {
						if len(avalDocs) == 0 {
							avalDocs = append(avalDocs, s.AvailableDocumentTypes...)
						}
						commonElements := intersection(avalDocs, s.AvailableDocumentTypes)
						var avalNew []string
						for _, el := range commonElements {
							avalNew = append(avalNew, fmt.Sprintf("%v", el))
						}
						avalDocs = avalNew
					}
					tripsmp[tid] = trp
					bookingService.BusByTripID = busbytripIDmp
					bookingService.ID = s.ID
					if bookingService.ObjectType != "serviceFlightOwn" {
						bookingService.Fee = s.AlternativePriceInfos[0].Fee
						bookingService.Price = s.AlternativePriceInfos[0].Price
						bookingService.SellerPrice = s.AlternativePriceInfos[0].SellerPrice
						bookingService.PriceByCustomer = priceByCustomerID
					}
					bookingService.TripIds = s.TripIds
					bookingService.ProviderServiceCode = s.ProviderServiceCode
					bookingService.SellingType = s.SellingType
					bookingService.ServiceType = s.ServiceType
					route.Routes.Backward = append(route.Routes.Backward, tid)
					servicesmp[s.ID] = bookingService
				}
			}
		}
	}

	for _, cust := range checkBookingResp.Order.Customers {
		customer.ID = cust.ID
		customer.Sex = cust.Sex
		customer.Age = cust.Age
		customer.SeatRequired = cust.SeatRequired
		customer.UserID = cust.UserID
		customerIDToInfo[cust.ID] = customer
		customerIDs = append(customerIDs, cust.ID)
	}

	if len(servicesmp) != 0 {
		route.TripIDToServiceID = tripIDToServiceIDmp
		route.AvailableDocumentTypes = avalDocs
		route.Trips = tripsmp
		route.Services = servicesmp
		route.Customers = customerIDToInfo
		route.CustomerIDs = customerIDs
		route.PriceInfo = checkBookingResp.Order.PriceInfo
		route.OrderStatus = checkBookingResp.Order.OrderStatus
		route.Owner = checkBookingResp.Order.Owner
		route.BeginDate = checkBookingResp.Order.BeginDate
		route.PriceInfo = checkBookingResp.Order.PriceInfo
		route.OrderStatus = checkBookingResp.Order.OrderStatus
		route.Owner = checkBookingResp.Order.Owner
		route.BeginDate = checkBookingResp.Order.BeginDate
		route.Places = checkBookingResp.Places
	}

	route.UID = in.UID
	route.Errors = rawResp.Error

	return c.JSON(http.StatusOK, echo.Map{"data": route})
}

func orderToMaasResp(ctx common.Context, rawResp *models.RawResponse) (maasOrderResp GetOrderMaasResp, err error) {
	var (
		order                    GetOrderResp
		customer                 OrderCustomer
		priceInfoToCustomerOrder PriceInfoToCustomerOrder
		booking                  OrderBooking
		orderSeatInfo            OrderSeatInfo
		sellerPrice              SellerPrice
		orderService             OrderServices
		customerIDs              []int32
	)

	err = json.Unmarshal(rawResp.Data, &order)
	if err != nil {
		err = fmt.Errorf("cannot unmarshal uapi response: %w", err)
		return
	}

	tripsmp := make(map[string]MaasTrip)
	servicesmp := make(map[string]interface{})
	tripIDToServiceIDmp := make(map[string]string)
	customerIDToInfo := make(map[int32]OrderCustomer)
	customerPricesMp := make(map[int32]PriceInfoToCustomerOrder)
	bookingsMp := make(map[int]OrderBooking)
	gBookingIDToFileStorage := make(map[int]string)

	// раскладывание сервисов и трипов
	for _, s := range order.Order.Services {
		priceByCustomerID := make(map[int]SellerPrice)
		for _, cp := range s.PriceInfo.CustomerPrices {
			sellerPrice.SellerPrice = cp.SellerPrice
			priceByCustomerID[cp.CustomerID] = sellerPrice
		}
		for _, tid := range s.TripIds {
			tripIDToServiceIDmp[tid] = s.ID
			for _, trp := range order.Order.Trips {
				optionsmp := make(map[string]string)
				if trp.ID == tid && !trp.IsReturnTrip {
					maasOrderResp.OrderRoutes.Forward = append(maasOrderResp.OrderRoutes.Forward, tid)
					for _, contData := range s.ContainerData.TripsData {
						for _, opt := range contData.GeneralData.Options {
							optionsmp[opt.Code] = opt.Name
						}
					}
					trp.Options = optionsmp
					tripsmp[tid] = trp
					orderService.Refund = s.PriceInfo.Refund
					orderService.ID = s.ID
					orderService.Fee = s.PriceInfo.Fee
					orderService.Price = s.PriceInfo.Price
					orderService.SellerPrice = s.PriceInfo.SellerPrice
					orderService.TripIds = s.TripIds
					orderService.ProviderServiceCode = s.ProviderServiceCode
					orderService.SellingType = s.SellingType
					orderService.ServiceType = s.ServiceType
					orderService.PriceByCustomer = priceByCustomerID
					servicesmp[s.ID] = orderService
				}
				if trp.ID == tid && trp.IsReturnTrip {
					maasOrderResp.OrderRoutes.Backward = append(maasOrderResp.OrderRoutes.Backward, tid)
					for _, contData := range s.ContainerData.TripsData {
						for _, opt := range contData.GeneralData.Options {
							optionsmp[opt.Code] = opt.Name
						}
					}
					trp.Options = optionsmp
					tripsmp[tid] = trp
					orderService.Refund = s.PriceInfo.Refund
					orderService.ID = s.ID
					orderService.Fee = s.PriceInfo.Fee
					orderService.Price = s.PriceInfo.Price
					orderService.SellerPrice = s.PriceInfo.SellerPrice
					orderService.TripIds = s.TripIds
					orderService.ProviderServiceCode = s.ProviderServiceCode
					orderService.SellingType = s.SellingType
					orderService.ServiceType = s.ServiceType
					orderService.PriceByCustomer = priceByCustomerID
					servicesmp[s.ID] = orderService
				}
			}
		}
	}

	for _, docForms := range order.Order.DocumentForms {
		gBookingIDToFileStorage[docForms.BookingID] = docForms.FileStorageID
	}

	for _, b := range order.Order.Bookings {
		priceInfoToCustomerOrder.Fee = b.PriceInfo.Fee
		booking.BookDocumentStatus = b.BookingStatus
		booking.ServiceID = b.ServiceID
		for bookID, fileStorage := range gBookingIDToFileStorage {
			if bookID == b.ID {
				booking.FileStorageID = fileStorage
			}
		}
		for _, bd := range b.BookDocuments {
			booking.TripIDs = bd.TripIds
			orderSeatInfo.SeatRequired = bd.SeatRequired
			orderSeatInfo.SeatNumber = bd.ContainerData.TicketData.SeatNumber
			booking.ID = bd.ID
			for _, bdc := range bd.BookDocumentCustomers {
				priceInfoToCustomerOrder.SellerPrice = bdc.CustomerPrice.SellerPrice
				orderSeatInfo.EquiveTariff = bdc.CustomerPrice.EquiveTariff
				orderSeatInfo.SellerPrice = bdc.CustomerPrice.SellerPrice
				customerPricesMp[bdc.CustomerID] = priceInfoToCustomerOrder
				customerIDToSeat := make(map[int32]OrderSeatInfo)
				customerIDToSeat[bdc.CustomerID] = orderSeatInfo
				booking.SeatByCustomerID = customerIDToSeat
				bookingsMp[bd.ID] = booking
			}
		}
	}
	// customers
	// order.Order.PriceInfo
	for _, cust := range order.Order.Customers {
		customer.Birthdate = cust.Birthdate
		customer.ID = cust.ID
		customer.Email = cust.Email
		customer.Sex = cust.Sex
		customer.Phone = cust.Phone
		for _, custDocs := range cust.CustomerDocuments {
			customer.Middlename = custDocs.Middlename
			customer.FirstName = custDocs.Firstname
			customer.LastName = custDocs.Lastname
			customer.Number = custDocs.Number
			customer.Type = custDocs.Type
			customer.Citizenship, _ = dictionary.GetCitizenshipByAlpha(ctx, custDocs.Citizenship)
			customer.ExpireDate = custDocs.ExpireDate
			for k, v := range customerPricesMp {
				if k == customer.ID {
					customer.SellerPrice = v.SellerPrice
					customer.Fee = v.Fee
				}
			}
		}
		customerIDToInfo[cust.ID] = customer
		customerIDs = append(customerIDs, cust.ID)
	}

	maasOrderResp.TripIDToServiceID = tripIDToServiceIDmp

	if ctx.UserID == order.Order.UserID {
		maasOrderResp.IsOwnedByUser = true
	}

	maasOrderResp.Bookings = bookingsMp
	maasOrderResp.Customers = customerIDToInfo
	maasOrderResp.CustomerIDs = customerIDs
	maasOrderResp.DateTimeToPay = order.Order.DateTimeToPay
	maasOrderResp.Trips = tripsmp
	maasOrderResp.Services = servicesmp
	maasOrderResp.OrderID = order.Order.ID
	maasOrderResp.SearchParams = order.SearchParams
	maasOrderResp.Places = order.Places
	maasOrderResp.OrderStatus = order.Order.OrderStatus
	maasOrderResp.Error = rawResp.Error

	return
}

func intersection(a interface{}, b interface{}) []interface{} {
	set := make([]interface{}, 0)
	av := reflect.ValueOf(a)

	for i := 0; i < av.Len(); i++ {
		el := av.Index(i).Interface()
		if containsEl(b, el) {
			set = append(set, el)
		}
	}

	return set
}

func containsEl(a interface{}, e interface{}) bool {
	v := reflect.ValueOf(a)

	for i := 0; i < v.Len(); i++ {
		if v.Index(i).Interface() == e {
			return true
		}
	}
	return false
}
