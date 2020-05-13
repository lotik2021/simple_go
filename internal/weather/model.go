package weather

import "bitbucket.movista.ru/maas/maasapi/internal/config"

type ForecastResponse struct {
	Code    string                 `json:"cod"`
	Message int                    `json:"message"`
	List    []*OpenWeatherResponse `json:"list"`
	Cnt     int                    `json:"cnt"`
}

type OpenWeatherResponse struct {
	Base   string `json:"base"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Cod   int `json:"cod"`
	Coord struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"coord"`
	Dt   int64 `json:"dt"`
	ID   int   `json:"id"`
	Main struct {
		FeelsLike float64 `json:"feels_like"`
		Humidity  int     `json:"humidity"`
		Pressure  int     `json:"pressure"`
		Temp      float64 `json:"temp"`
		TempMax   float64 `json:"temp_max"`
		TempMin   float64 `json:"temp_min"`
	} `json:"main"`
	Name string `json:"name"`
	Sys  struct {
		Country string `json:"country"`
		ID      int    `json:"id"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
		Type    int    `json:"type"`
	} `json:"sys"`
	Timezone   int `json:"timezone"`
	Visibility int `json:"visibility"`
	Weather    []struct {
		Description string `json:"description"`
		Icon        string `json:"icon"`
		ID          int    `json:"id"`
		Main        string `json:"main"`
	} `json:"weather"`
	Wind struct {
		Deg   int     `json:"deg"`
		Speed float64 `json:"speed"`
	} `json:"wind"`
	IosIconURL     []string `json:"ios_icon_url"`
	AndroidIconURL []string `json:"android_icon_url"`
}

func (r *OpenWeatherResponse) GetIcons() {
	for _, w := range r.Weather {
		switch w.Icon {
		case ClearSkyDay:
			r.IosIconURL = append(r.IosIconURL, config.C.Icons.Weather.ClearSky.Day.Ios)
			r.AndroidIconURL = append(r.AndroidIconURL, config.C.Icons.Weather.ClearSky.Day.Android)
		case ClearSkyNight:
			r.IosIconURL = append(r.IosIconURL, config.C.Icons.Weather.ClearSky.Night.Ios)
			r.AndroidIconURL = append(r.AndroidIconURL, config.C.Icons.Weather.ClearSky.Night.Android)
		case FewCloudsDay:
			r.IosIconURL = append(r.IosIconURL, config.C.Icons.Weather.FewClouds.Day.Ios)
			r.AndroidIconURL = append(r.AndroidIconURL, config.C.Icons.Weather.FewClouds.Day.Android)
		case FewCloudsNight:
			r.IosIconURL = append(r.IosIconURL, config.C.Icons.Weather.FewClouds.Night.Ios)
			r.AndroidIconURL = append(r.AndroidIconURL, config.C.Icons.Weather.FewClouds.Night.Android)
		case ScatteredCloudsDay:
			r.IosIconURL = append(r.IosIconURL, config.C.Icons.Weather.ScatteredClouds.Day.Ios)
			r.AndroidIconURL = append(r.AndroidIconURL, config.C.Icons.Weather.ScatteredClouds.Day.Android)
		case ScatteredCloudsNight:
			r.IosIconURL = append(r.IosIconURL, config.C.Icons.Weather.ScatteredClouds.Night.Ios)
			r.AndroidIconURL = append(r.AndroidIconURL, config.C.Icons.Weather.ScatteredClouds.Night.Android)
		case BrokenCloudsDay:
			r.IosIconURL = append(r.IosIconURL, config.C.Icons.Weather.BrokenClouds.Day.Ios)
			r.AndroidIconURL = append(r.AndroidIconURL, config.C.Icons.Weather.BrokenClouds.Day.Android)
		case BrokenCloudsNight:
			r.IosIconURL = append(r.IosIconURL, config.C.Icons.Weather.BrokenClouds.Night.Ios)
			r.AndroidIconURL = append(r.AndroidIconURL, config.C.Icons.Weather.BrokenClouds.Night.Android)
		case ShowerRainDay:
			r.IosIconURL = append(r.IosIconURL, config.C.Icons.Weather.ShowerRain.Day.Ios)
			r.AndroidIconURL = append(r.AndroidIconURL, config.C.Icons.Weather.ShowerRain.Day.Android)
		case ShowerRainNight:
			r.IosIconURL = append(r.IosIconURL, config.C.Icons.Weather.ShowerRain.Night.Ios)
			r.AndroidIconURL = append(r.AndroidIconURL, config.C.Icons.Weather.ShowerRain.Night.Android)
		case RainDay:
			r.IosIconURL = append(r.IosIconURL, config.C.Icons.Weather.Rain.Day.Ios)
			r.AndroidIconURL = append(r.AndroidIconURL, config.C.Icons.Weather.Rain.Day.Android)
		case RainNight:
			r.IosIconURL = append(r.IosIconURL, config.C.Icons.Weather.Rain.Night.Ios)
			r.AndroidIconURL = append(r.AndroidIconURL, config.C.Icons.Weather.Rain.Night.Android)
		case ThunderstormDay:
			r.IosIconURL = append(r.IosIconURL, config.C.Icons.Weather.ThunderStorm.Day.Ios)
			r.AndroidIconURL = append(r.AndroidIconURL, config.C.Icons.Weather.ThunderStorm.Day.Android)
		case ThunderstormNight:
			r.IosIconURL = append(r.IosIconURL, config.C.Icons.Weather.ThunderStorm.Night.Ios)
			r.AndroidIconURL = append(r.AndroidIconURL, config.C.Icons.Weather.ThunderStorm.Night.Android)
		case SnowDay:
			r.IosIconURL = append(r.IosIconURL, config.C.Icons.Weather.Snow.Day.Ios)
			r.AndroidIconURL = append(r.IosIconURL, config.C.Icons.Weather.Snow.Day.Android)
		case SnowNight:
			r.IosIconURL = append(r.IosIconURL, config.C.Icons.Weather.Snow.Night.Ios)
			r.AndroidIconURL = append(r.AndroidIconURL, config.C.Icons.Weather.Snow.Night.Android)
		case MistDay:
			r.IosIconURL = append(r.IosIconURL, config.C.Icons.Weather.Mist.Day.Ios)
			r.AndroidIconURL = append(r.AndroidIconURL, config.C.Icons.Weather.Mist.Day.Android)
		case MistNight:
			r.IosIconURL = append(r.IosIconURL, config.C.Icons.Weather.Mist.Night.Ios)
			r.AndroidIconURL = append(r.AndroidIconURL, config.C.Icons.Weather.Mist.Night.Android)
		}
	}
}
