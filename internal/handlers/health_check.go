package handlers

import (
	"gpt_presets_backend/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	c.IndentedJSON(http.StatusOK,
		models.BaseResponse{
			Message: "OK",
		},
	)
}
