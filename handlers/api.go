package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"location-tracking-shortlink/models"
)

type APIHandler struct{}

func NewAPIHandler() *APIHandler {
	return &APIHandler{}
}

func (h *APIHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, models.APIResponse{
		Code:    0,
		Message: "ok",
	})
}

func (h *APIHandler) NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, models.APIResponse{
		Code:    404,
		Message: "API endpoint not found",
	})
}
