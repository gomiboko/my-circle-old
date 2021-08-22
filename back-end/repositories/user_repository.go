package repositories

import (
	"errors"

	"github.com/gomiboko/my-circle/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUser(email string, password string) (*models.User, error)
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (ur *userRepository) GetUser(email string, password string) (*models.User, error) {
	var user models.User

	// email(UQ)で検索
	cond := models.User{Email: email}
	result := ur.DB.Where(&cond).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	// パスワードの比較
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, nil
		}
		return nil, err
	}

	// パスワードのハッシュ値は返却しない
	user.PasswordHash = ""

	return &user, err
}
