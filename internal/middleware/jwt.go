package middleware

import (
	"gpt_presets_backend/internal/handlers"
	"gpt_presets_backend/internal/models"
	"gpt_presets_backend/internal/utils/token"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware(r *handlers.UserHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer "

		bearer := c.GetHeader("Authorization")

		if !strings.Contains(bearer, BEARER_SCHEMA) {
			c.JSON(http.StatusUnauthorized, models.MessageResponse{
				Message: "Unauthorized",
			})
			c.Abort()
			return
		}

		payload, err := token.ParseUserToken(bearer[len(BEARER_SCHEMA):], os.Getenv("JWT_USER_SIGNATURE"))

		if err != nil {
			c.JSON(http.StatusUnauthorized, models.MessageResponse{
				Message: "Unauthorized",
			})
			c.Abort()
			return
		}

		user, err := r.UserRepository.FindUserByID(payload.PublicUser.ID)

		if err != nil {
			c.JSON(http.StatusUnauthorized, models.MessageResponse{
				Message: "Unauthorized",
			})
			c.Abort()
			return
		}

		c.Set("user", user)

		c.Next()
	}
}
