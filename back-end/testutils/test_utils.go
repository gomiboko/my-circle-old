package testutils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/gomiboko/my-circle/controllers/mocks"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ApiErrorReponse struct {
	Message string
}

func CreateRequestBodyStr(obj interface{}) (string, error) {
	if j, err := json.Marshal(obj); err != nil {
		return "", err
	} else {
		return string(j), nil
	}
}

// 指定の長さのメールアドレスを生成する
func CreateEmailAddress(length int) string {
	return strings.Repeat("a", length-len("@example.com")) + "@example.com"
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
	return gorm.Open(mysql.Open(getDSN()), &gorm.Config{})
}

func getDSN() string {
	return fmt.Sprintf("root:root@tcp(test-db:3306)/mycircle?charset=utf8mb3&parseTime=True&loc=%s",
		url.QueryEscape("Asia/Tokyo"))
}
