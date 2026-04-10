package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"location-tracking-shortlink/db"
	"location-tracking-shortlink/models"
	"location-tracking-shortlink/services"
	"location-tracking-shortlink/utils"
)

type VisitHandler struct {
	shortlinkService *services.ShortlinkService
	geoService       *services.GeoService
	deviceService    *services.DeviceService
}

func NewVisitHandler() *VisitHandler {
	return &VisitHandler{
		shortlinkService: services.NewShortlinkService(),
		geoService:       services.NewGeoService(),
		deviceService:    services.NewDeviceService(),
	}
}

func (h *VisitHandler) Submit(c *gin.Context) {
	var req models.VisitSubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    400,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	shortlink, err := h.shortlinkService.GetByCode(req.Code)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Code:    404,
			Message: "Shortlink not found",
		})
		return
	}

	ip := utils.GetClientIP(c.Request)
	userAgent := c.Request.UserAgent()
	referer := c.Request.Referer()

	deviceInfo := h.deviceService.Parse(userAgent)

	geoLocation, err := h.geoService.GetLocationByIP(ip)
	if err != nil {
		geoLocation = &models.GeoLocation{}
	}

	lat := req.Latitude
	lng := req.Longitude
	geoPrecision := req.GeoPrecision
	geoStatus := req.GeoStatus

	if lat != nil && lng != nil {
		geoPrecision = "gps"
		geoLocation = &models.GeoLocation{
			Latitude:  *lat,
			Longitude: *lng,
			Precision: "gps",
		}
	}

	result, err := db.GetDB().Exec(`
		INSERT INTO visits (
			shortlink_id, ip_address, country, city, latitude, longitude,
			geo_precision, geo_status, user_agent, os_type, os_version,
			browser_type, browser_version, device_type, visit_duration, visit_time, referer
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		shortlink.ID, ip, geoLocation.Country, geoLocation.City,
		geoLocation.Latitude, geoLocation.Longitude,
		geoPrecision, geoStatus, userAgent,
		deviceInfo.OS, deviceInfo.OSVersion, deviceInfo.Browser, deviceInfo.BrowserVer,
		deviceInfo.DeviceType, req.VisitDuration, time.Now(), referer,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Code:    500,
			Message: "Failed to record visit: " + err.Error(),
		})
		return
	}

	visitID, _ := result.LastInsertId()
	h.shortlinkService.IncrementVisits(shortlink.ID, visitID)

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    0,
		Message: "success",
	})
}

func (h *VisitHandler) UpdateDuration(c *gin.Context) {
	var req struct {
		Code     string `json:"code" binding:"required"`
		Duration int64  `json:"duration"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    400,
			Message: "Invalid request",
		})
		return
	}

	shortlink, err := h.shortlinkService.GetByCode(req.Code)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Code:    404,
			Message: "Shortlink not found",
		})
		return
	}

	_, err = db.GetDB().Exec(
		"UPDATE visits SET visit_duration = ?, exit_time = ? WHERE shortlink_id = ? AND exit_time IS NULL ORDER BY visit_time DESC LIMIT 1",
		req.Duration, time.Now(), shortlink.ID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Code:    500,
			Message: "Failed to update duration",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    0,
		Message: "success",
	})
}

func (h *VisitHandler) List(c *gin.Context) {
	code := c.Query("code")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 50
	}

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

	offset := (page - 1) * pageSize

	var total int
	var query string
	var args []interface{}

	if shortlinkID > 0 {
		query = "SELECT COUNT(*) FROM visits WHERE shortlink_id = ?"
		args = append(args, shortlinkID)
	} else {
		query = "SELECT COUNT(*) FROM visits"
	}

	if err := db.GetDB().QueryRow(query, args...).Scan(&total); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Code:    500,
			Message: "Failed to count visits",
		})
		return
	}

	var rows *sql.Rows
	var err error

	if shortlinkID > 0 {
		rows, err = db.GetDB().Query(`
			SELECT id, ip_address, country, city, latitude, longitude,
				geo_precision, geo_status, os_type, os_version,
				browser_type, browser_version, device_type, visit_duration, visit_time
			FROM visits WHERE shortlink_id = ? ORDER BY visit_time DESC LIMIT ? OFFSET ?`,
			shortlinkID, pageSize, offset)
	} else {
		rows, err = db.GetDB().Query(`
			SELECT id, ip_address, country, city, latitude, longitude,
				geo_precision, geo_status, os_type, os_version,
				browser_type, browser_version, device_type, visit_duration, visit_time
			FROM visits ORDER BY visit_time DESC LIMIT ? OFFSET ?`,
			pageSize, offset)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Code:    500,
			Message: "Failed to query visits",
		})
		return
	}
	defer rows.Close()

	var visits []models.VisitResponse
	for rows.Next() {
		var v models.VisitResponse
		var lat, lng sql.NullFloat64
		var country, city, geoPrecision, geoStatus sql.NullString
		var osType, osVersion, browserType, browserVersion, deviceType sql.NullString
		var visitDuration sql.NullInt64
		var visitTime time.Time

		if err := rows.Scan(&v.ID, &v.IPAddress, &country, &city, &lat, &lng,
			&geoPrecision, &geoStatus, &osType, &osVersion,
			&browserType, &browserVersion, &deviceType, &visitDuration, &visitTime); err != nil {
			continue
		}

		v.Country = country.String
		v.City = city.String
		v.GeoPrecision = geoPrecision.String
		v.GeoStatus = geoStatus.String
		v.OSType = osType.String
		v.OSVersion = osVersion.String
		v.BrowserType = browserType.String
		v.BrowserVersion = browserVersion.String
		v.DeviceType = deviceType.String
		if visitDuration.Valid {
			v.VisitDuration = visitDuration.Int64
		}
		if lat.Valid {
			v.Latitude = &lat.Float64
		}
		if lng.Valid {
			v.Longitude = &lng.Float64
		}
		v.VisitTime = visitTime.Format(time.RFC3339)

		visits = append(visits, v)
	}

	if visits == nil {
		visits = []models.VisitResponse{}
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    0,
		Message: "success",
		Data: models.VisitListResponse{
			Items:    visits,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}
