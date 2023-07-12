package api_v1

import (
	"gpt_presets_backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	r.GET("/health-check", handlers.HealthCheck)
}
