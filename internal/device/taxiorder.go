package device

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/logger"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
)

func SaveTaxiOrder(ctx common.Context, deeplinks map[string]models.DeepLink) {
	if len(deeplinks) == 0 {
		return
	}

	orders := make([]TaxiOrder, 0)
	for i, v := range deeplinks {
		orders = append(orders, TaxiOrder{Id: v.Id, DeviceID: ctx.DeviceID, Provider: i})
	}

	if _, err := ctx.DB.Model(&orders).Insert(); err != nil {
		logger.Log.Errorf("cannot save taxi orders %+v, err - %w", orders, err)
	}
}

func SetUsedLinkInTaxiOrder(ctx common.Context, orderId, usedLink string) (err error) {
	_, err = ctx.DB.Model(&TaxiOrder{Id: orderId, DeviceID: ctx.DeviceID, UsedLink: usedLink}).WherePK().UpdateNotZero()

	return
}
