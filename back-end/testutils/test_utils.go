package testutils

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/gomiboko/my-circle/consts"
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

// GETリクエスト用のGinコンテキストを生成する
func CreateGetContext(path string) (*httptest.ResponseRecorder, *gin.Context) {
	return createRequestContext(http.MethodGet, path, nil)
}

// POSTリクエスト用のGinコンテキストを生成する
func CreatePostContext(path, reqBody string) (*httptest.ResponseRecorder, *gin.Context) {
	return createRequestContext(http.MethodPost, path, strings.NewReader(reqBody))
}

func CreateFileHeader(filePath string) (*multipart.FileHeader, error) {
	const key = "tmp"

	// マルチパートリクエストを作成
	buf := new(bytes.Buffer)
	mw := multipart.NewWriter(buf)

	splittedPath := strings.Split(filePath, "/")
	fileName := splittedPath[len(splittedPath)-1]

	// リクエストボディに送信するファイルのヘッダを追加
	part, err := mw.CreateFormFile(key, fileName)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// リクエストボディにファイルデータを書き込み
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	mw.Close()

	r := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(r)
	// 指定のファイルをFileHeader型で取得するだけなのでURLはダミー
	c.Request, _ = http.NewRequest(http.MethodPost, "/dummy/url", buf)
	c.Request.Header.Set("Content-Type", mw.FormDataContentType())

	return c.FormFile(key)
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

func GetAwsConfig() (aws.Config, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   consts.AwsPartitionID,
			URL:           "http://test-localstack:4566",
			SigningRegion: "ap-northeast-1",
		}, nil
	})
	conf, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(os.Getenv(consts.EnvAWSRegion)),
		config.WithEndpointResolverWithOptions(customResolver),
	)

	return conf, err
}

func createRequestContext(httpMethod, path string, reqBody io.Reader) (*httptest.ResponseRecorder, *gin.Context) {
	r := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(r)
	c.Request, _ = http.NewRequest(httpMethod, path, reqBody)
	if httpMethod == http.MethodPost {
		c.Request.Header.Set("Content-Type", "application/json")
	}

	return r, c
}

func getDSN() string {
	return fmt.Sprintf("root:root@tcp(test-db:3306)/mycircle?charset=utf8mb3&parseTime=True&loc=%s",
		url.QueryEscape("Asia/Tokyo"))
}
