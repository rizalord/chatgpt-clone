package repository

import (
	"auth-service/app/model/entity"

	"gorm.io/gorm"
)

type MessageRepository interface {
    FindAllByChatID(tx *gorm.DB, messages *[]entity.Message, chatID int) error
    Create(tx *gorm.DB, message *entity.Message) error
}

type MessageRepositoryImpl struct {}

func NewMessageRepositoryImpl() *MessageRepositoryImpl {
    return &MessageRepositoryImpl{}
}

func (r *MessageRepositoryImpl) FindAllByChatID(tx *gorm.DB, messages *[]entity.Message, chatID int) error {
    return tx.Model(&entity.Message{}).Where("chat_id = ?", chatID).Find(&messages).Error
}

func (r *MessageRepositoryImpl) Create(tx *gorm.DB, message *entity.Message) error {
    return tx.Create(message).Error
}