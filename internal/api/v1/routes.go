package api_v1

import (
	"gpt_presets_backend/internal/handlers"
	"gpt_presets_backend/internal/middleware"
	"gpt_presets_backend/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routes(r *gin.Engine, db *gorm.DB) {
	userRepository := repository.NewGormUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepository)

	r.GET("/api_v1/health-check", handlers.HealthCheck)

	r.POST("/api_v1/auth/sign-up", userHandler.SignUp)
	r.POST("/api_v1/auth/sign-in", userHandler.SignIn)
	r.GET("/api_v1/auth/refresh-tokens", userHandler.RefreshTokens)
	r.GET("/api_v1/auth/wall", middleware.JwtAuthMiddleware(userHandler), handlers.HealthCheck)
}
