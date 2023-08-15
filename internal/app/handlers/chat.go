package handlers

import (
	"net/http"
	"strconv"
	"time"

	"gpt_presets_backend/internal/app/constants"
	"gpt_presets_backend/internal/app/models"
	"gpt_presets_backend/internal/app/repository"
	"gpt_presets_backend/internal/app/utils"
	"gpt_presets_backend/internal/pkg/openai"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	OpenAIClient *openai.Client
	Repos        *repository.Repos
	Logger       *utils.Logger
}

func NewChatHandler(r *repository.Repos, c *openai.Client) *ChatHandler {
	return &ChatHandler{
		Repos:        r,
		OpenAIClient: c,
		Logger:       utils.NewLogger("chat handler"),
	}
}

func (h *ChatHandler) CreateChat(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	body := models.ChatContent{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Message: constants.BadInput,
		})
		return
	}

	if overflowsTokenLimit := h.OpenAIClient.HasTokensOverflow(&body.Message); overflowsTokenLimit {
		c.JSON(http.StatusBadRequest, Response{
			Message: constants.BadInput,
			Error:   openai.TokensOverflowError,
		})
		return
	}

	res, err := h.OpenAIClient.TextCompletion(&[]openai.Message{*body.ToOpenAIMessage()})
	if err != nil {
		h.Logger.Log("create chat error from OpenAI", err)
		c.JSON(http.StatusInternalServerError, Response{
			Message: constants.UnableToProcess,
		})
		return
	}

	record, err := h.Repos.Chat.CreateChat(&models.Chat{
		UserID: user.ID,
		ChatContents: []models.ChatContent{
			body,
			{Message: res.Choices[len(res.Choices)-1].Message},
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: constants.UnableToProcess,
		})
		return
	}

	c.JSON(http.StatusCreated, Response{
		Message: "OK",
		Data:    record,
	})
}

func (h *ChatHandler) GetChat(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	chatID, err := parseChatID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Message: constants.BadInput,
		})
		return
	}

	record, err := h.Repos.Chat.GetChatByID(chatID, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Message: constants.NotFound,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Message: "OK",
		Data:    record,
	})
}

func (h *ChatHandler) PostInChat(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	chatID, err := parseChatID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Message: constants.BadInput,
		})
		return
	}

	body := models.ChatContent{
		ChatID:    chatID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Message: constants.BadInput,
		})
		return
	}

	record, err := h.Repos.Chat.GetChatRequestContentsByChatID(chatID, user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, Response{
			Message: constants.NotFound,
		})
		return
	}

	history := append(record.GetOpenAIMessages(), *body.ToOpenAIMessage())

	messages, err := h.OpenAIClient.BuildHistory(&history)
	if err != nil {
		if err.Error() == openai.TokensOverflowError {
			c.JSON(http.StatusBadRequest, Response{
				Message: constants.BadInput,
				Error:   openai.TokensOverflowError,
			})
			return
		}

		h.Logger.Log("post in chat error, cannot build history", err)
		c.JSON(http.StatusBadRequest, Response{
			Message: constants.UnableToProcess,
		})
		return
	}

	res, err := h.OpenAIClient.TextCompletion(messages)
	if err != nil {
		h.Logger.Log("post in chat error, text completion error", err)
		c.JSON(http.StatusInternalServerError, Response{
			Message: constants.UnableToProcess,
		})
		return
	}

	unusedContent := record.ChatContents[:len(*messages)]
	err = h.Repos.Chat.UpdateMessagesInRequestByIDs(&unusedContent, false)
	if err != nil {
		h.Logger.Log("post in chat error, update in_request in DB", err)
		c.JSON(http.StatusInternalServerError, Response{
			Message: constants.UnableToProcess,
		})
		return
	}

	err = h.Repos.Chat.WriteInChat(chatID, &[]models.ChatContent{
		body,
		{
			ChatID:  chatID,
			Message: res.Choices[0].Message,
		},
	})
	if err != nil {
		h.Logger.Log("post in chat error, store res in DB", err)
		c.JSON(http.StatusInternalServerError, Response{
			Message: constants.UnableToProcess,
		})
		return
	}

	c.JSON(http.StatusCreated, Response{
		Message: "OK",
		Data:    res.Choices[0].Message,
	})
}

func parseChatID(c *gin.Context) (uint, error) {
	id := c.Param("chat_id")
	u64, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(u64), nil
}
