package services

import (
	"database/sql"
	"time"

	"location-tracking-shortlink/db"
)

type StatsService struct{}

func NewStatsService() *StatsService {
	return &StatsService{}
}

type Stats struct {
	TotalVisits        int              `json:"total_visits"`
	TotalDuration      int64            `json:"total_duration"`
	AvgDuration        float64          `json:"avg_duration"`
	VisitTrend         []VisitTrendItem `json:"visit_trend"`
	GeoDistribution    []GeoPoint       `json:"geo_distribution"`
	DeviceDistribution map[string]int   `json:"device_distribution"`
}

type VisitTrendItem struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type GeoPoint struct {
	Lat   float64 `json:"lat"`
	Lng   float64 `json:"lng"`
	Count int     `json:"count"`
}

func (s *StatsService) GetStats() (*Stats, error) {
	var totalVisits int
	err := db.GetDB().QueryRow("SELECT COALESCE(SUM(total_visits), 0) FROM shortlinks WHERE is_deleted = FALSE").Scan(&totalVisits)
	if err != nil {
		return nil, err
	}

	var totalDuration int64
	err = db.GetDB().QueryRow("SELECT COALESCE(SUM(total_duration), 0) FROM shortlinks WHERE is_deleted = FALSE").Scan(&totalDuration)
	if err != nil {
		return nil, err
	}

	avgDuration := float64(0)
	if totalVisits > 0 {
		avgDuration = float64(totalDuration) / float64(totalVisits)
	}

	trend, err := s.getVisitTrend(7)
	if err != nil {
		return nil, err
	}

	geoDist, err := s.getGeoDistribution()
	if err != nil {
		return nil, err
	}

	deviceDist, err := s.getDeviceDistribution()
	if err != nil {
		return nil, err
	}

	return &Stats{
		TotalVisits:        totalVisits,
		TotalDuration:      totalDuration,
		AvgDuration:        avgDuration,
		VisitTrend:         trend,
		GeoDistribution:    geoDist,
		DeviceDistribution: deviceDist,
	}, nil
}

func (s *StatsService) getVisitTrend(days int) ([]VisitTrendItem, error) {
	var items []VisitTrendItem

	for i := days - 1; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")

		var count int
		err := db.GetDB().QueryRow(`
			SELECT COUNT(*) FROM visits 
			WHERE date(visit_time) = ?
		`, date).Scan(&count)
		if err != nil {
			return nil, err
		}

		items = append(items, VisitTrendItem{
			Date:  date,
			Count: count,
		})
	}

	return items, nil
}

func (s *StatsService) getGeoDistribution() ([]GeoPoint, error) {
	rows, err := db.GetDB().Query(`
		SELECT latitude, longitude, COUNT(*) as count 
		FROM visits 
		WHERE latitude IS NOT NULL AND longitude IS NOT NULL 
		GROUP BY ROUND(latitude, 2), ROUND(longitude, 2)
		LIMIT 100
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []GeoPoint
	for rows.Next() {
		var lat, lng float64
		var count int
		if err := rows.Scan(&lat, &lng, &count); err != nil {
			return nil, err
		}
		points = append(points, GeoPoint{
			Lat:   lat,
			Lng:   lng,
			Count: count,
		})
	}

	return points, nil
}

func (s *StatsService) getDeviceDistribution() (map[string]int, error) {
	dist := make(map[string]int)
	dist["pc"] = 0
	dist["mobile"] = 0
	dist["tablet"] = 0

	rows, err := db.GetDB().Query("SELECT device_type, COUNT(*) FROM visits GROUP BY device_type")
	if err != nil {
		if err == sql.ErrNoRows {
			return dist, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var deviceType string
		var count int
		if err := rows.Scan(&deviceType, &count); err != nil {
			continue
		}
		dist[deviceType] = count
	}

	return dist, nil
}
