package handlers

import (
	"fmt"
	"gpt_presets_backend/internal/models"
	"gpt_presets_backend/internal/repository"
	"gpt_presets_backend/internal/utils/password"
	"gpt_presets_backend/internal/utils/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserRepository repository.UserRepository
}

type UserWithTokens struct {
	AuthToken    string            `json:"auth_token"`
	RefreshToken string            `json:"refresh_token"`
	PublicUser   models.PublicUser `json:"user"`
}

func NewUserHandler(r repository.UserRepository) *UserHandler {
	return &UserHandler{
		UserRepository: r,
	}
}

func (h *UserHandler) SignUp(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Bad input",
			Error:   err.Error(),
		})
		return
	}

	tokens, err := token.CreateAccessTokens(user.PublicInfo())

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.MessageResponse{
			Message: "Access tokens creation failed",
		})
		return
	}

	user.Tokens = models.Tokens{UserID: user.ID, AuthToken: tokens.AuthToken, RefreshToken: tokens.RefreshToken}

	record, err := h.UserRepository.CreateUser(&user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Database error",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated,
		models.DataResponse{
			Message: "Created user",
			Data: UserWithTokens{
				AuthToken:    tokens.AuthToken,
				RefreshToken: tokens.RefreshToken,
				PublicUser:   record.PublicInfo(),
			},
		},
	)
}

func (h *UserHandler) SignIn(c *gin.Context) {
	var user models.LoginUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Bad input",
			Error:   err.Error(),
		})
	}

	record, err := h.UserRepository.FindUserByName(user.Name)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Database error",
			Error:   err.Error(),
		})
		return
	}

	if err := password.ComparePassword(record.Password, user.Password); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Wrong username or password",
			Error:   err.Error(),
		})
		return
	}

	tokens, err := token.CreateAccessTokens(record.PublicInfo())

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.MessageResponse{
			Message: "Access tokens creation failed",
		})
		return
	}

	h.UserRepository.StoreAccessTokens(
		models.Tokens{
			UserID:       record.ID,
			AuthToken:    tokens.AuthToken,
			RefreshToken: tokens.RefreshToken,
		},
	)

	c.JSON(http.StatusOK, models.DataResponse{
		Message: "Authorized",
		Data: UserWithTokens{
			AuthToken:    tokens.AuthToken,
			RefreshToken: tokens.RefreshToken,
			PublicUser:   record.PublicInfo(),
		},
	})
}

func (h *UserHandler) RefreshTokens(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token,omitempty"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusUnauthorized, models.MessageResponse{
			Message: "Bad token",
		})
		return
	}

	claim, err := token.ParseUserToken(body.RefreshToken, token.REFRESH_SIGNATURE)

	if err != nil {
		c.JSON(http.StatusUnauthorized, models.MessageResponse{
			Message: "Bad token",
		})
		return
	}

	userWithTokens, err := h.UserRepository.FindUserTokensByID(claim.PublicUser.ID)

	if err != nil {
		c.JSON(http.StatusUnauthorized, models.MessageResponse{
			Message: "Bad token",
		})
		return
	}

	fmt.Println(body.RefreshToken)
	if userWithTokens.Tokens.RefreshToken != body.RefreshToken {
		c.JSON(http.StatusUnauthorized, models.MessageResponse{
			Message: "Bad token",
		})
		return
	}

	tokens, err := token.CreateAccessTokens(claim.PublicUser)

	if err != nil {
		c.JSON(http.StatusUnauthorized, models.MessageResponse{
			Message: "Bad token",
		})
		return
	}

	h.UserRepository.StoreAccessTokens(tokens)

	c.JSON(http.StatusOK, UserWithTokens{
		AuthToken:    tokens.AuthToken,
		RefreshToken: tokens.RefreshToken,
		PublicUser:   userWithTokens.PublicInfo(),
	})
}
