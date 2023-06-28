package mocks

import (
	"github.com/gomiboko/my-circle/models"
	"github.com/gomiboko/my-circle/repositories"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type CircleRepositoryMock struct {
	mock.Mock
	repositories.CircleRepository
}

func NewCircleRepositoryMock(db *gorm.DB) *CircleRepositoryMock {
	return &CircleRepositoryMock{
		CircleRepository: repositories.NewCircleRepository(db),
	}
}

func (m *CircleRepositoryMock) CreateWithTran(tx *gorm.DB, circle *models.Circle) error {
	args := m.Called(tx, circle)
	return args.Error(0)
}
