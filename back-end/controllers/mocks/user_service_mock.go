package mocks

import (
	"github.com/gin-gonic/gin"
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

func (m *UserServiceMock) GetHomeInfo(userID uint) (gin.H, error) {
	args := m.Called(userID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(gin.H), args.Error(1)
	}
}
