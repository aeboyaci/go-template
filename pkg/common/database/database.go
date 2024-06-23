package database

import (
	"go-template/pkg/common/env"
	"go-template/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
)

var (
	db   *gorm.DB
	once sync.Once
)

func Initialize() error {
	var err error
	db, err = gorm.Open(postgres.Open(env.DB_URL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return err
	}

	return nil
}

func GetInstance() *gorm.DB {
	if db == nil {
		once.Do(func() {
			if err := Initialize(); err != nil {
				panic("database is not initialized")
			}
		})
	}

	return db
}
