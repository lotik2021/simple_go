package device

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/transport/travelcard"
	"github.com/spf13/cast"
)

const (
	PAYMENT_TYPE_TROIKA  = "troika"
	PAYMENT_TYPE_STRELKA = "strelka"
)

func PayTroika(ctx common.Context, amount int, cardNumber string) (redirectURL string, err error) {
	paymentId, err := CreatePayment(ctx, PAYMENT_TYPE_TROIKA, cardNumber, amount)
	if err != nil {
		return
	}

	return travelcard.PayTroika(ctx, amount, cardNumber, paymentId)
}

func PayStrelka(ctx common.Context, amount int, cardNumber, cardTypeID string) (redirectURL string, err error) {
	paymentId, err := CreatePayment(ctx, PAYMENT_TYPE_STRELKA, cardNumber, amount)
	if err != nil {
		return
	}

	return travelcard.PayStrelka(ctx, amount, paymentId, cardNumber, cardTypeID)
}

func CreatePayment(ctx common.Context, paymentType, cardNum string, amount int) (paymentId string, err error) {
	p := &TravelCardPayment{
		DeviceID:   ctx.DeviceID,
		Amount:     cast.ToFloat64(amount),
		CardType:   paymentType,
		CardNumber: cardNum,
		Currency:   "rub",
	}

	_, err = ctx.DB.Model(p).Insert()
	if err != nil {
		return
	}

	paymentId = p.ID

	return
}

func UpdatePaymentProcessedStatus(ctx common.Context, paymentId string, success bool) (err error) {
	payment := &TravelCardPayment{
		ID:        paymentId,
		Processed: success,
	}

	_, err = ctx.DB.Model(payment).WherePK().UpdateNotZero()
	if err != nil {
		return
	}

	return
}

func GetPayments(ctx common.Context) (items []TravelCardPayment, err error) {
	err = ctx.DB.Model(items).Where("device_id = ?", ctx.DeviceID).Order("created_at asc").Select()

	return
}
