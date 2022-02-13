package testutils

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"

	"github.com/go-sql-driver/mysql"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/gomiboko/my-circle/db"
	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	User1Email        = "user1@example.com"
	User1Password     = "password"
	User1PasswordHash = "$2a$10$5zIf9lXlK6F7eaMB38uRSes9ecydTeW/xDA53zADvQjrmxA/Q/BsG"
	User1Name         = "user1"

	UnregisteredEmail = "not-exist@example.com"

	FullWidthSpace = "　"
	HalfWidthSpace = " "
)

var ErrTest = errors.New("test error")

var (
	ErrDuplicateEntry = &mysql.MySQLError{
		Number:  uint16(db.ErrDuplicateEntry),
		Message: "db.ErrDuplicateEntry test message",
	}
)

func GetFixtures(fixturesDirPath string) (*testfixtures.Loader, error) {
	sqldb, err := sql.Open("mysql", getDSN())
	if err != nil {
		return nil, err
	}

	return testfixtures.New(
		testfixtures.Database(sqldb),
		testfixtures.Dialect("mysql"),
		testfixtures.Directory(fixturesDirPath),
		testfixtures.DangerousSkipTestDatabaseCheck(),
	)
}

func GetDB() (*gorm.DB, error) {
	return gorm.Open(gmysql.Open(getDSN()), &gorm.Config{})
}

func getDSN() string {
	return fmt.Sprintf("root:root@tcp(test-db:3306)/mycircle?charset=utf8mb3&parseTime=True&loc=%s",
		url.QueryEscape("Asia/Tokyo"))
}
