package testutils

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/gomiboko/my-circle/controllers/mocks"
	"github.com/gomiboko/my-circle/db"
	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ApiErrorReponse struct {
	Message string
}

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

func CreateRequestBodyStr(obj interface{}) (string, error) {
	if j, err := json.Marshal(obj); err != nil {
		return "", err
	} else {
		return string(j), nil
	}
}

func SetSessionMockToGin(c *gin.Context, sessMock *mocks.SessionMock) {
	// sessions.Sessions(string, sessions.Store) と同様の処理を実行
	c.Set(sessions.DefaultKey, sessMock)
}

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
