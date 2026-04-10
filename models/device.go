package models

type DeviceInfo struct {
	OS         string `json:"os_type"`
	OSVersion  string `json:"os_version"`
	Browser    string `json:"browser_type"`
	BrowserVer string `json:"browser_version"`
	DeviceType string `json:"device_type"`
}

type GeoLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Country   string  `json:"country"`
	City      string  `json:"city"`
	Precision string  `json:"precision"`
}

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
