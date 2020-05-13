package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var (
	C config
)

type (
	config struct {
		ENV      string
		LogLevel string
		Server   struct {
			URL string
		}
		RequestTimeout time.Duration
		ETCD           struct {
			URL      string
			BasePath string
			Paths    struct {
				PromptsPlacesTariffs string
				RateLimits           string
				BlackList            string
			}
		}
		CarrotQuestAuthKey string
		Redis              struct {
			URL string
		}
		DemoUser struct {
			Phone    string
			Code     string
			TokenTTL time.Duration
		}
		AsyncAvailable    bool
		SeatAutoSelection bool
		Database          struct {
			URL         string
			Poolsizemax int
		}
		YandexMoney struct {
			RequestTimeout time.Duration
			ClientID       string
			Urls           struct {
				StrelkaParameters string
				TroikaParameters  string
			}
		}
		Strelka struct {
			RequestTimeout time.Duration
			ClientID       string
			ApiAvailable   bool
			Urls           struct {
				DefaultRedirect string
				Redirect        string
				Pay             string
				Types           string
				Balance         string
			}
		}
		Troika struct {
			RequestTimeout time.Duration
			ClientID       string
			Urls           struct {
				FailRedirect           string
				SuccessRedirect        string
				RequestInstanceID      string
				ExternalPayment        string
				ProcessExternalPayment string
			}
		}
		Google struct {
			RequestTimeout time.Duration
			ApiKey         string
		}
		Taxi struct {
			RequestTimeout time.Duration
			CityMobile     struct {
				ApiKey   string
				DeepLink struct {
					PartnerId             string
					TrackingLinkMoscow    string
					TrackingLinkSamara    string
					TrackingLinkYaroslavl string
					TrackingLinkTolyatti  string
					Prefix                string
				}
				BaseURL string
				Urls    struct {
					Price string
					Auth  string
				}
			}
			Yandex struct {
				Clid     string
				RefId    string
				ApiKey   string
				DeepLink struct {
					Prefix     string
					TrackingId string
				}
				BaseURL string
				Urls    struct {
					Price string
				}
			}
			Maxim struct {
				Token    string
				AppCode  string
				DeepLink struct {
					TrackingLinkIos     string
					TrackingLinkAndroid string
					RefOrgID            string
					Prefix              string
				}
				BaseURL string
				Urls    struct {
					Price string
				}
			}
		}
		OpenWeather struct {
			BaseURL        string
			RequestTimeout time.Duration
			ApiKey         string
		}
		Auth struct {
			RequestTimeout  time.Duration
			TemporarySecret string
			Credentials     struct {
				ClientID     string
				ClientSecret string
				Internal     struct {
					GrantType string
					Scope     string
				}
				External struct {
					GrantType string
					Scope     string
				}
				RefreshToken struct {
					GrantType string
					Scope     string
				}
			}
			Core struct {
				BaseURL string
				Urls    struct {
					ConnectToken    string
					ConnectUserInfo string
				}
			}
			Account struct {
				BaseURL string
				Urls    struct {
					DevicesRegister      string
					Smscode              string
					Authorize            string
					CompleteRegistration string
					GetUser              string
					ConfirmEmail         string
				}
			}
		}
		Search struct {
			BaseURL string
			Urls    struct {
				RouteDetailsV3 string
			}
		}
		Analytics struct {
			Urls struct {
				Google  string
				Google2 string
			}
		}
		Web struct {
			BaseURL string
		}
		Content struct {
			BaseURL string
			Urls    struct {
				QueryCms string
			}
		}
		Notification struct {
			Provider     string
			SupportEmail string
			BaseURL      string
			Urls         struct {
				Feedback string
			}
		}
		WebTmp struct {
			BaseURL string
			Urls    struct {
				OrderEvent string
			}
		}
		Static struct {
			AboutCustomer string
		}
		Places struct {
			RequestTimeout time.Duration
			BaseURL        string
			Urls           struct {
				SearchPlaces      string
				SearchPlacesByGeo string
				SearchPlacesByIds string
			}
		}
		Dialog struct {
			RequestTimeout time.Duration
			BaseURL        string
			Urls           struct {
				Dialog string
			}
		}
		PersonalArea struct {
			RequestTimeout time.Duration
			BaseURL        string
			Urls           struct {
				DocumentsV2           string
				DocumentsDeactivateV2 string
				CustomersV2           string
			}
		}
		FapiAdapter struct {
			RequestTimeout time.Duration
			BaseURL        string
			Urls           struct {
				SearchAsyncV5        string
				GetSearchStatusV5    string
				GetSearchResultsV5   string
				GetPathGroupV5       string
				GetSegmentRoutesV5   string
				SaveSelectedRoutesV5 string
				SaveSelectedRoutesV4 string
				GetSelectedRoutesV5  string
				GetSelectedRoutesV4  string
				SearchV4             string
				CreateBookingV4      string
				CheckBookingV4       string
				GetOrderV2           string
				CreateOrderV4        string
				RefundOrderAsyncV2   string
				AddServicesToOrder   string
				CheckRefund          string
				Blank                string
				BlanksZip            string
				BlanksPdf            string
				GetOrderByEmail      string
			}
		}
		Push struct {
			RequestTimeout time.Duration
			BaseURL        string
			Urls           struct {
				Send                  string
				CarrotQuestEvent      string
				CarrotQuestProps      string
				CarrotQuestOrderEvent string
			}
		}
		Booking struct {
			RequestTimeout time.Duration
			BaseURL        string
			Urls           struct {
				GetFareRulesV2        string
				GetFareRulesV1        string
				CancelOrder           string
				CheckRefundV1         string
				CreateBookingV2       string
				ChangeERegistrationV1 string
			}
		}
		IpGeo struct {
			BaseURL string
		}
		Payment struct {
			RequestTimeout time.Duration
			BaseURL        string
			Processing     string
			Urls           struct {
				Pay             string
				OrderState      string
				GetOrderStateV1 string
				RefundOrderV1   string
				AdditionalPay   string
			}
		}
		Recaptcha struct {
			Enabled          bool
			BaseURL          string
			SecretKeyWeb     string
			SecretKeyIos     string
			SecretKeyAndroid string
		}
		AllowOrigins []string
		Icons        struct {
			TravelCards struct {
				Troika struct {
					Ios     string
					Android string
				}
				Strelka struct {
					Ios     string
					Android string
				}
			}
			Taxi struct {
				Yandex struct {
					Ios     string
					Android string
				}
				Maxim struct {
					Ios     string
					Android string
				}
				Citymobil struct {
					Ios     string
					Android string
				}
			}
			Driving struct {
				IconUrl struct {
					Ios     string
					Android string
				}
			}
			Trips struct {
				Avia struct {
					Ios     string
					Android string
				}
				Bus struct {
					Ios     string
					Android string
				}
				Train struct {
					Ios     string
					Android string
				}
			}
			Places struct {
				Airport struct {
					Ios     string
					Android string
				}
				Bus struct {
					Ios     string
					Android string
				}
				City struct {
					Ios     string
					Android string
				}
				Default struct {
					Ios     string
					Android string
				}
				Hub struct {
					Ios     string
					Android string
				}
				Home struct {
					Ios     string
					Android string
				}
				Metro struct {
					Ios     string
					Android string
				}
				Train struct {
					Ios     string
					Android string
				}
				Work struct {
					Ios     string
					Android string
				}
				DarkTheme struct {
					Airport struct {
						Ios     string
						Android string
					}
					Bus struct {
						Ios     string
						Android string
					}
					City struct {
						Ios     string
						Android string
					}
					Default struct {
						Ios     string
						Android string
					}
					Hub struct {
						Ios     string
						Android string
					}
					Home struct {
						Ios     string
						Android string
					}
					Metro struct {
						Ios     string
						Android string
					}
					Train struct {
						Ios     string
						Android string
					}
					Work struct {
						Ios     string
						Android string
					}
				}
			}
			GenericData struct {
				DetailsIcons struct {
					Bus struct {
						Ios     string
						Android string
					}
					Ferry struct {
						Ios     string
						Android string
					}
					Metro struct {
						Ios     string
						Android string
					}
					Train struct {
						Ios     string
						Android string
					}
					Tram struct {
						Ios     string
						Android string
					}
					Trolleybus struct {
						Ios     string
						Android string
					}
					Mcd struct {
						Ios     string
						Android string
					}
					Aeroexpress struct {
						Ios     string
						Android string
					}
				}
				ShortIcons struct {
					Bus struct {
						Ios     string
						Android string
					}
					ShareTaxi struct {
						Ios     string
						Android string
					}
					Train struct {
						Ios     string
						Android string
					}
					Tram struct {
						Ios     string
						Android string
					}
					Trolleybus struct {
						Ios     string
						Android string
					}
					CommuterTrain struct {
						Ios     string
						Android string
					}
					McdD1 struct {
						Ios     string
						Android string
					}
					McdD2 struct {
						Ios     string
						Android string
					}
					Aeroexpress struct {
						Ios     string
						Android string
					}
					Ferry struct {
						Ios     string
						Android string
					}
					MetroDefault struct {
						Ios     string
						Android string
					}
				}
			}
			Weather struct {
				ClearSky struct {
					Day struct {
						Ios     string
						Android string
					}
					Night struct {
						Ios     string
						Android string
					}
				}
				FewClouds struct {
					Day struct {
						Ios     string
						Android string
					}
					Night struct {
						Ios     string
						Android string
					}
				}
				ScatteredClouds struct {
					Day struct {
						Ios     string
						Android string
					}
					Night struct {
						Ios     string
						Android string
					}
				}
				BrokenClouds struct {
					Day struct {
						Ios     string
						Android string
					}
					Night struct {
						Ios     string
						Android string
					}
				}
				ShowerRain struct {
					Day struct {
						Ios     string
						Android string
					}
					Night struct {
						Ios     string
						Android string
					}
				}
				Rain struct {
					Day struct {
						Ios     string
						Android string
					}
					Night struct {
						Ios     string
						Android string
					}
				}
				ThunderStorm struct {
					Day struct {
						Ios     string
						Android string
					}
					Night struct {
						Ios     string
						Android string
					}
				}
				Snow struct {
					Day struct {
						Ios     string
						Android string
					}
					Night struct {
						Ios     string
						Android string
					}
				}
				Mist struct {
					Day struct {
						Ios     string
						Android string
					}
					Night struct {
						Ios     string
						Android string
					}
				}
			}
		}
	}
)

