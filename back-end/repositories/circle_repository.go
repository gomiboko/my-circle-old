package repositories

import (
	"github.com/gomiboko/my-circle/models"
	"gorm.io/gorm"
)

type CircleRepository interface {
	RepositoryBase
	CreateWithTran(tx *gorm.DB, circle *models.Circle) error
}

type circleRepository struct {
	RepositoryBaseImpl
}

func NewCircleRepository(db *gorm.DB) CircleRepository {
	return &circleRepository{
		RepositoryBaseImpl: RepositoryBaseImpl{db},
	}
}

func (cr *circleRepository) CreateWithTran(tx *gorm.DB, circle *models.Circle) error {
	result := tx.Create(circle)

	return result.Error
}
