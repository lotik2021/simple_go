package models

type DeepLink struct {
	Id           string `json:"guid"`
	Ios          string `json:"ios"`
	Android      string `json:"android"`
	IosStore     string `json:"ios_store"`
	AndroidStore string `json:"android_store"`
	Title        string `json:"title,omitempty"`
	IconName     string `json:"icon_name,omitempty"`
}
