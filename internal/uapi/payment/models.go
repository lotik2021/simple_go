package payment

type PayRequest struct {
	OrderId     int    `json:"order_id" validate:"required"`
	Uid         string `json:"uid" validate:"required"`
	PaymentType string `json:"payment_type"`
	Email       string `json:"email" validate:"required"`
}

type OrderStateRequest struct {
	OrderId int    `json:"order_id" validate:"required"`
	Uid     string `json:"uid" validate:"required"`
}

type UapiPayRequest struct {
	OrderId     int    `json:"orderId"`
	RequestId   string `json:"requestId"`
	BackUrl     string `json:"backUrl"`
	PaymentType string `json:"paymentType"`
}

type AdditionalPayRequest struct {
	EventId     int64  `json:"eventId" validate:"required"`
	Uid         string `json:"uid"`
	PaymentType string `json:"paymentType"`
}

type UAdditionalPayRequest struct {
	EventId     int64  `json:"eventId"`
	BackUrl     string `json:"backUrl"`
	PaymentType string `json:"paymentType"`
}

type UPayData struct {
	RedirectUrl string `json:"redirectUrl"`
}

type RefundOrderRequest struct {
	OrderId         uint32 `json:"orderId,omitempty"`
	BookDocumentIds []int  `json:"bookDocumentIds"`
}