func init() {
	C.ENV = strings.ToLower(os.Getenv("ENVIRONMENT"))
	if C.ENV == "" {
		C.ENV = "local"
	}

	viper.AddConfigPath("./configs")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config.default")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Running with %s configuration\n", C.ENV)
	viper.SetConfigName(fmt.Sprintf("config.%s", C.ENV))
	err = viper.MergeInConfig()
	if err != nil {
		log.Fatal(err)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err = viper.Unmarshal(&C)
	if err != nil {
		log.Fatal(err)
	}

	C.Places.Urls.SearchPlaces = C.Places.BaseURL + C.Places.Urls.SearchPlaces
	C.Places.Urls.SearchPlacesByGeo = C.Places.BaseURL + C.Places.Urls.SearchPlacesByGeo
	C.Places.Urls.SearchPlacesByIds = C.Places.BaseURL + C.Places.Urls.SearchPlacesByIds
	C.Dialog.Urls.Dialog = C.Dialog.BaseURL + C.Dialog.Urls.Dialog

	C.Auth.Core.Urls.ConnectToken = C.Auth.Core.BaseURL + C.Auth.Core.Urls.ConnectToken
	C.Auth.Core.Urls.ConnectUserInfo = C.Auth.Core.BaseURL + C.Auth.Core.Urls.ConnectUserInfo
	C.Auth.Account.Urls.DevicesRegister = C.Auth.Account.BaseURL + C.Auth.Account.Urls.DevicesRegister
	C.Auth.Account.Urls.Authorize = C.Auth.Account.BaseURL + C.Auth.Account.Urls.Authorize
	C.Auth.Account.Urls.Smscode = C.Auth.Account.BaseURL + C.Auth.Account.Urls.Smscode
	C.Auth.Account.Urls.CompleteRegistration = C.Auth.Account.BaseURL + C.Auth.Account.Urls.CompleteRegistration
	C.Auth.Account.Urls.GetUser = C.Auth.Account.BaseURL + C.Auth.Account.Urls.GetUser
	C.Auth.Account.Urls.ConfirmEmail = C.Auth.Account.BaseURL + C.Auth.Account.Urls.ConfirmEmail

	C.Strelka.Urls.Redirect = C.Web.BaseURL + C.Strelka.Urls.Redirect
	C.Troika.Urls.SuccessRedirect = C.Web.BaseURL + C.Troika.Urls.SuccessRedirect
	C.Troika.Urls.FailRedirect = C.Web.BaseURL + C.Troika.Urls.FailRedirect

	C.Taxi.CityMobile.Urls.Auth = C.Taxi.CityMobile.BaseURL + C.Taxi.CityMobile.Urls.Auth
	C.Taxi.CityMobile.Urls.Price = C.Taxi.CityMobile.BaseURL + C.Taxi.CityMobile.Urls.Price
	C.Taxi.Yandex.Urls.Price = C.Taxi.Yandex.BaseURL + C.Taxi.Yandex.Urls.Price
	C.Taxi.Maxim.Urls.Price = C.Taxi.Maxim.BaseURL + C.Taxi.Maxim.Urls.Price

	C.FapiAdapter.Urls.SearchAsyncV5 = C.FapiAdapter.BaseURL + C.FapiAdapter.Urls.SearchAsyncV5
	C.FapiAdapter.Urls.GetSearchStatusV5 = C.FapiAdapter.BaseURL + C.FapiAdapter.Urls.GetSearchStatusV5
	C.FapiAdapter.Urls.GetSearchResultsV5 = C.FapiAdapter.BaseURL + C.FapiAdapter.Urls.GetSearchResultsV5
	C.FapiAdapter.Urls.GetPathGroupV5 = C.FapiAdapter.BaseURL + C.FapiAdapter.Urls.GetPathGroupV5
	C.FapiAdapter.Urls.GetSegmentRoutesV5 = C.FapiAdapter.BaseURL + C.FapiAdapter.Urls.GetSegmentRoutesV5
	C.FapiAdapter.Urls.SaveSelectedRoutesV5 = C.FapiAdapter.BaseURL + C.FapiAdapter.Urls.SaveSelectedRoutesV5
	C.FapiAdapter.Urls.SaveSelectedRoutesV4 = C.FapiAdapter.BaseURL + C.FapiAdapter.Urls.SaveSelectedRoutesV4
	C.FapiAdapter.Urls.GetSelectedRoutesV5 = C.FapiAdapter.BaseURL + C.FapiAdapter.Urls.GetSelectedRoutesV5
	C.FapiAdapter.Urls.GetSelectedRoutesV4 = C.FapiAdapter.BaseURL + C.FapiAdapter.Urls.GetSelectedRoutesV4
	C.FapiAdapter.Urls.SearchV4 = C.FapiAdapter.BaseURL + C.FapiAdapter.Urls.SearchV4
	C.FapiAdapter.Urls.CreateBookingV4 = C.FapiAdapter.BaseURL + C.FapiAdapter.Urls.CreateBookingV4
	C.FapiAdapter.Urls.CheckBookingV4 = C.FapiAdapter.BaseURL + C.FapiAdapter.Urls.CheckBookingV4
	C.FapiAdapter.Urls.GetOrderV2 = C.FapiAdapter.BaseURL + C.FapiAdapter.Urls.GetOrderV2
	C.FapiAdapter.Urls.CreateOrderV4 = C.FapiAdapter.BaseURL + C.FapiAdapter.Urls.CreateOrderV4
	C.FapiAdapter.Urls.AddServicesToOrder = C.FapiAdapter.BaseURL + C.FapiAdapter.Urls.AddServicesToOrder
	C.FapiAdapter.Urls.CheckRefund = C.FapiAdapter.BaseURL + C.FapiAdapter.Urls.CheckRefund
	C.FapiAdapter.Urls.GetOrderByEmail = C.FapiAdapter.BaseURL + C.FapiAdapter.Urls.GetOrderByEmail

	C.PersonalArea.Urls.DocumentsV2 = C.PersonalArea.BaseURL + C.PersonalArea.Urls.DocumentsV2

	C.PersonalArea.Urls.DocumentsDeactivateV2 = C.PersonalArea.BaseURL + C.PersonalArea.Urls.DocumentsDeactivateV2

	C.PersonalArea.Urls.CustomersV2 = C.PersonalArea.BaseURL + C.PersonalArea.Urls.CustomersV2

	C.FapiAdapter.Urls.Blank = C.FapiAdapter.BaseURL + C.FapiAdapter.Urls.Blank
	C.FapiAdapter.Urls.BlanksPdf = C.FapiAdapter.BaseURL + C.FapiAdapter.Urls.BlanksPdf
	C.FapiAdapter.Urls.BlanksZip = C.FapiAdapter.BaseURL + C.FapiAdapter.Urls.BlanksZip

	C.Push.Urls.Send = C.Push.BaseURL + C.Push.Urls.Send
	C.Push.Urls.CarrotQuestEvent = C.Push.BaseURL + C.Push.Urls.CarrotQuestEvent
	C.Push.Urls.CarrotQuestProps = C.Push.BaseURL + C.Push.Urls.CarrotQuestProps
	C.Push.Urls.CarrotQuestOrderEvent = C.Push.BaseURL + C.Push.Urls.CarrotQuestOrderEvent

	C.Booking.Urls.GetFareRulesV2 = C.Booking.BaseURL + C.Booking.Urls.GetFareRulesV2
	C.Booking.Urls.GetFareRulesV1 = C.Booking.BaseURL + C.Booking.Urls.GetFareRulesV1
	C.Booking.Urls.CancelOrder = C.Booking.BaseURL + C.Booking.Urls.CancelOrder
	C.Booking.Urls.CreateBookingV2 = C.Booking.BaseURL + C.Booking.Urls.CreateBookingV2
	C.Booking.Urls.CheckRefundV1 = C.Booking.BaseURL + C.Booking.Urls.CheckRefundV1
	C.Booking.Urls.ChangeERegistrationV1 = C.Booking.BaseURL + C.Booking.Urls.ChangeERegistrationV1

	C.ETCD.Paths.PromptsPlacesTariffs = C.ETCD.BasePath + C.ETCD.Paths.PromptsPlacesTariffs
	C.ETCD.Paths.RateLimits = C.ETCD.BasePath + C.ETCD.Paths.RateLimits
	C.ETCD.Paths.BlackList = C.ETCD.BasePath + C.ETCD.Paths.BlackList

	C.Payment.Urls.Pay = C.Payment.BaseURL + C.Payment.Urls.Pay
	C.Payment.Urls.OrderState = C.Payment.BaseURL + C.Payment.Urls.OrderState
	C.Payment.Urls.AdditionalPay = C.Payment.BaseURL + C.Payment.Urls.AdditionalPay
	C.Payment.Urls.GetOrderStateV1 = C.Payment.BaseURL + C.Payment.Urls.GetOrderStateV1
	C.Payment.Urls.RefundOrderV1 = C.Payment.BaseURL + C.Payment.Urls.RefundOrderV1

	C.Content.Urls.QueryCms = C.Content.BaseURL + C.Content.Urls.QueryCms

	C.WebTmp.Urls.OrderEvent = C.WebTmp.BaseURL + C.WebTmp.Urls.OrderEvent

	C.Notification.Urls.Feedback = C.Notification.BaseURL + C.Notification.Urls.Feedback

	C.Search.Urls.RouteDetailsV3 = C.Search.BaseURL + C.Search.Urls.RouteDetailsV3

	LoadRemote()
}

func IsProd() bool {
	return C.ENV == "prod"
}
