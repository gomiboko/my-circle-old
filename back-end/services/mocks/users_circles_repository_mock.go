package mocks

import (
	"github.com/gomiboko/my-circle/models"
	"github.com/gomiboko/my-circle/repositories"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type UsersCirclesRepositoryMock struct {
	mock.Mock
	repositories.UsersCirclesRepository
}

func NewUsersCirclesRepositoryMock(db *gorm.DB) *UsersCirclesRepositoryMock {
	return &UsersCirclesRepositoryMock{
		UsersCirclesRepository: repositories.NewUsersCirclesRepository(db),
	}
}

func (m *UsersCirclesRepositoryMock) CreateWithTran(tx *gorm.DB, usersCircles *models.UsersCircles) error {
	args := m.Called(tx, usersCircles)
	return args.Error(0)
}
