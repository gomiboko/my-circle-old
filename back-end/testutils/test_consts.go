package testutils

import (
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/gomiboko/my-circle/db"
)

// テスト用DB登録済みのusersデータ
const (
	User1ID           uint = 1
	User1Email             = "user1@example.com"
	User1Password          = "password"
	User1PasswordHash      = "$2a$10$5zIf9lXlK6F7eaMB38uRSes9ecydTeW/xDA53zADvQjrmxA/Q/BsG"
	User1Name              = "user1"
	User1RowVersion   uint = 2

	User2ID uint = 2
	User3ID uint = 3
)

var User1CreatedAt = time.Date(2021, 8, 24, 12, 34, 56, 0, locale())
var User1UpdatedAt = time.Date(2021, 8, 25, 23, 45, 01, 0, locale())

// テスト用DB登録済みのcirclesデータ
const (
	Circle1ID      uint = 1
	Circle1Name         = "Circle1"
	Circle1IconUrl      = "http://localhost:4566/my-circle-bucket/circles/54a2e0a21c246a49c4b2f3057ea78da4a38952dbbfa450bc120bde5d99f0a7eb"
)

var Circle1CreatedAt = time.Date(2022, 3, 28, 12, 34, 56, 0, locale())
var Circle1UpdatedAt = time.Date(2022, 3, 29, 23, 45, 01, 0, locale())

// テスト用データ
const (
	UnregisteredEmail = "not-exist@example.com"
	ValidEmail        = "foo@example.com"
	ValidPassword     = "password"
	ValidUserName     = "username"
	ValidUrl          = "https://example.com"
	FullWidthSpace    = "　"
	HalfWidthSpace    = " "
	HalfWidthSymbol   = "`~!@#$%^&*()-_=+[]{}\\|;:'\",./<>?"
	FullWidthA        = "Ａ"
	UnexpectedErrMsg  = "予期せぬエラーが発生しました"
)

const (
	EmailMaxLength    = 254
	PasswordMinLength = 8
	PasswordMaxLength = 128
	UsernameMaxLength = 45
)

// テスト用エラーオブジェクト
var (
	ErrTest           = errors.New("test error")
	ErrDuplicateEntry = &mysql.MySQLError{
		Number:  uint16(db.ErrDuplicateEntry),
		Message: "db.ErrDuplicateEntry test message",
	}
)

func locale() *time.Location {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	return jst
}
