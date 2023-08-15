package middleware

import (
	"net/http"

	"gpt_presets_backend/internal/app/constants"
	"gpt_presets_backend/internal/app/handlers"
	"gpt_presets_backend/internal/app/models"

	"github.com/gin-gonic/gin"
)

func UserMiddleware(c *gin.Context) {
	user, exists := c.Get("user")

	if !exists {
		c.JSON(http.StatusUnauthorized, handlers.Response{
			Message: constants.Unauthorized,
		})
		c.Abort()
		return
	}

	u := user.(*models.User)

	if u.ID == 0 {
		c.JSON(http.StatusUnauthorized, handlers.Response{
			Message: constants.Unauthorized,
		})
		c.Abort()
		return
	}

	c.Next()
}
