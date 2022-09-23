package services

import (
	"github.com/gomiboko/my-circle/forms"
	"github.com/gomiboko/my-circle/models"
	"github.com/gomiboko/my-circle/repositories"
)

type CircleService interface {
	CreateCircle(circleForm forms.CircleForm, userId uint) (*models.Circle, error)
	UpdateCircle(circle *models.Circle) error
}

type circleService struct {
	circleRepository       repositories.CircleRepository
	usersCirclesRepository repositories.UsersCirclesRepository
}

func NewCircleService(cr repositories.CircleRepository, ucr repositories.UsersCirclesRepository) CircleService {
	return &circleService{cr, ucr}
}

func (cs *circleService) CreateCircle(circleForm forms.CircleForm, userId uint) (*models.Circle, error) {
	circle := models.Circle{
		Name: circleForm.CircleName,
	}

	tx := cs.circleRepository.BeginTransaction()

	// サークル登録
	err := cs.circleRepository.CreateWithTran(tx, &circle)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 登録したサークルに自分を所属させる
	relation := models.UsersCircles{
		CircleID: circle.ID,
		UserID:   userId,
	}
	err = cs.usersCirclesRepository.CreateWithTran(tx, &relation)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &circle, nil
}

func (cs *circleService) UpdateCircle(circle *models.Circle) error {
	return cs.circleRepository.Update(circle)
}
