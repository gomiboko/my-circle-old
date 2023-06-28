package mocks

import (
	"github.com/gomiboko/my-circle/models"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) Get(email string) (*models.User, error) {
	args := m.Called(email)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *UserRepositoryMock) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserRepositoryMock) GetHomeInfo(userID uint) (*models.User, error) {
	args := m.Called(userID)
	return args.Get(0).(*models.User), args.Error(1)
}
