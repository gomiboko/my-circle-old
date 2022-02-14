package testutils

import (
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/gomiboko/my-circle/db"
)

// テスト用DB登録済みのデータ
const (
	User1Email        = "user1@example.com"
	User1Password     = "password"
	User1PasswordHash = "$2a$10$5zIf9lXlK6F7eaMB38uRSes9ecydTeW/xDA53zADvQjrmxA/Q/BsG"
	User1Name         = "user1"
)

const (
	UnregisteredEmail = "not-exist@example.com"
	ValidEmail        = "foo@example.com"
	ValidPassword     = "password123"
	FullWidthSpace    = "　"
	HalfWidthSpace    = " "
)

const (
	EmailMaxLength    = 254
	PasswordMinLength = 8
	PasswordMaxLength = 128
)

var (
	ErrTest           = errors.New("test error")
	ErrDuplicateEntry = &mysql.MySQLError{
		Number:  uint16(db.ErrDuplicateEntry),
		Message: "db.ErrDuplicateEntry test message",
	}
)
