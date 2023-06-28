package services

import (
	"github.com/gin-gonic/gin"
	"github.com/gomiboko/my-circle/consts"
	"github.com/gomiboko/my-circle/forms"
	"github.com/gomiboko/my-circle/models"
	"github.com/gomiboko/my-circle/repositories"
	"github.com/gomiboko/my-circle/responses"
	"golang.org/x/crypto/bcrypt"
)

// bcryptのコストパラメータ
const cost = 10

type UserService interface {
	CreateUser(userForm forms.UserForm) (*models.User, error)
	GetHomeInfo(userID uint) (gin.H, error)
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

func (us *userService) GetHomeInfo(userID uint) (gin.H, error) {
	user, err := us.userRepository.GetHomeInfo(userID)
	if err != nil {
		return nil, err
	}

	// レスポンス用に整形
	var resCircles []responses.Circle
	for _, c := range user.Circles {
		resC := responses.Circle{
			ID:      c.ID,
			Name:    c.Name,
			IconUrl: c.IconUrl,
		}
		resCircles = append(resCircles, resC)
	}

	res := gin.H{
		consts.ResKeyUserName:    user.Name,
		consts.ResKeyUserIconUrl: user.IconUrl,
		consts.ResKeyCircles:     resCircles,
	}

	return res, nil
}
