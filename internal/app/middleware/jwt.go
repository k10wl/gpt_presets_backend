package middleware

import (
	"net/http"
	"strings"

	"gpt_presets_backend/internal/app/constants"
	"gpt_presets_backend/internal/app/handlers"
	"gpt_presets_backend/internal/app/utils/token"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware(r *handlers.UserHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer "

		bearer := c.GetHeader("Authorization")

		if !strings.Contains(bearer, BEARER_SCHEMA) {
			c.JSON(http.StatusUnauthorized, handlers.Response{
				Message: constants.Unauthorized,
			})
			c.Abort()
			return
		}

		authToken := bearer[len(BEARER_SCHEMA):]
		payload, err := token.ParseUserToken(authToken, token.AUTH_SIGNATURE)
		if err != nil {
			c.JSON(http.StatusUnauthorized, handlers.Response{
				Message: constants.Unauthorized,
			})
			c.Abort()
			return
		}

		user, err := r.Repos.User.FindUserTokensByID(payload.PublicUser.ID)

		if err != nil || user.Tokens.AuthToken != authToken {
			c.JSON(http.StatusUnauthorized, handlers.Response{
				Message: constants.Unauthorized,
			})
			c.Abort()
			return
		}

		c.Set("user", user)

		c.Next()
	}
}
