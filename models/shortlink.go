package models

import (
	"time"
)

type Shortlink struct {
	ID            int64     `json:"id"`
	Code          string    `json:"code"`
	OriginalURL   string    `json:"original_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	IsDeleted     bool      `json:"is_deleted"`
	IsDisabled    bool      `json:"is_disabled"`
	TotalVisits   int       `json:"total_visits"`
	TotalDuration int       `json:"total_duration"`
}

type ShortlinkCreateRequest struct {
	URL string `json:"url" binding:"required,url"`
}

type ShortlinkResponse struct {
	ShortURL    string    `json:"short_url"`
	Code        string    `json:"code"`
	OriginalURL string    `json:"original_url"`
	CreatedAt   time.Time `json:"created_at"`
}

type ShortlinkListResponse struct {
	Items    []Shortlink `json:"items"`
	Total    int         `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}
