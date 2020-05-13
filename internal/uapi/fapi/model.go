package fapi

import (
	"bitbucket.movista.ru/maas/maasapi/internal/dictionary"
	"time"

	"bitbucket.movista.ru/maas/maasapi/internal/models"
	"bitbucket.movista.ru/maas/maasapi/internal/uapi/place"
)

// searchAsync
type AsyncRequest struct {
	UID          string       `json:"uid"`
	SearchParams SearchParams `json:"searchParams"`
}

type SearchParams struct {
	SearchID       int         `json:"searchId"`
	CurrencyCode   string      `json:"currencyCode"`
	CultureCode    string      `json:"cultureCode"`
	DepartureBegin string      `json:"departureBegin"`
	FromPlace      Place       `json:"fromPlace"`
	ToPlace        Place       `json:"toPlace"`
	Customers      []Customers `json:"customers"`
	TripTypes      []string    `json:"tripTypes"`
}

type AsyncRequestV5 struct {
	SearchParams SearchParamsV5 `json:"searchParams"`
}

type SearchParamsV5 struct {
	CurrencyCode    string       `json:"currencyCode"`
	DepartureBegin  string       `json:"departureBegin"`
	ArrivalDate     string       `json:"returnDepartureBegin"`
	CultureCode     string       `json:"cultureCode"`
	From            int          `json:"from"`
	To              int          `json:"to"`
	Customers       []Customers  `json:"customers"`
	TripTypes       []string     `json:"tripTypes"`
	AirServiceClass string       `json:"airServiceClass"`
	SearchID        int          `json:"searchId"`
	FromPlace       *place.Place `json:"fromPlace,omitempty"`
	ToPlace         *place.Place `json:"toPlace,omitempty"`
}

