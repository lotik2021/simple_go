package weather

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"github.com/parnurzeal/gorequest"
)

const (
	ClearSkyDay          = "01d"
	ClearSkyNight        = "01n"
	FewCloudsDay         = "02d"
	FewCloudsNight       = "02n"
	ScatteredCloudsDay   = "03d"
	ScatteredCloudsNight = "03n"
	BrokenCloudsDay      = "04d"
	BrokenCloudsNight    = "04n"
	ShowerRainDay        = "09d"
	ShowerRainNight      = "09n"
	RainDay              = "10d"
	RainNight            = "10n"
	ThunderstormDay      = "11d"
	ThunderstormNight    = "11n"
	SnowDay              = "13d"
	SnowNight            = "13n"
	MistDay              = "50d"
	MistNight            = "50n"
)

var (
	units   = "metric" // For temperature in Celsius
	lang    = "ru"
	wclient *gorequest.SuperAgent
)

func init() {
	wclient = common.DefaultRequest.Clone().Timeout(config.C.OpenWeather.RequestTimeout)
}
