package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Setup(dsn string) (db *gorm.DB, err error) {
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return
}
