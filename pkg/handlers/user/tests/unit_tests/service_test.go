package unit_tests

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go-template/pkg/common/logger"
	"go-template/pkg/handlers/user"
	"go-template/pkg/handlers/user/tests"
	"go-template/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

func Test_Login_DatabaseError(t *testing.T) {
	mockRepository := tests.NewMockRepository()
	mockRepository.MFindUserByEmail = func(db *gorm.DB, email string) (models.User, error) {
		return models.User{}, gorm.ErrInvalidDB
	}

	underTest := user.NewService(nil, logger.NoOpLogger(), mockRepository)
	token, err := underTest.Login("email@test.app", "123")
	assert.NotNil(t, err)
	assert.EqualError(t, err, "database error: internal server error")
	assert.Equal(t, token, "")
}

func Test_Login_IncorrectEmail(t *testing.T) {
	mockRepository := tests.NewMockRepository()
	mockRepository.MFindUserByEmail = func(db *gorm.DB, email string) (models.User, error) {
		return models.User{}, gorm.ErrRecordNotFound
	}

	underTest := user.NewService(nil, logger.NoOpLogger(), mockRepository)
	token, err := underTest.Login("incorrect_email@test.app", "password")
	assert.NotNil(t, err)
	assert.EqualError(t, err, "incorrect email or password: not found")
	assert.Equal(t, token, "")
}

func Test_Login_IncorrectPassword(t *testing.T) {
	mockRepository := tests.NewMockRepository()
	mockRepository.MFindUserByEmail = func(db *gorm.DB, email string) (models.User, error) {
		return models.User{
			Model: gorm.Model{
				ID: 1,
			},
			Email:    "email@test.app",
			Password: "$2a$10$3Q7...",
		}, nil
	}

	underTest := user.NewService(nil, logger.NoOpLogger(), mockRepository)
	token, err := underTest.Login("email@test.app", "incorrect_password")
	assert.NotNil(t, err)
	assert.EqualError(t, err, "incorrect email or password: not found")
	assert.Equal(t, token, "")
}

func Test_Login_Success(t *testing.T) {
	mockRepository := tests.NewMockRepository()
	mockRepository.MFindUserByEmail = func(db *gorm.DB, email string) (models.User, error) {
		return models.User{
			Model: gorm.Model{
				ID: 1,
			},
			Email:    "email@test.app",
			Password: "$2a$10$jo.D0ZDxxfGhyB4hcgNdnuW9mnSytcsDf8w8vhbHTtIscNmp5OGl.",
		}, nil
	}

	underTest := user.NewService(nil, logger.NoOpLogger(), mockRepository)
	token, err := underTest.Login("email@test.app", "123")
	assert.Nil(t, err)
	assert.NotEmpty(t, token)
}

func Test_Register_DatabaseError_FindUserByEmail(t *testing.T) {
	mockRepository := tests.NewMockRepository()
	mockRepository.MFindUserByEmail = func(db *gorm.DB, email string) (models.User, error) {
		return models.User{}, gorm.ErrInvalidDB
	}

	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.Nil(t, err)

	mock.ExpectBegin()
	mock.ExpectRollback()

	underTest := user.NewService(gormDB, logger.NoOpLogger(), mockRepository)
	err = underTest.Register("email@test.app", "123")
	assert.NotNil(t, err)
	assert.EqualError(t, err, "database error: internal server error")
}

func Test_Register_DatabaseError_SaveUser(t *testing.T) {
	mockRepository := tests.NewMockRepository()
	mockRepository.MFindUserByEmail = func(db *gorm.DB, email string) (models.User, error) {
		return models.User{}, gorm.ErrRecordNotFound
	}
	mockRepository.MSaveUser = func(db *gorm.DB, user models.User) error {
		return gorm.ErrInvalidDB
	}

	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.Nil(t, err)

	mock.ExpectBegin()
	mock.ExpectRollback()

	underTest := user.NewService(gormDB, logger.NoOpLogger(), mockRepository)
	err = underTest.Register("email@test.app", "123")
	assert.NotNil(t, err)
	assert.EqualError(t, err, "database error: internal server error")
}

func Test_Register_EmailAlreadyExists(t *testing.T) {
	mockRepository := tests.NewMockRepository()
	mockRepository.MFindUserByEmail = func(db *gorm.DB, email string) (models.User, error) {
		return models.User{
			Model: gorm.Model{
				ID: 1,
			},
			Email:    "email@test.app",
			Password: "$2a$10$3Q7...",
		}, nil
	}

	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.Nil(t, err)

	mock.ExpectBegin()
	mock.ExpectRollback()

	underTest := user.NewService(gormDB, logger.NoOpLogger(), mockRepository)
	err = underTest.Register("email@test.app", "123")
	assert.NotNil(t, err)
	assert.EqualError(t, err, "email already exists: bad request")
}

func Test_Register_Success(t *testing.T) {
	mockRepository := tests.NewMockRepository()
	mockRepository.MFindUserByEmail = func(db *gorm.DB, email string) (models.User, error) {
		return models.User{}, gorm.ErrRecordNotFound
	}
	mockRepository.MSaveUser = func(db *gorm.DB, user models.User) error {
		return nil
	}

	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.Nil(t, err)

	mock.ExpectBegin()
	mock.ExpectCommit()
	mock.ExpectRollback()

	underTest := user.NewService(gormDB, logger.NoOpLogger(), mockRepository)
	err = underTest.Register("email@test.app", "123")
	assert.Nil(t, err)
}
