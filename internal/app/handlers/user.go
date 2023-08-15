package handlers

import (
	"net/http"

	"gpt_presets_backend/internal/app/constants"
	"gpt_presets_backend/internal/app/models"
	"gpt_presets_backend/internal/app/repository"
	"gpt_presets_backend/internal/app/utils"
	"gpt_presets_backend/internal/app/utils/password"
	"gpt_presets_backend/internal/app/utils/token"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Repos  *repository.Repos
	Logger *utils.Logger
}

type UserWithTokens struct {
	AuthToken    string            `json:"auth_token"`
	RefreshToken string            `json:"refresh_token"`
	PublicUser   models.PublicUser `json:"user"`
}

func NewUserHandler(r *repository.Repos) *UserHandler {
	return &UserHandler{
		Repos:  r,
		Logger: utils.NewLogger("user handler"),
	}
}

// TODO: `[ ] remember me` checkbox for infinite refresh token duration
func (h *UserHandler) SignUp(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Message: constants.BadInput,
		})
		return
	}

	record, err := h.Repos.User.CreateUser(&user)
	if err != nil {
		h.Logger.Log("sign up create user database error", err)
		c.JSON(http.StatusInternalServerError, Response{
			Message: constants.UnableToProcess,
		})
		return
	}

	tokens, err := token.CreateAccessTokens(record.PublicInfo())
	if err != nil {
		h.Logger.Log("sign up create access tokens database error", err)
		c.JSON(http.StatusInternalServerError, Response{
			Message: constants.UnableToProcess,
		})
		return
	}

	err = h.Repos.Token.CreateToken(&tokens)
	if err != nil {
		h.Logger.Log("sign up store access tokens database error", err)
		c.JSON(http.StatusInternalServerError, Response{
			Message: constants.UnableToProcess,
		})
		return
	}

	c.JSON(http.StatusCreated,
		Response{
			Message: "Created user",
			Data: UserWithTokens{
				AuthToken:    tokens.AuthToken,
				RefreshToken: tokens.RefreshToken,
				PublicUser:   record.PublicInfo(),
			},
		},
	)
}

// TODO: [ ] remember me checkbox for infinite refresh token duration
func (h *UserHandler) SignIn(c *gin.Context) {
	var user models.LoginUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Message: constants.BadInput,
		})
		return
	}

	record, err := h.Repos.User.FindUserByName(user.Name)
	if err != nil {
		c.JSON(http.StatusNotFound, Response{
			Message: constants.WrongCredentials,
		})
		return
	}

	if err := password.ComparePassword(record.Password, user.Password); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Message: constants.WrongCredentials,
		})
		return
	}

	tokens, err := token.CreateAccessTokens(record.PublicInfo())
	if err != nil {
		h.Logger.Log("error upon tokens creation", err)
		c.JSON(http.StatusInternalServerError, Response{
			Message: constants.UnableToProcess,
		})
		return
	}

	h.Repos.Token.UpdateToken(
		&models.Tokens{
			UserID:       record.ID,
			AuthToken:    tokens.AuthToken,
			RefreshToken: tokens.RefreshToken,
		},
	)

	c.JSON(http.StatusCreated, Response{
		Message: "Authorized",
		Data: UserWithTokens{
			AuthToken:    tokens.AuthToken,
			RefreshToken: tokens.RefreshToken,
			PublicUser:   record.PublicInfo(),
		},
	})
}

func (h *UserHandler) SignOut(c *gin.Context) {
	u, exists := c.Get("user")
	user := u.(*models.User)
	if !exists || user.ID == 0 {
		c.JSON(http.StatusUnauthorized, Response{
			Message: "Could not retrieve user info",
		})
	}

	err := h.Repos.Token.InvalidateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: "Failed",
			Error:   err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *UserHandler) RefreshTokens(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token,omitempty"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusUnauthorized, Response{
			Message: "Bad token",
		})
		return
	}

	claim, err := token.ParseUserToken(body.RefreshToken, token.REFRESH_SIGNATURE)
	if err != nil {
		c.JSON(http.StatusUnauthorized, Response{
			Message: "Bad token",
		})
		return
	}

	userWithTokens, err := h.Repos.User.FindUserTokensByID(claim.PublicUser.ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, Response{
			Message: "Bad token",
		})
		return
	}

	if userWithTokens.Tokens.RefreshToken != body.RefreshToken {
		c.JSON(http.StatusUnauthorized, Response{
			Message: "Bad token",
		})
		return
	}

	tokens, err := token.CreateAccessTokens(claim.PublicUser)
	if err != nil {
		c.JSON(http.StatusUnauthorized, Response{
			Message: "Bad token",
		})
		return
	}

	h.Repos.Token.UpdateToken(&tokens)

	c.JSON(http.StatusCreated, UserWithTokens{
		AuthToken:    tokens.AuthToken,
		RefreshToken: tokens.RefreshToken,
		PublicUser:   userWithTokens.PublicInfo(),
	})
}
