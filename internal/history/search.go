package history

import (
	"time"

	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/logger"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"bitbucket.movista.ru/maas/maasapi/internal/uapi/fapi"
	"github.com/google/uuid"
)

func SaveMovistaWebSearch(ctx common.Context, searchParams fapi.SearchParamsV5, searchUID string) {

	const LongTime = "2006-01-02T15:04:05"
	const ShortTime = "2006-01-02"

	departureTime, err := time.Parse(LongTime, searchParams.DepartureBegin)
	if err != nil {
		err = nil
		departureTime, _ = time.Parse(ShortTime, searchParams.DepartureBegin)
	}

	var arrivalTime time.Time
	if searchParams.ArrivalDate != "" {
		arrivalTime, err = time.Parse(LongTime, searchParams.ArrivalDate)
		if err != nil {
			err = nil
			arrivalTime, _ = time.Parse(ShortTime, searchParams.ArrivalDate)
		}
	}

	customers := make([]Customer, len(searchParams.Customers), len(searchParams.Customers))
	for i, customer := range searchParams.Customers {
		customers[i] = Customer{
			Age:          customer.Age,
			SeatRequired: customer.SeatRequired,
		}
	}

	ms := &MovistaSearch{
		ID:            searchUID,
		DeviceID:      ctx.DeviceID,
		ToID:          searchParams.To,
		FromID:        searchParams.From,
		TripTypes:     searchParams.TripTypes,
		CultureCode:   searchParams.CultureCode,
		CurrencyCode:  searchParams.CurrencyCode,
		ComfortType:   searchParams.AirServiceClass,
		ArrivalTime:   models.Time{Time: arrivalTime},
		DepartureTime: models.Time{Time: departureTime},
		Customers:     customers,
	}

	SaveMovistaSearch(ctx, ms)
}

func SaveMovistaSearch(ctx common.Context, ms *MovistaSearch) (err error) {
	if ms.ID == "" {
		// searchV4 не передает uid в параметрах запроса к fapiadapter
		ms.ID = uuid.New().String()
	}

	if ms.Customers == nil || len(ms.Customers) == 0 {
		ms.Customers = []Customer{DefaultCustomer}
	}

	for _, customer := range ms.Customers {
		customer.MovistaSearchId = ms.ID
		_, err = ctx.DB.Model(&customer).Insert()
		if err != nil {
			logger.Log.Warnf("cannot save customer from the search: %v, %+v", err, customer)
		}
	}

	_, err = ctx.DB.Model(ms).Insert()
	if err != nil {
		logger.Log.Warnf("cannot save search history: %v, %+v", err, ms)
	}

	return err
}

func SaveGoogleSearch(ctx common.Context, gs *GoogleSearch) (err error) {
	_, err = ctx.DB.Model(gs).Insert()
	return
}
