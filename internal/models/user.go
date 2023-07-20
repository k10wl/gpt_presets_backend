package models

import (
	"gpt_presets_backend/internal/utils/password"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	ID        uint       `json:"id"`
	Name      string     `gorm:"type:varchar(255);uniqueIndex" json:"name,omitempty" binding:"required"`
	Password  string     `json:"password,omitempty" binding:"required"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`

	Tokens Tokens
}

type Tokens struct {
	gorm.Model
	AuthToken    string
	RefreshToken string
	UserID       uint
}

type PublicUser struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type UserResponse struct {
	User PublicUser `json:"user"`
}

type LoginUser struct {
	Name     string `json:"name,omitempty" binding:"required"`
	Password string `json:"password,omitempty" binding:"required"`
}

func (u *User) BeforeCreate(db *gorm.DB) (err error) {
	password, err := password.HashPassword(u.Password)

	if err != nil {
		return err
	}

	u.Password = password

	return
}

func (u *User) PublicInfo() PublicUser {
	return PublicUser{
		ID:   u.ID,
		Name: u.Name,
	}

}
