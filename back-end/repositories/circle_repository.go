package repositories

import (
	"github.com/gomiboko/my-circle/consts"
	"github.com/gomiboko/my-circle/models"
	"github.com/gomiboko/my-circle/utils"
	"gorm.io/gorm"
)

type CircleRepository interface {
	RepositoryBase
	CreateWithTran(tx *gorm.DB, circle *models.Circle) error
	Update(circle *models.Circle) error
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

func (cr *circleRepository) Update(circle *models.Circle) error {
	result := cr.db.Model(circle).
		Where(utils.CreateRowVersionCond(circle.RowVersion)).
		Updates(models.Circle{
			Name:       circle.Name,
			IconUrl:    circle.IconUrl,
			RowVersion: circle.RowVersion + 1,
		})

	if result.RowsAffected == 0 {
		return consts.ErrUpdatedByOthers
	}

	return result.Error
}
