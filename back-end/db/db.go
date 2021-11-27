package db

import (
	"fmt"
	"net/url"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func Init() error {
	dbSvcNm := os.Getenv("DB_SERVICE_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbTz := os.Getenv("DB_TIME_ZONE")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb3&parseTime=True&loc=%s",
		dbUser,
		dbPass,
		dbSvcNm,
		dbName,
		url.QueryEscape(dbTz))
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	return err
}

func GetDB() *gorm.DB {
	return db
}
