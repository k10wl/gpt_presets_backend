package repository

import (
	"errors"
	"fmt"
	"time"

	"gpt_presets_backend/internal/app/models"

	"gorm.io/gorm"
)

type TokenRepository interface {
	CreateToken(t *models.Tokens) error
	UpdateToken(t *models.Tokens) error
	InvalidateToken(userID uint) error
}

type GormTokenRepository struct {
	db *gorm.DB
}

func NewGormTokenRepository(db *gorm.DB) TokenRepository {
	return &GormTokenRepository{db: db}
}

func (r *GormTokenRepository) CreateToken(t *models.Tokens) error {
	res := r.db.Create(&t)

	if res.Error != nil {
		fmt.Printf("TokenRepository error: %+v", res.Error)
		return errors.New("operation failed")
	}

	if res.RowsAffected == 0 {
		fmt.Printf("TokenRepository error: no rows where created")
		return errors.New("operation failed")
	}
	return nil
}

func (r *GormTokenRepository) UpdateToken(t *models.Tokens) error {
	res := r.db.Model(&t).Where("user_id = ?", t.UserID).Updates(&t)

	if res.Error != nil {
		fmt.Printf("TokenRepository error: %+v", res.Error)
		return errors.New("operation failed")
	}

	if res.RowsAffected == 0 {
		fmt.Printf("TokenRepository error: no rows where affected")
		return errors.New("operation failed")
	}
	return nil
}

func (r *GormTokenRepository) InvalidateToken(userID uint) error {
	fmt.Printf("userID: %v\n", userID)
	res := r.db.Model(&models.Tokens{}).Where("user_id = ?", userID).Updates(map[string]interface{}{
		"auth_token":    "",
		"refresh_token": "",
		"deleted_at":    time.Now(),
	})

	if res.Error != nil {
		fmt.Printf("TokenRepository error: %+v", res.Error)
		return errors.New("operation failed")
	}

	if res.RowsAffected == 0 {
		fmt.Printf("TokenRepository error: no rows where affected")
		return errors.New("operation failed")
	}

	return nil
}
