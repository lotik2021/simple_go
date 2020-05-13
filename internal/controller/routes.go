package controller

import (
	"net/http"

	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"bitbucket.movista.ru/maas/maasapi/internal/middlewares"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"bitbucket.movista.ru/maas/maasapi/internal/static"
	"bitbucket.movista.ru/maas/maasapi/internal/uapi/booking"
	"bitbucket.movista.ru/maas/maasapi/internal/uapi/fapi"
	"bitbucket.movista.ru/maas/maasapi/internal/uapi/payment"
	"bitbucket.movista.ru/maas/maasapi/internal/uapi/search"
	"bitbucket.movista.ru/maas/maasapi/internal/uapi/web_tmp"
	"github.com/labstack/echo/v4"
)

func Add(api *echo.Group) {
	// AUTH
	api.POST("/auth/create", authCreate)
	api.POST("/auth/refresh", authRefreshToken)

	api.POST("/auth/sendSms", ClosedRoute)
	api.POST("/auth/regenerateCode", ClosedRoute)

	api.POST("/auth/sendSmsMobile", authSendSmsSecure, middlewares.Token)
	api.POST("/auth/regenerateCodeMobile", authRegenerateCodeSecure, middlewares.Token)

	api.POST("/auth/sendSmsWeb", authSendSmsSecure, middlewares.Token)
	api.POST("/auth/regenerateCodeWeb", authRegenerateCodeSecure, middlewares.Token)

	api.POST("/auth/authorize", authAuthorize, middlewares.Token)
	api.POST("/auth/completeRegistration", authCompleteRegistration, middlewares.Token)
	api.GET("/auth/proxy_user_info", authProxyUserInfo, middlewares.Token)
	api.POST("/auth/proxy_user_info", authProxyUserInfo, middlewares.Token)
	api.POST("/auth/logout", authLogout, middlewares.Token)

	api.POST("/user/favorites/create", createFavorite, middlewares.Token)
	api.POST("/user/favorites/delete/:id", deleteFavorite, middlewares.Token)
	api.PUT("/user/favorites/:id", updateFavorite, middlewares.Token)
	api.GET("/user/favorites", listFavorite, middlewares.Token)
	api.POST("/user/favorites", listFavorite, middlewares.Token)
	api.POST("/user/favorites/:id/toggleinactions", toggleFavoriteInActions, middlewares.Token)

	api.POST("/user/ping/location", pingLocation, middlewares.Token)

	api.POST("/user/pay/troika", payTroika, middlewares.Token)
	api.GET("/user/pay/troika", getTroikaParameters, middlewares.Token)
	api.GET("/user/pay/strelka/cardparameters", getStrelkaCardParameters, middlewares.Token)
	api.POST("/user/pay/strelka/paymentparameters", getStrelkaPaymentParameters, middlewares.Token)
	api.POST("/user/pay/strelka", payStrelka, middlewares.Token)
	api.GET("/user/pay/travelcards", getTravelCards, middlewares.Token)
	api.POST("/user/pay/travelcards", getTravelCards, middlewares.Token)
	api.POST("/user/pay/save", savePayment, middlewares.Token)
	api.POST("/user/payments", getPayments, middlewares.Token)

	api.POST("/user/deeplink/click", setUsedLinkInTaxiOrder, middlewares.Token)

	api.POST("/user/getprofile", getDeviceProfile, middlewares.Token)
	api.POST("/user/updateusergoogletransports", updateDeviceGoogleTransportSettings, middlewares.Token)
	api.POST("/user/transport/userfilters", getDeviceTransportSettings, middlewares.Token)

	api.POST("/user/setosplayerid", setOSPlayerID, middlewares.Token)

	if !config.IsProd() {
		api.POST("/admin/push/send", sendPushNotification)
		//api.DELETE("/user/:id", i.userRemoveUser)
	}

	api.GET("/places/autosuggestions", webPlacesAutoSuggestions)

	api.POST("/dialogs/dialog", dialog, middlewares.Token)

	api.POST("/booking/v2/GetFareRulesV2", booking.GetFareRulesV2, middlewares.Token)
	api.POST("/booking/v4/CreateBooking", fapi.CreateBookingV4, middlewares.Token)
	api.POST("/booking/v4/CheckBooking", checkBookingV4, middlewares.Token)
	api.POST("/booking/v2/GetOrder", fapi.GetOrderV2, middlewares.Token)
	api.POST("/booking/v4/CreateOrder", fapi.CreateOrderV4, middlewares.Token)
	api.POST("/booking/v1/AddServicesToOrder", fapi.AddServicesToOrder, middlewares.Token)
	api.POST("/booking/v4/CheckBookingMobile", fapi.CheckBookingMobileV4, middlewares.Token)
	api.POST("/booking/v1/cancelorder", booking.CancelOrder, middlewares.Token)
	api.POST("/booking/v2/GetOrderByEmailMobile", fapi.GetOrderByEmailMobileV2, middlewares.Token)
	api.POST("/booking/v2/GetOrderByEmail", fapi.GetOrderByEmailV2, middlewares.Token)
	api.POST("/booking/v2/GetOrderMobile", fapi.GetOrderMobileV2, middlewares.Token)
	api.POST("/booking/v2/CheckBookingAvB", fapi.CheckBookingAviaBus, middlewares.Token)

	api.POST("/blank/v1/base64", fapi.GetBlankBase64, middlewares.Token)
	api.POST("/blank/v1/openPdfFile", fapi.GetOpenedBlankFile, middlewares.Token)
	api.POST("/blank/v1/pdf", fapi.GetBlankPDF, middlewares.Token)
	api.POST("/blank/v1/zip", fapi.GetBlankZIP, middlewares.Token)
	api.POST("/blank/v1/pdfByIds", fapi.GetBlankPDFs, middlewares.Token)

	api.GET("/blank/v1/:fileStorageId/base64", fapi.WebGetBlankBase64)
	api.GET("/blank/v1/:fileStorageId/openPdfFile", fapi.WebGetOpenedBlankFile)
	api.GET("/blank/v1/:fileStorageId/pdf", fapi.WebGetBlankPDF)
	api.GET("/blank/v1/zip", fapi.WebGetBlankZIP)
	api.GET("/blank/v1/pdfByIds", fapi.WebGetBlankAllPDF)

	api.POST("/order/v1/checkRefund", fapi.CheckRefund, middlewares.Token)

	api.GET("/skills", getSills)
	api.POST("/skills", getSills)

	// SEARCH
	api.POST("/search/v2/routes", searchFindRoutesByLocationV2, middlewares.Token)
	api.POST("/search/v2/routesByLocation", searchFindRoutesByLocationV2, middlewares.Token)

	api.POST("/search/autocomplete", searchFindAutocompletePlaces, middlewares.Token)
	api.POST("/search/routes", searchFindRoutesByLocation, middlewares.Token)
	api.POST("/search/routesByLocation", searchFindRoutesByLocation, middlewares.Token)

	api.POST("/search/places/findByName", searchFindPlaceByName, middlewares.Token)
	api.POST("/search/places/findByIP", searchFindPlaceByIP, middlewares.Token)

	api.POST("/search/places/:placeID", searchFindPlaceByID, middlewares.Token)
	api.POST("/search/places/history", searchGetPlaceHistory, middlewares.Token)
	api.POST("/search/places/favorite", searchGetPlaceFavorite, middlewares.Token)
	api.POST("/search/places/save", searchSavePlace, middlewares.Token)

	api.POST("/search/searchAsync", searchSearchAsync, middlewares.Token)
	api.POST("/search/getSearchStatus", searchGetSearchStatus, middlewares.Token)
	api.POST("/search/getSearchResults", searchGetSearchResults, middlewares.Token)
	api.POST("/search/getPathGroup", searchGetPathGroup, middlewares.Token)
	api.POST("/search/getSegmentRoutes", searchGetSegmentRoutes, middlewares.Token)
	api.POST("/search/saveselectedroutes", fapi.SaveSelectedRoutesV5, middlewares.Token)
	api.POST("/search/getselectedroutes", fapi.GetSelectedRoutesV5, middlewares.Token)
	api.POST("/search/v5/getselectedroutes", fapi.GetSelectedRoutesV5, middlewares.Token)
	api.POST("/search/refundOrderAsync", searchRefundOrderAsync, middlewares.Token)

	api.POST("/search/schedule/train", searchGetTrainSchedule, middlewares.Token)

	api.POST("/search/asyncAvailable", searchAsyncStatus)
	api.POST("/seatAutoSelection", seatSelection)

	api.POST("/weather/getWeather", getWeather, middlewares.Token)
	api.POST("/weather/getForecast", getForecast, middlewares.Token)

	api.POST("/purchase/hints", getPurchaseHints)

	// методы записной книжки
	api.POST("/notebook/v2/documents", getDocumentsV2, middlewares.Token)
	api.POST("/notebook/v2/deactivate", documentsDeactivateV2, middlewares.Token)
	api.POST("/notebook/v2/edit", documentsEditV2, middlewares.Token)
	api.POST("/notebook/v2/create", documentsCreateV2, middlewares.Token)

	api.POST("/payment/v1/pay", paymentPay, middlewares.Token)
	api.POST("/payment/v1/orderState", paymentOrderState, middlewares.Token)

	api.POST("/references/v1/citizenships", Citizenships)

	// новые методы
	api.GET("/confirm/v1/email", confirmEmail, middlewares.Token)
	api.GET("/content/v1/cms", queryCms, middlewares.Token)
	api.POST("/ga/v1/sendGa", sendGA)
	api.POST("/notification/v1/feedback", sendFeedback, middlewares.Token)
	api.POST("/places/v1/placesByIds", placesByIds, middlewares.Token)
	api.POST("/places/v1/searchByName", searchFindPlaceByName)
	api.POST("/payment/v1/additionalPay", additionalPay, middlewares.Token)

	api.GET("/booking/v2/CheckBooking/:uid", checkBookingV2)
	api.POST("/booking/v2/CreateBooking", fapi.CreateBookingV2)
	api.GET("/search/v4/getselectedroutes/:uid", fapi.GetSelectedRoutesV4)
	api.GET("/payment/v1/getorderstate", payment.GetOrderStateV1, middlewares.Token)
	api.POST("/booking/v1/changeERegistration", booking.ChangeERegistrationV1, middlewares.Token)
	api.POST("/booking/v1/checkRefund", booking.CheckRefundV1, middlewares.Token)
	api.GET("/v1/events/:eventId", web_tmp.GetEventFromOrderV2, middlewares.Token)
	api.POST("/booking/v1/GetFareRules", booking.GetFareRulesV1, middlewares.Token)
	api.POST("/payment/v1/refundOrder", payment.RefundOrderV1, middlewares.Token)
	api.POST("/search/v3/routeDetails", search.RouteDetailsV3, middlewares.Token)
	api.GET("/search/v4/saveselectedroutes/:uid/:forwardRouteId", fapi.SaveSelectedRoutesV4)

	api.GET("/static/v1/content", static.LoadContent)

	return
}

func seatSelection(c echo.Context) (err error) {
	return c.JSON(http.StatusOK, models.UapiResponse{
		Data: echo.Map{"auto_select": config.C.SeatAutoSelection},
	})
}

func ClosedRoute(c echo.Context) (err error) {
	err = models.Error{
		StatusCode: http.StatusInternalServerError,
		Code:       "CLOSED_ROUTE",
		Message:    "CLOSED_ROUTE",
		Type:       models.SystemError,
	}

	return err
}
