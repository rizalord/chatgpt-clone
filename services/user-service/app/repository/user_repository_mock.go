package repository

import (
	"user-service/app/model/entity"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type UserRepositoryMock struct {
    mock.Mock
}

func (m *UserRepositoryMock) FindByEmail(tx *gorm.DB, user *entity.User, email string) error {
    args := m.Called(tx, user, email)
    return args.Error(0)
}

func (m *UserRepositoryMock) FindByID(tx *gorm.DB, user *entity.User, id uint) error {
    args := m.Called(tx, user, id)
    return args.Error(0)
}

func (m *UserRepositoryMock) FindAll(tx *gorm.DB, users *[]entity.User) error {
    args := m.Called(tx, users)
    return args.Error(0)
}

func (m *UserRepositoryMock) Create(tx *gorm.DB, user *entity.User) error {
    args := m.Called(tx, user)
    return args.Error(0)
}

func (m *UserRepositoryMock) Update(tx *gorm.DB, user *entity.User) error {
    args := m.Called(tx, user)
    return args.Error(0)
}

func (m *UserRepositoryMock) Delete(tx *gorm.DB, user *entity.User) error {
    args := m.Called(tx, user)
    return args.Error(0)
}