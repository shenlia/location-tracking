package handlers

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"location-tracking-shortlink/db"
	"location-tracking-shortlink/models"
	"location-tracking-shortlink/services"
)

type AdminHandler struct {
	shortlinkService *services.ShortlinkService
	statsService     *services.StatsService
}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{
		shortlinkService: services.NewShortlinkService(),
		statsService:     services.NewStatsService(),
	}
}

func (h *AdminHandler) ListShortlinks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	shortlinks, total, err := h.shortlinkService.List(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Code:    500,
			Message: "Failed to list shortlinks: " + err.Error(),
		})
		return
	}

	if shortlinks == nil {
		shortlinks = []models.Shortlink{}
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    0,
		Message: "success",
		Data: models.ShortlinkListResponse{
			Items:    shortlinks,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

func (h *AdminHandler) DeleteShortlink(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    400,
			Message: "Invalid ID",
		})
		return
	}

	if err := h.shortlinkService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Code:    500,
			Message: "Failed to delete shortlink: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    0,
		Message: "success",
	})
}

func (h *AdminHandler) ToggleShortlink(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    400,
			Message: "Invalid ID",
		})
		return
	}

	shortlink, err := h.shortlinkService.Toggle(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Code:    500,
			Message: "Failed to toggle shortlink: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    0,
		Message: "success",
		Data:    shortlink,
	})
}

func (h *AdminHandler) GetStats(c *gin.Context) {
	stats, err := h.statsService.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Code:    500,
			Message: "Failed to get stats: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    0,
		Message: "success",
		Data:    stats,
	})
}

func (h *AdminHandler) ExportCSV(c *gin.Context) {
	code := c.Query("code")

	var shortlinkID int64
	if code != "" {
		shortlink, err := h.shortlinkService.GetByCode(code)
		if err != nil {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Code:    404,
				Message: "Shortlink not found",
			})
			return
		}
		shortlinkID = shortlink.ID
	}

	var rows *sql.Rows
	var err error

	if shortlinkID > 0 {
		rows, err = db.GetDB().Query(`
			SELECT ip_address, country, city, latitude, longitude,
				geo_precision, geo_status, os_type, os_version,
				browser_type, browser_version, device_type, visit_duration, visit_time
			FROM visits WHERE shortlink_id = ? ORDER BY visit_time DESC`,
			shortlinkID)
	} else {
		rows, err = db.GetDB().Query(`
			SELECT ip_address, country, city, latitude, longitude,
				geo_precision, geo_status, os_type, os_version,
				browser_type, browser_version, device_type, visit_duration, visit_time
			FROM visits ORDER BY visit_time DESC`)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Code:    500,
			Message: "Failed to query visits",
		})
		return
	}
	defer rows.Close()

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=visits_%s.csv", code))

	writer := csv.NewWriter(c.Writer)
	writer.Write([]string{
		"IP Address", "Country", "City", "Latitude", "Longitude",
		"Geo Precision", "Geo Status", "OS Type", "OS Version",
		"Browser Type", "Browser Version", "Device Type", "Visit Duration", "Visit Time",
	})

	for rows.Next() {
		var ip, country, city, geoPrecision, geoStatus string
		var lat, lng sql.NullFloat64
		var osType, osVersion, browserType, browserVersion, deviceType string
		var duration sql.NullInt64
		var visitTime []byte

		rows.Scan(&ip, &country, &city, &lat, &lng,
			&geoPrecision, &geoStatus, &osType, &osVersion,
			&browserType, &browserVersion, &deviceType, &duration, &visitTime)

		latStr := ""
		if lat.Valid {
			latStr = fmt.Sprintf("%f", lat.Float64)
		}
		lngStr := ""
		if lng.Valid {
			lngStr = fmt.Sprintf("%f", lng.Float64)
		}
		durStr := ""
		if duration.Valid {
			durStr = fmt.Sprintf("%d", duration.Int64)
		}

		writer.Write([]string{
			ip, country, city, latStr, lngStr,
			geoPrecision, geoStatus, osType, osVersion,
			browserType, browserVersion, deviceType, durStr, string(visitTime),
		})
	}

	writer.Flush()
}

func (h *AdminHandler) AdminPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin.html", nil)
}

func (h *AdminHandler) StatsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "stats.html", nil)
}
