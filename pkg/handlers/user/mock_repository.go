package user

import (
	"go-template/pkg/models"
	"gorm.io/gorm"
)

type MockRepository struct {
	MFindUserById    func(tx *gorm.DB, userId uint) (models.User, error)
	MFindUserByEmail func(tx *gorm.DB, email string) (models.User, error)
	MSaveUser        func(tx *gorm.DB, account models.User) error
}

func NewMockRepository() *MockRepository {
	return &MockRepository{}
}

func (m *MockRepository) FindUserById(tx *gorm.DB, userId uint) (models.User, error) {
	return m.MFindUserById(tx, userId)
}

func (m *MockRepository) FindUserByEmail(tx *gorm.DB, email string) (models.User, error) {
	return m.MFindUserByEmail(tx, email)
}

func (m *MockRepository) SaveUser(tx *gorm.DB, account models.User) error {
	return m.MSaveUser(tx, account)
}
