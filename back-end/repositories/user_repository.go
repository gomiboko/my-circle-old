package repositories

import (
	"github.com/gomiboko/my-circle/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUser(email string) (*models.User, error)
	CreateUser(user *models.User) error
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (ur *userRepository) GetUser(email string) (*models.User, error) {
	var user models.User

	// email(UQ)で検索
	cond := models.User{Email: email}
	result := ur.DB.Where(&cond).First(&user)

	return &user, result.Error
}

func (ur *userRepository) CreateUser(user *models.User) error {
	result := ur.DB.Create(&user)

	return result.Error
}
