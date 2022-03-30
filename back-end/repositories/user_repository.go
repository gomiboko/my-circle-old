package repositories

import (
	"github.com/gomiboko/my-circle/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Get(email string) (*models.User, error)
	Create(user *models.User) error
	GetHomeInfo(userId uint) (*models.User, error)
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (ur *userRepository) Get(email string) (*models.User, error) {
	var user models.User

	// email(UQ)で検索
	cond := models.User{Email: email}
	result := ur.DB.Where(&cond).First(&user)

	return &user, result.Error
}

func (ur *userRepository) Create(user *models.User) error {
	result := ur.DB.Create(&user)

	return result.Error
}

func (ur *userRepository) GetHomeInfo(userId uint) (*models.User, error) {
	var user models.User

	cond := models.User{ID: userId}
	result := ur.DB.Where(&cond).
		Preload("Circles", func(db *gorm.DB) *gorm.DB {
			return db.Order("name")
		}).
		Table("Users").
		Select("id, name, created_at, updated_at").
		First(&user)

	return &user, result.Error
}
