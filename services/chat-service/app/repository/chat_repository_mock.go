package repository

import (
	"auth-service/app/model/entity"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type ChatRepositoryMock struct {
    mock.Mock
}

func (m *ChatRepositoryMock) FindByID(tx *gorm.DB, chat *entity.Chat, id int) error {
    args := m.Called(tx, chat, id)
    return args.Error(0)
}

func (m *ChatRepositoryMock) FindAllByUserID(tx *gorm.DB, chats *[]entity.Chat, userID int) error {
    args := m.Called(tx, chats, userID)
    return args.Error(0)
}

func (m *ChatRepositoryMock) Create(tx *gorm.DB, chat *entity.Chat) error {
    args := m.Called(tx, chat)
    return args.Error(0)
}

func (m *ChatRepositoryMock) Update(tx *gorm.DB, chat *entity.Chat) error {
    args := m.Called(tx, chat)
    return args.Error(0)
}