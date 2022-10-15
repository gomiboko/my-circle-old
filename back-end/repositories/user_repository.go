package repositories

import (
	"github.com/gomiboko/my-circle/consts"
	"github.com/gomiboko/my-circle/models"
	"github.com/gomiboko/my-circle/utils"
	"gorm.io/gorm"
)

type UserRepository interface {
	Get(email string) (*models.User, error)
	Create(user *models.User) error
	GetHomeInfo(userID uint) (*models.User, error)
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

func (ur *userRepository) GetHomeInfo(userID uint) (*models.User, error) {
	var user models.User

	cond := models.User{ID: userID}
	result := ur.DB.Where(&cond).
		Preload(consts.ModelNameCircles, func(db *gorm.DB) *gorm.DB {
			orderStr := utils.CreateOrderStr(consts.ColumnNameUsersName, consts.ColumnNameCommonID)
			return db.
				Select(
					consts.ColumnNameCommonID,
					consts.ColumnNameUsersName,
					consts.ColumnNameUsersIconUrl).
				Order(orderStr)
		}).
		Table(consts.TableNameUsers).
		Select(
			consts.ColumnNameCommonID,
			consts.ColumnNameUsersName,
			consts.ColumnNameUsersEmail,
			consts.ColumnNameCommonCreatedAt,
			consts.ColumnNameCommonUpdatedAt).
		First(&user)

	return &user, result.Error
}
