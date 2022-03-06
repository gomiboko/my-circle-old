package services

import (
	"github.com/gomiboko/my-circle/forms"
	"github.com/gomiboko/my-circle/models"
	"github.com/gomiboko/my-circle/repositories"
	"golang.org/x/crypto/bcrypt"
)

// bcryptのコストパラメータ
const cost = 10

type UserService interface {
	CreateUser(userForm forms.UserForm) (*models.User, error)
}

type userService struct {
	userRepository repositories.UserRepository
}

func NewUserService(ur repositories.UserRepository) UserService {
	return &userService{ur}
}

func (us *userService) CreateUser(userForm forms.UserForm) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(userForm.Password), cost)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Email:        userForm.Email,
		Name:         userForm.Username,
		PasswordHash: string(hash),
	}

	err = us.userRepository.Create(&user)

	// パスワードのハッシュ値はログイン認証以外で使わないのでクリア
	user.PasswordHash = ""

	return &user, err
}
