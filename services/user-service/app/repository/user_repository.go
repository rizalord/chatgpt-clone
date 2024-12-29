package repository

import (
	"user-service/app/model/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
    Find(tx *gorm.DB, user *entity.User, id int) error
    FindByEmail(tx *gorm.DB, user *entity.User, email string) error
    Create(tx *gorm.DB, user *entity.User) error
    Update(tx *gorm.DB, user *entity.User) error
}

type UserRepositoryImpl struct {}

func NewUserRepositoryImpl() *UserRepositoryImpl {
    return &UserRepositoryImpl{}
}

func (r *UserRepositoryImpl) Find(tx *gorm.DB, user *entity.User, id int) error {
    return tx.Model(&entity.User{}).Find(&user, id).Error
}

func (r *UserRepositoryImpl) FindByEmail(tx *gorm.DB, user *entity.User, email string) error {
    return tx.Model(&entity.User{}).Where("email = ?", email).First(&user).Error
}

func (r *UserRepositoryImpl) Create(tx *gorm.DB, user *entity.User) error {
    return tx.Create(user).Error
}

func (r *UserRepositoryImpl) Update(tx *gorm.DB, user *entity.User) error {
    return tx.Save(user).Error
}
