package handlers

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"location-tracking-shortlink/models"
	"location-tracking-shortlink/services"
)

type ShortlinkHandler struct {
	service *services.ShortlinkService
}

func NewShortlinkHandler() *ShortlinkHandler {
	return &ShortlinkHandler{
		service: services.NewShortlinkService(),
	}
}

func (h *ShortlinkHandler) Redirect(c *gin.Context) {
	code := c.Param("code")

	shortlink, err := h.service.GetByCode(code)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Code:    404,
			Message: "Shortlink not found",
		})
		return
	}

	if shortlink.IsDeleted {
		c.JSON(http.StatusGone, models.APIResponse{
			Code:    410,
			Message: "Shortlink has been deleted",
		})
		return
	}

	if shortlink.IsDisabled {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Code:    403,
			Message: "Shortlink is disabled",
		})
		return
	}

	title := shortlink.InduceTitle
	subtitle := shortlink.InduceSubtitle
	imageURL := shortlink.InduceImageURL
	template := shortlink.InduceTemplate

	for _, t := range models.TemplateLibrary {
		if t.ID == template {
			if title == "" {
				title = t.Title
			}
			if subtitle == "" {
				subtitle = t.Subtitle
			}
			break
		}
	}

	c.HTML(http.StatusOK, "redirect.html", gin.H{
		"Code":           shortlink.Code,
		"OriginalURL":    shortlink.OriginalURL,
		"InduceTitle":    title,
		"InduceSubtitle": subtitle,
		"InduceImageURL": imageURL,
		"InduceTemplate": template,
		"Icon":           getTemplateIcon(template),
	})
}

func getTemplateIcon(templateID string) string {
	for _, t := range models.TemplateLibrary {
		if t.ID == templateID {
			return t.Icon
		}
	}
	return "📢"
}

func (h *ShortlinkHandler) Create(c *gin.Context) {
	var req models.ShortlinkCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    400,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	parsedURL, err := url.Parse(req.URL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    400,
			Message: "Invalid URL format",
		})
		return
	}

	shortlink, err := h.service.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Code:    500,
			Message: "Failed to create shortlink: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    0,
		Message: "success",
		Data: models.ShortlinkResponse{
			ShortURL:    "/" + shortlink.Code,
			Code:        shortlink.Code,
			OriginalURL: shortlink.OriginalURL,
			CreatedAt:   shortlink.CreatedAt,
		},
	})
}

func (h *ShortlinkHandler) GetTemplates(c *gin.Context) {
	c.JSON(http.StatusOK, models.APIResponse{
		Code:    0,
		Message: "success",
		Data:    models.TemplateLibrary,
	})
}
