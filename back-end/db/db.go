package db

import (
	"fmt"
	"net/url"
	"os"

	"github.com/gomiboko/my-circle/consts"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func Init() error {
	dbSvcNm := os.Getenv(consts.EnvDbServiceName)
	dbUser := os.Getenv(consts.EnvDbUser)
	dbPass := os.Getenv(consts.EnvDbPassword)
	dbName := os.Getenv(consts.EnvDbName)
	dbTz := os.Getenv(consts.EnvDbTimeZone)

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
