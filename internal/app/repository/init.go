package repository

import "gorm.io/gorm"

type Repos struct {
	User  UserRepository
	Chat  ChatRepository
	Token TokenRepository
}

func InitRepos(db *gorm.DB) *Repos {
	return &Repos{
		User:  NewGormUserRepository(db),
		Chat:  NewGormChatRepository(db),
		Token: NewGormTokenRepository(db),
	}
}