type Place struct {
	ID        int     `json:"id,omitempty"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Accuracy  string  `json:"accuracy,omitempty"`
}

type Customers struct {
	ID           int  `json:"id"`
	Age          int  `json:"age"`
	SeatRequired bool `json:"seat_required"`
}

type AsyncResponseV5 struct {
	UID            string `json:"uid"`
	DateTimeUpdate string `json:"dateTimeUpdate"`
	IsComplete     bool   `json:"isComplete"`
}

type SearchResult struct {
	PathGroups   []PathGroup    `json:"pathGroups"`
	SearchParams SearchParamsV5 `json:"searchParams"`
}

type PathGroup struct {
	ID               string         `json:"id"`
	Title            string         `json:"title"`
	Descr            string         `json:"descr"`
	MinDuration      string         `json:"minDuration"`
	MinDurationTitle string         `json:"minDurationTitle"`
	MinPrice         float64        `json:"minPrice"`
	MinPriceTitle    string         `json:"minPriceTitle"`
	TripTypes        []string       `json:"tripTypes"`
	RoutesCount      int            `json:"routesCount"`
	Status           string         `json:"bookableStatus"`
	State            string         `json:"state"`
	PlaceIDs         []int          `json:"placeIds"`
	Places           interface{}    `json:"places"`
	SearchParams     SearchParamsV5 `json:"searchParams"`
	UID              string         `json:"uid"`
	SearchUID        string         `json:"search_uid"`
	TripTypeSequence [][]string     `json:"tripTypeSequence"`
	CreatedAt        *models.Time   `json:"created_at"`
}

type GetPathGroupRequest struct {
	UID         string `json:"uid"`
	PathGroupID string `json:"pathGroupId"`
}

type PathGroupResponse struct {
	Trips map[string]PathGroupTrip
}

type PathGroupTrip struct {
	ObjectType        string         `json:"objectType"`
	TripTrain         *TrainSchedule `json:"tripTrain,omitempty"`
	TripTrainSuburban *TrainSchedule `json:"tripTrainSuburban,omitempty"`
}

type TrainSchedule struct {
	Arrival            string `json:"arrival"`
	BoardingDuration   int    `json:"boardingDuration"`
	CarTypeName        string `json:"carTypeName"`
	CarrierName        string `json:"carrierName"`
	Departure          string `json:"departure"`
	Descr              string `json:"descr"`
	Direction          string `json:"direction"`
	Distance           int    `json:"distance"`
	Duration           int    `json:"duration"`
	FromDescr          string `json:"fromDescr"`
	FromID             int    `json:"fromId"`
	ID                 string `json:"id"`
	IsReturnTrip       bool   `json:"isReturnTrip"`
	Number             string `json:"number"`
	Status             string `json:"status"`
	ToDescr            string `json:"toDescr"`
	ToID               int    `json:"toId"`
	TripType           string `json:"tripType"`
	UnboardingDuration int    `json:"unboardingDuration"`
}

type CreateBookingResponse struct {
	Order struct {
		ID             int    `json:"id"`
		DateTimeToPay  string `json:"dateTimeToPay"`
		DateTimeCreate string `json:"dateTimeCreate"`
		UserID         int    `json:"userId"`
	}
}

type CheckBookingResponse struct {
	Order struct {
		BeginDate   time.Time   `json:"beginDate"`
		CustomerIds []int       `json:"customerIds"`
		Customers   interface{} `json:"customers"`
		ID          int         `json:"id"`
		OrderStatus string      `json:"orderStatus"`
		Owner       string      `json:"owner"`
		PriceInfo   struct {
			AgencyVAT      float64 `json:"agencyVAT"`
			Commission     float64 `json:"commission"`
			CurrencyCode   string  `json:"currencyCode"`
			CustomerPrices []struct {
				CurrencyCode string        `json:"currencyCode"`
				CustomerID   float64       `json:"customerId"`
				EquiveTariff float64       `json:"equiveTariff"`
				SellerPrice  float64       `json:"sellerPrice"`
				Tariff       float64       `json:"tariff"`
				Taxes        []interface{} `json:"taxes"`
				Vat          float64       `json:"vat"`
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
				Commission     float64 `json:"commission"`
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
			ContainerData          struct {
				TripsData []struct {
					GeneralData struct {
						NeedPrintDocument bool `json:"needPrintDocument"`
						Options           []struct {
							Code string `json:"code"`
							Name string `json:"name"`
						} `json:"options"`
					} `json:"generalData"`
					SeatsData struct {
						FreeSeats      interface{} `json:"freeSeats"`
						FreeSeatsCount int64       `json:"freeSeatsCount"`
					} `json:"seatsData"`
				} `json:"tripsData"`
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
		Trips  []MaasTrip `json:"trips"`
		UserID int        `json:"userId"`
	} `json:"order"`
	Places       map[string]BookingPlace `json:"places"`
	SearchParams struct {
		AirServiceClass string `json:"airServiceClass"`
		CultureCode     string `json:"cultureCode"`
		CurrencyCode    string `json:"currencyCode"`
		Customers       []struct {
			Age          int  `json:"age"`
			ID           int  `json:"id"`
			SeatRequired bool `json:"seatRequired"`
		} `json:"customers"`
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
}

type BookingPlace struct {
	CityName     string  `json:"cityName"`
	CountryName  string  `json:"countryName"`
	Description  string  `json:"description"`
	FullName     string  `json:"fullName"`
	ID           int     `json:"id"`
	Lat          float64 `json:"lat"`
	Lon          float64 `json:"lon"`
	Name         string  `json:"name"`
	PlaceclassID int     `json:"placeclassId"`
	StateName    string  `json:"stateName"`
	TimeZone     string  `json:"timeZone"`
	TypePlace    int     `json:"typePlace"`
}

type CheckBookingMobile struct {
	BeginDate         time.Time              `json:"beginDate,omitempty"`
	Forward           []TripIDToServiceID    `json:"forward,omitempty"`
	Backward          []TripIDToServiceID    `json:"backward,omitempty"`
	Trips             map[string]MaasTrip    `json:"trips,omitempty"`
	Services          map[string]interface{} `json:"services,omitempty"`
	Customers         interface{}            `json:"customers,omitempty"`
	PriceInfo         interface{}            `json:"priceInfo,omitempty"`
	Places            interface{}            `json:"places,omitempty"`
	Citizenship       interface{}            `json:"citizenship,omitempty"`
	OrderStatus       string                 `json:"orderStatus,omitempty"`
	Owner             string                 `json:"owner,omitempty"`
	Errors            interface{}            `json:"errors,omitempty"`
	UserInfo          interface{}            `json:"user_info,omitempty"`
	CustomerDocuments interface{}            `json:"customer_documents,omitempty"`
}

type CheckBookingAviaBusMDL struct {
	BeginDate              time.Time                        `json:"beginDate,omitempty"`
	Routes                 OrderRoutes                      `json:"routes"`
	TripIDToServiceID      map[string]string                `json:"tripIdToServiceId"`
	Trips                  map[string]MaasTripAviaBus       `json:"trips,omitempty"`
	Services               map[string]CBServicesAviaBus     `json:"services,omitempty"`
	AvailableDocumentTypes []string                         `json:"availableDocumentTypes"`
	Customers              map[int32]CheckBookingABCustomer `json:"customers"`
	PriceInfo              interface{}                      `json:"priceInfo,omitempty"`
	Places                 interface{}                      `json:"places,omitempty"`
	UID                    string                           `json:"uid"`
	OrderStatus            string                           `json:"orderStatus,omitempty"`
	Owner                  string                           `json:"owner,omitempty"`
	Errors                 interface{}                      `json:"errors,omitempty"`
	CustomerIDs            []int32                          `json:"customerIds"`
}

type TripIDToServiceID struct {
	ServiceID string `json:"serviceId"`
	TripID    string `json:"tripId"`
}

type GetOrderByEmailReq struct {
	Email   string `json:"email"`
	OrderID int    `json:"orderId"`
}

type GetOrderResp struct {
	Order struct {
		Bookings []struct {
			BookDocuments []struct {
				BookDocumentCustomers []struct {
					CustomerID    int32 `json:"customerId"`
					CustomerPrice struct {
						CurrencyCode string        `json:"currencyCode"`
						CustomerID   int           `json:"customerId"`
						EquiveTariff float64       `json:"equiveTariff"`
						SellerPrice  float64       `json:"sellerPrice"`
						Tariff       float64       `json:"tariff"`
						Taxes        []interface{} `json:"taxes"`
						Vat          float64       `json:"vat"`
					} `json:"customerPrice"`
					ID           int      `json:"id"`
					SeatRequired bool     `json:"seatRequired"`
					TariffInfo   struct{} `json:"tariffInfo"`
				} `json:"bookDocumentCustomers"`
				BookDocumentStatus string `json:"bookDocumentStatus"`
				ContainerData      struct {
					TicketData struct {
						SeatNumber string `json:"seatNumber"`
					} `json:"ticketData"`
				} `json:"containerData"`
				Refund       interface{} `json:"refund,omitempty"`
				ID           int         `json:"id"`
				Number       string      `json:"number"`
				SeatRequired bool        `json:"seatRequired"`
				TripIds      []string    `json:"tripIds"`
			} `json:"bookDocuments"`
			BookingStatus       string `json:"bookingStatus"`
			CustomerDocumentIds []int  `json:"customerDocumentIds"`
			DateTimeConfirm     string `json:"dateTimeConfirm"`
			DateTimeCreate      string `json:"dateTimeCreate"`
			DateTimeUpdate      string `json:"dateTimeUpdate"`
			ID                  int    `json:"id"`
			PriceInfo           struct {
				AgencyVAT      float64 `json:"agencyVAT"`
				Commission     float64 `json:"commission"`
				CurrencyCode   string  `json:"currencyCode"`
				CustomerPrices []struct {
					CurrencyCode string        `json:"currencyCode"`
					CustomerID   int           `json:"customerId"`
					EquiveTariff float64       `json:"equiveTariff"`
					SellerPrice  float64       `json:"sellerPrice"`
					Tariff       float64       `json:"tariff"`
					Taxes        []interface{} `json:"taxes"`
					Vat          float64       `json:"vat"`
				} `json:"customerPrices"`
				Fee         float64 `json:"fee"`
				Price       float64 `json:"price"`
				PriceStatus string  `json:"priceStatus"`
				Reward      float64 `json:"reward"`
				SellerPrice float64 `json:"sellerPrice"`
			} `json:"priceInfo"`
			PrivateData struct {
				PolicyID string `json:"policyId"`
			} `json:"privateData"`
			ServiceID string `json:"serviceId"`
		} `json:"bookings"`
		BeginDate     time.Time  `json:"beginDate"`
		CustomerIds   []int      `json:"customerIds"`
		Customers     []Customer `json:"customers"`
		ID            int        `json:"id"`
		OrderStatus   string     `json:"orderStatus"`
		Owner         string     `json:"owner"`
		DateTimeToPay string     `json:"dateTimeToPay"`
		DocumentForms []struct {
			BookingID          int    `json:"bookingId"`
			DocumentFormFormat string `json:"documentFormFormat"`
			DocumentFormType   string `json:"documentFormType"`
			FileStorageID      string `json:"fileStorageId,omitempty"`
			ID                 int    `json:"id"`
			URLDataEnd         string `json:"urlDataEnd"`
		} `json:"documentForms"`
		PriceInfo struct {
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
			ContainerData struct {
				TripsData []struct {
					GeneralData struct {
						NeedPrintDocument bool `json:"needPrintDocument"`
						Options           []struct {
							Code string `json:"code"`
							Name string `json:"name"`
						} `json:"options"`
					} `json:"generalData"`
					SeatsData struct {
						FreeSeatsCount int64 `json:"freeSeatsCount"`
					} `json:"seatsData"`
				} `json:"tripsData"`
			} `json:"containerData"`
			ID                  string   `json:"id"`
			ProviderServiceCode string   `json:"providerServiceCode"`
			SellingType         string   `json:"sellingType"`
			ServiceType         string   `json:"serviceType"`
			TripIds             []string `json:"tripIds"`
			PriceInfo           struct {
				AgencyVAT      float64 `json:"agencyVAT"`
				Commission     float64 `json:"commission"`
				CurrencyCode   string  `json:"currencyCode"`
				CustomerPrices []struct {
					CurrencyCode string        `json:"currencyCode"`
					CustomerID   int           `json:"customerId"`
					EquiveTariff float64       `json:"equiveTariff"`
					SellerPrice  float64       `json:"sellerPrice"`
					Tariff       float64       `json:"tariff"`
					Taxes        []interface{} `json:"taxes"`
					Vat          float64       `json:"vat"`
				} `json:"customerPrices"`
				Fee         float64     `json:"fee"`
				Price       float64     `json:"price"`
				PriceStatus string      `json:"priceStatus"`
				Reward      float64     `json:"reward"`
				SellerPrice float64     `json:"sellerPrice"`
				Refund      interface{} `json:"refund,omitempty"`
			} `json:"priceInfo"`
		} `json:"services"`
		Trips  []MaasTrip `json:"trips"`
		UserID int        `json:"userId"`
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
	Places interface{} `json:"places"`
}

type Customer struct {
	SeatRequired      bool               `json:"seatRequired,omitempty"`
	Age               int                `json:"age,omitempty"`
	ID                int32              `json:"id,omitempty"`
	Birthdate         string             `json:"birthDate,omitempty"`
	Sex               string             `json:"sex, omitempty"`
	CustomerDocuments []CustomerDocument `json:"customerDocuments,omitempty"`
	Phone             string             `json:"phone,omitempty "`
	Email             string             `json:"email,omitempty"`
	UserID            int                `json:"userId,omitempty"`
}

type MaasTrip struct {
	Arrival            time.Time         `json:"arrival"`
	BoardingDuration   int               `json:"boardingDuration"`
	CarrierName        string            `json:"carrierName"`
	Departure          time.Time         `json:"departure"`
	Descr              string            `json:"descr"`
	Direction          string            `json:"direction"`
	Distance           int               `json:"distance"`
	Duration           int               `json:"duration"`
	FromDescr          string            `json:"fromDescr"`
	FromID             int               `json:"fromId"`
	ID                 string            `json:"id"`
	IsReturnTrip       bool              `json:"isReturnTrip"`
	IsTransfer         bool              `json:"isTransfer"`
	ToDescr            string            `json:"toDescr"`
	ToID               int               `json:"toId"`
	TransportClass     string            `json:"transportClass"`
	TripType           string            `json:"tripType"`
	UnboardingDuration int               `json:"unboardingDuration"`
	Options            map[string]string `json:"options,omitempty"`
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

type PriceInfoToCustomerOrder struct {
	SellerPrice float64 `json:"sellerPrice"`
	Fee         float64 `json:"fee"`
}

type OrderCustomer struct {
	ID          int32                  `json:"id, omitempty"`
	Birthdate   string                 `json:"birthDate, omitempty"`
	Sex         string                 `json:"sex, omitempty"`
	SellerPrice float64                `json:"sellerPrice,omitempty"`
	Fee         float64                `json:"fee,omitempty"`
	Citizenship dictionary.Citizenship `json:"citizenship"`
	ExpireDate  string                 `json:"expireDate"`
	Phone       string                 `json:"phone, omitempty "`
	Email       string                 `json:"email, omitempty"`
	FirstName   string                 `json:"firstName"`
	Middlename  string                 `json:"middleName,omitempty"`
	LastName    string                 `json:"lastName"`
	Number      string                 `json:"number"`
	Type        string                 `json:"type"`
}

type CheckBookingABCustomer struct {
	ID           int32  `json:"id"`
	Sex          string `json:"sex"`
	SeatRequired bool   `json:"seatRequired"`
	Age          int    `json:"age"`
	UserID       int    `json:"userId"`
}

type GetOrderMaasResp struct {
	OrderRoutes       OrderRoutes             `json:"routes"`
	TripIDToServiceID map[string]string       `json:"tripIdToServiceId"`
	Trips             map[string]MaasTrip     `json:"trips"`
	Services          map[string]interface{}  `json:"services"`
	Places            interface{}             `json:"places"`
	SearchParams      interface{}             `json:"searchParams"`
	CustomerIDs       []int32                 `json:"customerIds"`
	Customers         map[int32]OrderCustomer `json:"customers"`
	OrderStatus       string                  `json:"orderStatus"`
	OrderID           int                     `json:"orderId"`
	DateTimeToPay     string                  `json:"dateTimeToPay"`
	IsOwnedByUser     bool                    `json:"isOwnedByUser"`
	Bookings          map[int]OrderBooking    `json:"bookings"`
	Error             models.UapiError        `json:"errors,omitempty"`
}

type OrderServices struct {
	ID                  string              `json:"id"`
	ProviderServiceCode string              `json:"providerServiceCode"`
	SellingType         string              `json:"sellingType"`
	ServiceType         string              `json:"serviceType"`
	TripIds             []string            `json:"tripIds"`
	PriceByCustomer     map[int]SellerPrice `json:"priceByCustomerId"`
	Fee                 float64             `json:"fee"`
	Price               float64             `json:"price"`
	Refund              interface{}         `json:"refund,omitempty"`
	SellerPrice         float64             `json:"sellerPrice"`
}

type CBServicesAviaBus struct {
	Alternatives        map[int]CBAlternative  `json:"alternatives,omitempty"`
	ID                  string                 `json:"id"`
	ProviderServiceCode string                 `json:"providerServiceCode"`
	SellingType         string                 `json:"sellingType"`
	ServiceType         string                 `json:"serviceType"`
	TripIds             []string               `json:"tripIds"`
	PriceByCustomer     map[int]SellerPrice    `json:"priceByCustomerId,omitempty"`
	ObjectType          string                 `json:"objectType,omitempty"`
	Fee                 float64                `json:"fee,omitempty"`
	Price               float64                `json:"price,omitempty"`
	SellerPrice         float64                `json:"sellerPrice,omitempty"`
	BusByTripID         map[string]BusByTripID `json:"busByTripId,omitempty"`
}

type CBAlternative struct {
	SegmentByTripID   map[string]CBSegment `json:"segmentByTripId,omitempty"`
	PriceByCustomerID map[int]interface{}  `json:"priceByCustomerId"`
	Fee               float64              `json:"fee"`
	Price             float64              `json:"price"`
	PriceStatus       string               `json:"priceStatus"`
	Reward            float64              `json:"reward"`
	SellerPrice       float64              `json:"sellerPrice"`
	AgencyVAT         float64              `json:"agencyVAT"`
	Commission        float64              `json:"commission"`
	CurrencyCode      string               `json:"currencyCode"`
	ID                int                  `json:"id"`
}

type CBSegment struct {
	ComfortType      string                `json:"comfortType"`
	FareByCustomerID map[int64]interface{} `json:"fareByCustomerId"`
	FareFamily       bool                  `json:"fareFamily"`
	FareFamilyName   string                `json:"fareFamilyName"`
	FreeSeatsCount   int64                 `json:"freeSeatsCount"`
	MarkCarrierCode  string                `json:"markCarrierCode"`
	MarkCarrierName  string                `json:"markCarrierName"`
	OpCarrierCode    string                `json:"opCarrierCode"`
	Options          interface{}           `json:"options"`
	ValCarrierCode   string                `json:"valCarrierCode"`
	ValCarrierName   string                `json:"valCarrierName"`
}

type CBFare struct {
	BookingClass   string      `json:"bookingClass"`
	Code           string      `json:"code"`
	ComfortType    string      `json:"comfortType"`
	FareCalcLine   string      `json:"fareCalcLine"`
	FareFamily     bool        `json:"fareFamily"`
	FareFamilyName string      `json:"fareFamilyName"`
	Options        interface{} `json:"options"`
}

type BusByTripID struct {
	FreeSeats         interface{} `json:"freeSeats"`
	FreeSeatsCount    int64       `json:"freeSeatsCount"`
	Options           interface{} `json:"options"`
	ID                string      `json:"id"`
	NeedPrintDocument bool        `json:"needPrintDocument"`
}

type SellerPrice struct {
	SellerPrice float64 `json:"sellerPrice"`
}

type OrderBooking struct {
	SeatByCustomerID   map[int32]OrderSeatInfo `json:"seatByCustomerId"`
	BookDocumentStatus string                  `json:"bookDocumentStatus"`
	ID                 int                     `json:"id"`
	TripIDs            []string                `json:"tripIds"`
	Refund             interface{}             `json:"refund,omitempty"`
	ServiceID          string                  `json:"serviceId"`
	FileStorageID      string                  `json:"fileStorageId,omitempty"`
}

type OrderRoutes struct {
	Forward  []string `json:"forward,omitempty"`
	Backward []string `json:"backward,omitempty"`
}

type OrderSeatInfo struct {
	EquiveTariff float64 `json:"equiveTariff"`
	SellerPrice  float64 `json:"sellerPrice"`
	SeatRequired bool    `json:"seatRequired"`
	SeatNumber   string  `json:"seatNumber,omitempty"`
}

type CheckBookingAviaBusResponse struct {
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
				CurrencyCode string        `json:"currencyCode"`
				CustomerID   float64       `json:"customerId"`
				EquiveTariff float64       `json:"equiveTariff"`
				SellerPrice  float64       `json:"sellerPrice"`
				Tariff       float64       `json:"tariff"`
				Taxes        []interface{} `json:"taxes"`
				Vat          float64       `json:"vat"`
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
				Commission     float64 `json:"commission"`
				CurrencyCode   string  `json:"currencyCode"`
				CustomerPrices []struct {
					CurrencyCode string      `json:"currencyCode"`
					CustomerID   int         `json:"customerId"`
					EquiveTariff float64     `json:"equiveTariff"`
					SellerPrice  float64     `json:"sellerPrice"`
					Tariff       float64     `json:"tariff"`
					Taxes        interface{} `json:"taxes"`
					Vat          float64     `json:"vat"`
				} `json:"customerPrices"`
				Fee         float64 `json:"fee"`
				Price       float64 `json:"price"`
				PriceStatus string  `json:"priceStatus"`
				Reward      float64 `json:"reward"`
				SellerPrice float64 `json:"sellerPrice"`
			} `json:"alternativePriceInfos"`
			AvailableDocumentTypes []string `json:"availableDocumentTypes"`
			ContainerData          struct {
				AlternativesData []struct {
					Segments []struct {
						Fares []struct {
							BookingClass   string      `json:"bookingClass"`
							Code           string      `json:"code"`
							ComfortType    string      `json:"comfortType"`
							CustomerID     int64       `json:"customerId"`
							FareCalcLine   string      `json:"fareCalcLine"`
							FareFamily     bool        `json:"fareFamily"`
							FareFamilyName string      `json:"fareFamilyName"`
							Options        interface{} `json:"options"`
						} `json:"fares"`
					} `json:"segments"`
				} `json:"alternativesData"`
				GeneralData struct {
					FlightType string `json:"flightType"`
					Segments   []struct {
						FreeSeatsCount  int64  `json:"freeSeatsCount"`
						IsFareFamily    bool   `json:"isFareFamily"`
						MarkCarrierCode string `json:"markCarrierCode"`
						MarkCarrierName string `json:"markCarrierName"`
						OpCarrierCode   string `json:"opCarrierCode"`
						ValCarrierCode  string `json:"valCarrierCode"`
						ValCarrierName  string `json:"valCarrierName"`
					} `json:"segments"`
				} `json:"generalData,omitempty"`
				TripsData []struct {
					GeneralData struct {
						NeedPrintDocument bool        `json:"needPrintDocument"`
						Options           interface{} `json:"options"`
					} `json:"generalData"`
					SeatsData struct {
						FreeSeats []struct {
							Number string `json:"number"`
							SeatID string `json:"seatId"`
						} `json:"freeSeats"`
						FreeSeatsCount int64 `json:"freeSeatsCount"`
					} `json:"seatsData"`
				} `json:"tripsData"`
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
		Trips  []MaasTripAviaBus `json:"trips"`
		UserID int               `json:"userId"`
	} `json:"order"`
	Places       map[string]BookingPlace `json:"places"`
	SearchParams struct {
		AirServiceClass string `json:"airServiceClass"`
		CultureCode     string `json:"cultureCode"`
		CurrencyCode    string `json:"currencyCode"`
		Customers       []struct {
			Age          int  `json:"age"`
			ID           int  `json:"id"`
			SeatRequired bool `json:"seatRequired"`
		} `json:"customers"`
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
}

type MaasTripAviaBus struct {
	Arrival            time.Time   `json:"arrival"`
	BoardingDuration   int         `json:"boardingDuration"`
	CarrierName        string      `json:"carrierName"`
	Departure          time.Time   `json:"departure"`
	Descr              string      `json:"descr"`
	Direction          string      `json:"direction"`
	Distance           int         `json:"distance"`
	Duration           int         `json:"duration"`
	FromDescr          string      `json:"fromDescr"`
	FromID             int         `json:"fromId"`
	ID                 string      `json:"id"`
	IsReturnTrip       bool        `json:"isReturnTrip"`
	IsTransfer         bool        `json:"isTransfer"`
	ToDescr            string      `json:"toDescr"`
	ToID               int         `json:"toId"`
	TransportClass     string      `json:"transportClass"`
	TripType           string      `json:"tripType"`
	UnboardingDuration int         `json:"unboardingDuration"`
	Options            interface{} `json:"options,omitempty"`
}
