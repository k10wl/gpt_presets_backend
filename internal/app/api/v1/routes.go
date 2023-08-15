package api_v1

import (
	"os"

	"gpt_presets_backend/internal/app/handlers"
	"gpt_presets_backend/internal/app/middleware"
	"gpt_presets_backend/internal/app/repository"
	"gpt_presets_backend/internal/pkg/openai"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routes(r *gin.Engine, db *gorm.DB) {
	openaiClient := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	repos := repository.InitRepos(db)

	// TODO: add handlers container similar to repos
	userHandler := handlers.NewUserHandler(repos)
	chatHandler := handlers.NewChatHandler(repos, openaiClient)

	r.GET("/api_v1/health-check", handlers.HealthCheck)

	r.POST("/api_v1/auth/sign-up", userHandler.SignUp)
	r.POST("/api_v1/auth/sign-in", userHandler.SignIn)
	r.POST(
		"/api_v1/auth/sign-out",
		middleware.JwtAuthMiddleware(userHandler),
		middleware.UserMiddleware,
		userHandler.SignOut,
	)
	r.GET("/api_v1/auth/refresh-tokens", userHandler.RefreshTokens)
	r.GET(
		"/api_v1/auth/wall",
		middleware.JwtAuthMiddleware(userHandler),
		middleware.UserMiddleware,
		handlers.HealthCheck,
	)

	r.POST(
		"/api_v1/chat",
		middleware.JwtAuthMiddleware(userHandler),
		middleware.UserMiddleware,
		chatHandler.CreateChat,
	)
	r.GET(
		"/api_v1/chat/:chat_id",
		middleware.JwtAuthMiddleware(userHandler),
		middleware.UserMiddleware,
		chatHandler.GetChat,
	)
	r.POST(
		"/api_v1/chat/:chat_id",
		middleware.JwtAuthMiddleware(userHandler),
		middleware.UserMiddleware,
		chatHandler.PostInChat,
	)
}
