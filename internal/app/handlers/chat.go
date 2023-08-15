package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"gpt_presets_backend/internal/app/models"
	"gpt_presets_backend/internal/app/repository"
	"gpt_presets_backend/internal/pkg/openai"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	OpenAIClient   openai.Client
	ChatRepository repository.ChatRepository
}

func NewChatHandler(r *repository.ChatRepository, c *openai.Client) *ChatHandler {
	return &ChatHandler{
		ChatRepository: *r,
		OpenAIClient:   *c,
	}
}

func (h *ChatHandler) CreateChat(c *gin.Context) {
	body := models.ChatContent{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Bad request",
			Error:   err.Error(),
		})
		return
	}

	res, err := h.OpenAIClient.TextCompletion(&[]openai.Message{body.ToOpenAIMessage()})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "OpenAI error",
			Error:   err.Error(),
		})
		return
	}

	u, exists := c.Get("user")
	user := u.(*models.User)
	if !exists || user.ID == 0 {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Message: "Could not retrieve user info",
		})
	}

	record, err := h.ChatRepository.CreateChat(&models.Chat{
		UserID: user.ID,
		ChatContents: []models.ChatContent{
			body,
			{Message: res.Choices[len(res.Choices)-1].Message},
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Error upon chat creation",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.DataResponse{
		Message: "OK",
		Data:    record,
	})
}

func (h *ChatHandler) GetChat(c *gin.Context) {
	id := c.Param("chat_id")

	u64, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Bad request",
			Error:   err.Error(),
		})
		return
	}

	record, err := h.ChatRepository.GetByID(uint(u64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Error while reading record",
			Error:   err.Error(),
		})
		return
	}

	u, exists := c.Get("user")
	user := u.(*models.User)
	if !exists || user.ID == 0 {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Message: "Could not retrieve user info",
		})
	}

	if record.UserID != user.ID {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Message: "Error while reading record",
			Error:   errors.New("record not found").Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{
		Message: "OK",
		Data:    record,
	})
}

// TODO allow to change temperature, max tokens, and clear history in request.
func (h *ChatHandler) PostInChat(c *gin.Context) {
	id := c.Param("chat_id")

	u64, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Bad request",
			Error:   err.Error(),
		})
		return
	}
	chatID := uint(u64)

	u, exists := c.Get("user")
	user := u.(*models.User)
	if !exists || user.ID == 0 {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Message: "Could not retrieve user info",
		})
	}

	record, err := h.ChatRepository.GetByID(chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Error while reading record",
			Error:   err.Error(),
		})
		return
	}

	body := models.ChatContent{
		ChatID:    chatID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Bad request",
			Error:   err.Error(),
		})
		return
	}

	if record.UserID != user.ID {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Message: "Error while reading record",
			Error:   errors.New("record not found").Error(),
		})
		return
	}

	history := record.GetOpenAIMessages()
	history = append(history, body.ToOpenAIMessage())

	messages, err := h.OpenAIClient.BuildHistory(&history)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "OpenAI error",
			Error:   err.Error(),
		})
		return
	}

	res, err := h.OpenAIClient.TextCompletion(messages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "OpenAI error",
			Error:   err.Error(),
		})
		return
	}

	err = h.ChatRepository.WriteInChat(chatID, &[]models.ChatContent{
		body,
		{
			ChatID:  chatID,
			Message: res.Choices[0].Message,
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "OpenAI error",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"body":  body,
		"res":   res,
		"input": messages,
	})
}
