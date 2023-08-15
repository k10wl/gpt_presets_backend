package models

import (
	"time"

	"gpt_presets_backend/internal/pkg/openai"

	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model `json:"-"`

	ID        uint       `json:"id,omitempty"`
	UserID    uint       `json:"user_id,omitempty"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	ChatContents []ChatContent `json:"chat_contents,omitempty"`
}

type ChatContent struct {
	gorm.Model `json:"-"`
	openai.Message

	ID        uint       `json:"id,omitempty"`
	ChatID    uint       `json:"chat_id,omitempty"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func (c *Chat) GetOpenAIMessages() []openai.Message {
	messages := []openai.Message{}

	for _, cc := range c.ChatContents {
		messages = append(messages, cc.ToOpenAIMessage())
	}

	return messages
}

func (c *ChatContent) ToOpenAIMessage() openai.Message {
	return openai.Message{
		Role:    c.Role,
		Content: c.Content,
	}
}
