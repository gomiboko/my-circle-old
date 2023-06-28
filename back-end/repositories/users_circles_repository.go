package repositories

import (
	"github.com/gomiboko/my-circle/models"
	"gorm.io/gorm"
)

type UsersCirclesRepository interface {
	RepositoryBase
	CreateWithTran(tx *gorm.DB, usersCircles *models.UsersCircles) error
}

type usersCirclesRepository struct {
	RepositoryBaseImpl
}

func NewUsersCirclesRepository(db *gorm.DB) UsersCirclesRepository {
	return &usersCirclesRepository{
		RepositoryBaseImpl: RepositoryBaseImpl{db},
	}
}

func (ucr *usersCirclesRepository) CreateWithTran(tx *gorm.DB, usersCircles *models.UsersCircles) error {
	result := tx.Create(&usersCircles)

	return result.Error
}
