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
	GetChatByID(chatID uint, userID uint) (*models.Chat, error)
	WriteInChat(chatID uint, content *[]models.ChatContent) error
	GetChatRequestContentsByChatID(chatID uint, userID uint) (*models.Chat, error)
	UpdateMessagesInRequestByIDs(content *[]models.ChatContent, inRequest bool) error
}

func NewGormChatRepository(db *gorm.DB) ChatRepository {
	return &GormChatRepository{db: db}
}

func (r *GormChatRepository) GetChatByID(chatID uint, userID uint) (*models.Chat, error) {
	var chat models.Chat

	res := r.db.Preload("ChatContents").Find(&chat, "id = ? AND user_id = ?", chatID, userID)

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

func (r *GormChatRepository) UpdateMessagesInRequestByIDs(content *[]models.ChatContent, inRequest bool) error {
	var ids []uint
	for _, cc := range *content {
		ids = append(ids, cc.ID)
	}

	err := r.db.Model(&models.ChatContent{}).Where("id IN ?", ids).Update("in_request", inRequest).Error

	return err
}

func (r *GormChatRepository) GetChatRequestContentsByChatID(chatID uint, userID uint) (*models.Chat, error) {
	var chat models.Chat
	err := r.db.Preload("ChatContents", "in_request = true").Find(&chat, "id = ? AND user_id = ?", chatID, userID).Error

	return &chat, err
}

func (r *GormChatRepository) RemoveFromRequest(messageID ...[]uint) error {
	r.db.Model(&models.ChatContent{}).Where("id IN ?", messageID).Updates(&models.ChatContent{InRequest: false})

	return nil
}
