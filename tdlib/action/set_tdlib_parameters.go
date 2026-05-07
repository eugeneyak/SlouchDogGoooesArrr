package action

type SetTdlibParameters struct {
	Type               string `json:"@type"`
	APIID              int32  `json:"api_id"`
	APIHash            string `json:"api_hash"`
	SystemLanguageCode string `json:"system_language_code"`
	DeviceModel        string `json:"device_model"`
	ApplicationVersion string `json:"application_version"`
}
