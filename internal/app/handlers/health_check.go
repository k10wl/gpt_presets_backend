package handlers

import (
	"gpt_presets_backend/internal/app/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK,
		models.MessageResponse{
			Message: "OK",
		},
	)
}
