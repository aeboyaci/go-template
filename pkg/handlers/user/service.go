package user

import (
	"github.com/pkg/errors"
	"go-template/pkg/common/apierrors"
	"go-template/pkg/common/utils"
	"go-template/pkg/models"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type Service interface {
	Login(email, password string) (string, error)
	Register(email, password string) error
}

type serviceImpl struct {
	db         *gorm.DB
	logger     *zap.Logger
	repository Repository
}

func NewService(db *gorm.DB, logger *zap.Logger, repository Repository) Service {
	return &serviceImpl{
		db:         db,
		logger:     logger,
		repository: repository,
	}
}

func (s *serviceImpl) Login(email, password string) (string, error) {
	user, err := s.repository.FindUserByEmail(s.db, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.Wrap(apierrors.ErrorNotFound, "incorrect email or password")
		}

		s.logger.Error("failed to find user by email", zap.String("location", "Login"), zap.String("email", email), zap.Error(err))
		return "", errors.Wrap(apierrors.ErrorInternalServer, "database error")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.Wrap(apierrors.ErrorNotFound, "incorrect email or password")
	}

	token := utils.SignToken(user.ID, time.Now().Add(time.Hour*24))
	s.logger.Debug("user logged in", zap.String("location", "Login"), zap.String("email", user.Email))
	return token, nil
}

func (s *serviceImpl) Register(email, password string) error {
	tx := s.db.Begin()
	defer tx.Rollback()

	user, err := s.repository.FindUserByEmail(tx, email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error("failed to find user by email", zap.String("location", "Register"), zap.String("email", email), zap.Error(err))
		return errors.Wrap(apierrors.ErrorInternalServer, "database error")
	}
	if user.ID != 0 {
		return errors.Wrap(apierrors.ErrorBadRequest, "email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("failed to hash password", zap.String("location", "Register"), zap.String("email", email), zap.Error(err))
		return apierrors.ErrorInternalServer
	}

	user = models.User{
		Email:    email,
		Password: string(hashedPassword),
	}
	if err := s.repository.SaveUser(tx, user); err != nil {
		s.logger.Error("failed to save user", zap.String("location", "Register"), zap.String("email", email), zap.Error(err))
		return errors.Wrap(apierrors.ErrorInternalServer, "database error")
	}

	if err := tx.Commit().Error; err != nil {
		s.logger.Error("failed to commit transaction", zap.String("location", "Register"), zap.String("email", email), zap.Error(err))
		return errors.Wrap(apierrors.ErrorInternalServer, "transaction commit error")
	}

	s.logger.Debug("user registered", zap.String("location", "Register"), zap.String("email", email))
	return nil
}
