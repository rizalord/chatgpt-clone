package repository

import (
	"auth-service/app/model/entity"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MessageRepositoryMock struct {
    mock.Mock
}

func (m *MessageRepositoryMock) FindAllByChatID(tx *gorm.DB, messages *[]entity.Message, chatID int) error {
    args := m.Called(tx, messages, chatID)
    return args.Error(0)
}

func (m *MessageRepositoryMock) Create(tx *gorm.DB, message *entity.Message) error {
    args := m.Called(tx, message)
    return args.Error(0)
}