package repository

import (
	"auth-service/app/model/entity"

	"gorm.io/gorm"
)

type ChatRepository interface {
    FindByID(tx *gorm.DB, chat *entity.Chat, id int) error
    FindByIDAndUserID(tx *gorm.DB, chat *entity.Chat, id, userID int) error
    FindAllByUserID(tx *gorm.DB, chats *[]entity.Chat, userID int) error
    Create(tx *gorm.DB, user *entity.Chat) error
    Update(tx *gorm.DB, user *entity.Chat) error
}

type ChatRepositoryImpl struct {}

func NewChatRepositoryImpl() *ChatRepositoryImpl {
    return &ChatRepositoryImpl{}
}

func (r *ChatRepositoryImpl) FindByID(tx *gorm.DB, chat *entity.Chat, id int) error {
    return tx.Model(&entity.Chat{}).First(&chat, id).Error
}

func (r *ChatRepositoryImpl) FindByIDAndUserID(tx *gorm.DB, chat *entity.Chat, id, userID int) error {
    return tx.Model(&entity.Chat{}).
        Preload("Messages").
        Where("id = ? AND user_id = ?", id, userID).
        First(&chat).Error
}

func (r *ChatRepositoryImpl) FindAllByUserID(tx *gorm.DB, chats *[]entity.Chat, userID int) error {
    return tx.Model(&entity.Chat{}).Where("user_id = ?", userID).Find(&chats).Error
}

func (r *ChatRepositoryImpl) Create(tx *gorm.DB, chat *entity.Chat) error {
    return tx.Create(chat).Error
}

func (r *ChatRepositoryImpl) Update(tx *gorm.DB, chat *entity.Chat) error {
    return tx.Save(chat).Error
}