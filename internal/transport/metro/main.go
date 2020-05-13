package metro

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/logger"
)

var (
	metroData = make([]Metro, 0)
)

func Init() {
	db := common.GetDatabaseConnection()
	if err := db.Model(&metroData).Select(); err != nil {
		logger.Log.Fatal(err)
	}
}

type Metro struct {
	tableName      struct{} `pg:"maasapi.metro,alias:metro"`
	ID             int      `pg:"id,pk" json:"id"`
	Language       string   `pg:"language" json:"language"`
	CityName       string   `pg:"city_name" json:"city_name"`
	MovistaCityID  int      `pg:"movista_city_id" json:"movista_city_id"`
	LogoURI        string   `pg:"logo_uri" json:"logo_uri" `
	LineNumber     string   `pg:"line_number" json:"line_number"`
	FullName       string   `pg:"full_name" json:"full_name"`
	ShortName      string   `pg:"short_name" json:"short_name"`
	ColorHex       string   `pg:"color_hex" json:"color_hex"`
	AgencyUrl      string   `pg:"agency_url" json:"agency_uri"`
	AgencyName     string   `pg:"agency_name" json:"agency_name"`
	Type           string   `pg:"type" json:"type"`
	IosIconUrl     string   `pg:"ios_icon_url" json:"ios_icon_url"`
	AndroidIconUrl string   `pg:"android_icon_url" json:"android_icon_url"`
}

func (m *Metro) GetIcons() {

}

func GetLineColor(lineName string, lineShortName string, lineurl string) (lineColor, movistaLineShortName, metroIcon string, androidMetroIcon string) {
	for _, color := range metroData {
		if lineName == color.ShortName || (lineShortName == color.LineNumber && lineurl == color.AgencyUrl) {
			lineColor = color.ColorHex
			movistaLineShortName = color.ShortName
			metroIcon = color.IosIconUrl
			androidMetroIcon = color.AndroidIconUrl
			return
		}
	}
	return
}
