package services

import (
	"errors"

	"github.com/gomiboko/my-circle/repositories"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Authenticate(email string, password string) (bool, error)
}

type authService struct {
	ur repositories.UserRepository
}

func NewAuthService(ur repositories.UserRepository) AuthService {
	return &authService{ur}
}

func (as *authService) Authenticate(email string, password string) (bool, error) {
	// ユーザ検索
	user, err := as.ur.GetUser(email)
	if err != nil {
		// ユーザが取得できなかった場合、認証失敗
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	// パスワード照合
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		// パスワードが異なる場合、認証失敗
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
