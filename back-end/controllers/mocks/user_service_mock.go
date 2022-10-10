package mocks

import (
	"github.com/gomiboko/my-circle/forms"
	"github.com/gomiboko/my-circle/models"
	"github.com/stretchr/testify/mock"
)

type UserServiceMock struct {
	mock.Mock
}

func (m *UserServiceMock) CreateUser(userForm forms.UserForm) (*models.User, error) {
	args := m.Called(userForm)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *UserServiceMock) GetHomeInfo(userID uint) (*models.User, error) {
	args := m.Called(userID)
	return args.Get(0).(*models.User), args.Error(1)
}
