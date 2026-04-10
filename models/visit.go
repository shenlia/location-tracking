package models

import (
	"database/sql"
	"time"
)

type Visit struct {
	ID             int64           `json:"id"`
	ShortlinkID    int64           `json:"shortlink_id"`
	IPAddress      string          `json:"ip_address"`
	Country        sql.NullString  `json:"country"`
	City           sql.NullString  `json:"city"`
	Latitude       sql.NullFloat64 `json:"latitude"`
	Longitude      sql.NullFloat64 `json:"longitude"`
	GeoPrecision   sql.NullString  `json:"geo_precision"`
	GeoStatus      sql.NullString  `json:"geo_status"`
	UserAgent      sql.NullString  `json:"user_agent"`
	OSType         sql.NullString  `json:"os_type"`
	OSVersion      sql.NullString  `json:"os_version"`
	BrowserType    sql.NullString  `json:"browser_type"`
	BrowserVersion sql.NullString  `json:"browser_version"`
	DeviceType     sql.NullString  `json:"device_type"`
	VisitDuration  sql.NullInt64   `json:"visit_duration"`
	VisitTime      time.Time       `json:"visit_time"`
	ExitTime       sql.NullTime    `json:"exit_time"`
	Referer        sql.NullString  `json:"referer"`
}

type VisitSubmitRequest struct {
	Code          string   `json:"code" binding:"required"`
	Latitude      *float64 `json:"latitude"`
	Longitude     *float64 `json:"longitude"`
	GeoPrecision  string   `json:"geo_precision"`
	GeoStatus     string   `json:"geo_status"`
	VisitDuration int64    `json:"visit_duration"`
	PageLoaded    bool     `json:"page_loaded"`
}

type VisitResponse struct {
	ID             int64    `json:"id"`
	IPAddress      string   `json:"ip_address"`
	Country        string   `json:"country"`
	City           string   `json:"city"`
	Latitude       *float64 `json:"latitude"`
	Longitude      *float64 `json:"longitude"`
	GeoPrecision   string   `json:"geo_precision"`
	GeoStatus      string   `json:"geo_status"`
	OSType         string   `json:"os_type"`
	OSVersion      string   `json:"os_version"`
	BrowserType    string   `json:"browser_type"`
	BrowserVersion string   `json:"browser_version"`
	DeviceType     string   `json:"device_type"`
	VisitDuration  int64    `json:"visit_duration"`
	VisitTime      string   `json:"visit_time"`
}

type VisitListResponse struct {
	Items    []VisitResponse `json:"items"`
	Total    int             `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
}
