package analytics

import (
	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"encoding/json"
	"time"
)

type CarrotEventRequest struct {
	ID       string        `json:"id"`
	Event    string        `json:"event" validate:"required"`
	Params   *CarrotParams `json:"params,omitempty"`
	Client   string        `json:"client"`
	Created  int64         `json:"created,omitempty"`
	ByUserID bool          `json:"by_user_id,omitempty"`
}

type CarrotPropsRequest struct {
	ID         string              `json:"id"`
	Operations []map[string]string `json:"operations,omitempty"`
	ByUserID   bool                `json:"by_user_id"`
}

type CarrotParams struct {
	OrderID   string `json:"orderID"`
	BirthDate string `json:"birthDate"`
	Sex       string `json:"sex"`
	FirstName string `json:"firstName"`
	Departure string `json:"departure"`
	Arrival   string `json:"arrival"`
	From      string `json:"from"`
	To        string `json:"to"`
	Price     string `json:"price"`
	Status    string `json:"status"`
	Flight    string `json:"flight,omitempty"`
	Bus       string `json:"bus,omitempty"`
	Train     string `json:"train,omitempty"`
	Taxi      string `json:"taxi,omitempty"`
	CityFrom  string `json:"cityFrom,omitempty"`
	CityTo    string `json:"cityTo,omitempty"`
}

type CarrotDBEvent struct {
	tableName           struct{}        `pg:"maasapi.events"`
	ID                  int             `pg:"id,pk"`
	OrderID             int             `pg:"order_id"`
	PaymentStatus       bool            `pg:"payment_status"`
	Event               json.RawMessage `pg:"event"`
	EventDeliveryStatus bool            `pg:"event_delivery_status"`
	CreatedAt           models.Time     `pg:"created_at"`
	UpdatedAt           models.Time     `pg:"updated_at"`
}

type OrderInfo struct {
	Data struct {
		Order struct {
			BeginDate   time.Time  `json:"beginDate"`
			CustomerIds []int      `json:"customerIds"`
			Customers   []Customer `json:"customers"`
			ID          int        `json:"id"`
			OrderStatus string     `json:"orderStatus"`
			Owner       string     `json:"owner"`
			PriceInfo   struct {
				AgencyVAT      float64 `json:"agencyVAT"`
				Commission     float64 `json:"commission"`
				CurrencyCode   string  `json:"currencyCode"`
				CustomerPrices []struct {
					CurrencyCode string  `json:"currencyCode"`
					CustomerID   float64 `json:"customerId"`
					EquiveTariff float64 `json:"equiveTariff"`
					SellerPrice  float64 `json:"sellerPrice"`
					Tariff       float64 `json:"tariff"`
					Taxes        []interface {
					} `json:"taxes"`
					Vat float64 `json:"vat"`
				} `json:"customerPrices"`
				Fee         float64 `json:"fee"`
				Price       float64 `json:"price"`
				PriceStatus string  `json:"priceStatus"`
				Reward      float64 `json:"reward"`
				SellerPrice float64 `json:"sellerPrice"`
			} `json:"priceInfo"`
			Services []struct {
				AlternativePriceInfos []struct {
					AgencyVAT      float64 `json:"agencyVAT"`
					Commission     int     `json:"commission"`
					CurrencyCode   string  `json:"currencyCode"`
					CustomerPrices []struct {
						CurrencyCode string  `json:"currencyCode"`
						CustomerID   int     `json:"customerId"`
						EquiveTariff float64 `json:"equiveTariff"`
						SellerPrice  float64 `json:"sellerPrice"`
						Tariff       float64 `json:"tariff"`
						Vat          float64 `json:"vat"`
					} `json:"customerPrices"`
					Fee         float64 `json:"fee"`
					Price       float64 `json:"price"`
					PriceStatus string  `json:"priceStatus"`
					Reward      float64 `json:"reward"`
					SellerPrice float64 `json:"sellerPrice"`
				} `json:"alternativePriceInfos"`
				AvailableDocumentTypes []string `json:"availableDocumentTypes"`
				ContainerData          interface {
				} `json:"containerData"`
				Descr       string `json:"descr"`
				ID          string `json:"id"`
				PrivateData struct {
					RouteID string `json:"routeId"`
				} `json:"privateData"`
				ProviderServiceCode string   `json:"providerServiceCode"`
				SellingType         string   `json:"sellingType"`
				ServiceKind         string   `json:"serviceKind"`
				ServiceStatus       string   `json:"serviceStatus"`
				ServiceType         string   `json:"serviceType"`
				TripIds             []string `json:"tripIds"`
			} `json:"services"`
			Trips []struct {
				Arrival            time.Time `json:"arrival"`
				BoardingDuration   int       `json:"boardingDuration"`
				CarrierName        string    `json:"carrierName"`
				Departure          time.Time `json:"departure"`
				Descr              string    `json:"descr"`
				Direction          string    `json:"direction"`
				Distance           int       `json:"distance"`
				Duration           int       `json:"duration"`
				FromDescr          string    `json:"fromDescr"`
				FromID             int       `json:"fromId"`
				ID                 string    `json:"id"`
				IsReturnTrip       bool      `json:"isReturnTrip"`
				IsTransfer         bool      `json:"isTransfer"`
				ToDescr            string    `json:"toDescr"`
				ToID               int       `json:"toId"`
				TransportClass     string    `json:"transportClass"`
				TripType           string    `json:"tripType"`
				UnboardingDuration int       `json:"unboardingDuration"`
			} `json:"trips"`
			UserID int `json:"userId"`
		} `json:"order"`
		SearchParams struct {
			DepartureBegin string `json:"departureBegin"`
			From           int    `json:"from"`
			FromPlace      struct {
				ID int `json:"id"`
			} `json:"fromPlace"`
			ReturnDepartureBegin string `json:"returnDepartureBegin"`
			SearchID             int    `json:"searchId"`
			To                   int    `json:"to"`
			ToPlace              struct {
				ID int `json:"id"`
			} `json:"toPlace"`
			TripTypes []string `json:"tripTypes"`
		} `json:"searchParams"`
	} `json:"data"`
	Errors []struct {
		Code    string `json:"code"`
		IsError bool   `json:"isError"`
		Message string `json:"message"`
		Type    string `json:"type"`
		Values  struct {
			ServiceID string `json:"serviceId"`
			Source    string `json:"source"`
		} `json:"values"`
	} `json:"errors"`
}

type Customer struct {
	ID                int32              `json:"id, omitempty"`
	Birthdate         string             `json:"birthDate, omitempty"`
	Sex               string             `json:"sex, omitempty"`
	CustomerDocuments []CustomerDocument `json:"customerDocuments, omitempty"`
	Phone             string             `json:"phone, omitempty "`
	Email             string             `json:"email, omitempty"`
}

type CustomerDocument struct {
	ID               int64  `json:"id"`
	Type             string `json:"type"`
	Number           string `json:"number"`
	Firstname        string `json:"firstName"`
	Middlename       string `json:"middleName"`
	Lastname         string `json:"lastName"`
	Citizenship      string `json:"citizenship"`
	ExpireDate       string `json:"expireDate, omitempty"`
	IssueDate        string `json:"issueDate, omitempty"`
	IssuingAuthority string `json:"issuingAuthority, omitempty"`
	IsActive         bool   `json:"isActive"`
}
