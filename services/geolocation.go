package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"location-tracking-shortlink/config"
	"location-tracking-shortlink/models"
)

type GeoService struct {
	ipAPIURL string
}

func NewGeoService() *GeoService {
	cfg := config.Get()
	ipAPIURL := "http://ip-api.com/json/"
	if cfg != nil && cfg.Geo.IPAPIURL != "" {
		ipAPIURL = cfg.Geo.IPAPIURL
	}
	return &GeoService{
		ipAPIURL: ipAPIURL,
	}
}

type IPAPIResponse struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	ISP         string  `json:"isp"`
	Org         string  `json:"org"`
	AS          string  `json:"as"`
	Query       string  `json:"query"`
}

func (s *GeoService) GetLocationByIP(ip string) (*models.GeoLocation, error) {
	url := fmt.Sprintf("%s%s", s.ipAPIURL, ip)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("IP API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ipResp IPAPIResponse
	if err := json.Unmarshal(body, &ipResp); err != nil {
		return nil, err
	}

	if ipResp.Status != "success" {
		return &models.GeoLocation{
			Precision: "ip",
		}, nil
	}

	return &models.GeoLocation{
		Latitude:  ipResp.Lat,
		Longitude: ipResp.Lon,
		Country:   ipResp.Country,
		City:      ipResp.City,
		Precision: "ip",
	}, nil
}

func (s *GeoService) CombineLocation(gpsLat, gpsLng *float64, ipLat, ipLng float64, country, city string) *models.GeoLocation {
	if gpsLat != nil && gpsLng != nil {
		return &models.GeoLocation{
			Latitude:  *gpsLat,
			Longitude: *gpsLng,
			Country:   country,
			City:      city,
			Precision: "gps",
		}
	}

	if ipLat != 0 && ipLng != 0 {
		return &models.GeoLocation{
			Latitude:  ipLat,
			Longitude: ipLng,
			Country:   country,
			City:      city,
			Precision: "ip",
		}
	}

	return nil
}
