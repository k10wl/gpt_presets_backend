package repository

import (
	"gpt_presets_backend/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(u *models.User) (*models.User, error)
	FindUserByID(id uint) (*models.User, error)
	FindUserByName(name string) (*models.User, error)
	StoreAccessTokens(t models.Tokens) error
	FindUserTokensByID(id uint) (*models.User, error)
}

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) UserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) CreateUser(u *models.User) (*models.User, error) {
	err := r.db.Create(u).Error
	return u, err
}

func (r *GormUserRepository) FindUserByName(name string) (*models.User, error) {
	var user *models.User
	err := r.db.Where(models.User{Name: name}).First(&user).Error
	return user, err
}

func (r *GormUserRepository) FindUserByID(id uint) (*models.User, error) {
	var user *models.User
	err := r.db.Where(models.User{ID: id}).First(&user).Error
	return user, err
}

func (r *GormUserRepository) FindUserTokensByID(id uint) (*models.User, error) {
	var user *models.User
	err := r.db.Where(models.User{ID: id}).Preload("Tokens").First(&user).Error
	return user, err
}

func (r *GormUserRepository) StoreAccessTokens(t models.Tokens) error {
	res := r.db.Model(&models.Tokens{}).Where("user_id = ?", t.UserID).Updates(t)

	if res.RowsAffected == 0 {
		res = r.db.Model(&models.Tokens{}).Save(t)
	}

	return res.Error
}
