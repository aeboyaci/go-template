package user

import (
	"go-template/pkg/models"
	"gorm.io/gorm"
)

type Repository interface {
	FindUserById(tx *gorm.DB, userId uint) (models.User, error)
	FindUserByEmail(tx *gorm.DB, email string) (models.User, error)
	SaveUser(tx *gorm.DB, account models.User) error
}

type repositoryImpl struct{}

func NewRepository() Repository {
	return &repositoryImpl{}
}

func (r *repositoryImpl) FindUserById(tx *gorm.DB, userId uint) (models.User, error) {
	var result models.User
	err := tx.
		Model(&models.User{}).
		Where("id = ?", userId).
		Take(&result).Error
	return result, err
}

func (r *repositoryImpl) FindUserByEmail(tx *gorm.DB, email string) (models.User, error) {
	var result models.User
	err := tx.
		Model(&models.User{}).
		Where("email = ?", email).
		Take(&result).Error
	return result, err
}

func (r *repositoryImpl) SaveUser(tx *gorm.DB, user models.User) error {
	return tx.Save(&user).Error
}
