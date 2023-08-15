package models

import "gorm.io/gorm"

type Tokens struct {
	gorm.Model
	AuthToken    string
	RefreshToken string
	UserID       uint `gorm:index"`
}
