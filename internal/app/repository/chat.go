package repository

import (
	"errors"

	"gpt_presets_backend/internal/app/models"

	"gorm.io/gorm"
)

type GormChatRepository struct {
	db *gorm.DB
}

type ChatRepository interface {
	CreateChat(chat *models.Chat) (*models.Chat, error)
	GetByID(id uint) (*models.Chat, error)
	WriteInChat(chatID uint, content *[]models.ChatContent) error
}

func NewGormChatRepository(db *gorm.DB) ChatRepository {
	return &GormChatRepository{db: db}
}

func (r *GormChatRepository) GetByID(id uint) (*models.Chat, error) {
	var chat models.Chat

	res := r.db.Preload("ChatContents").Find(&chat, id)

	if chat.ID == 0 {
		return nil, errors.New("record not found")
	}

	return &chat, res.Error
}

func (r *GormChatRepository) CreateChat(chat *models.Chat) (*models.Chat, error) {
	res := r.db.Create(&chat)

	return chat, res.Error
}

func (r *GormChatRepository) WriteInChat(chatID uint, content *[]models.ChatContent) error {
	return r.db.Model(&models.ChatContent{}).Create(content).Error
}
