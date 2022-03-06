package services

import (
	"testing"

	"github.com/gomiboko/my-circle/forms"
	"github.com/gomiboko/my-circle/services/mocks"
	"github.com/gomiboko/my-circle/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserServiceTestSuite struct {
	suite.Suite
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

func (s *UserServiceTestSuite) TestCreateUser() {
	form := forms.UserForm{
		Username: "testName",
		Email:    "testemail@example.com",
		Password: "password",
	}

	s.Run("usersテーブルへの登録に成功した場合", func() {
		urMock := new(mocks.UserRepositoryMock)
		urMock.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

		us := NewUserService(urMock)

		user, err := us.CreateUser(form)

		assert.Nil(s.T(), err)
		assert.Empty(s.T(), user.PasswordHash)
	})

	s.Run("usersテーブルへの登録に失敗した場合", func() {
		urMock := new(mocks.UserRepositoryMock)
		urMock.On("Create", mock.AnythingOfType("*models.User")).Return(testutils.ErrTest)

		us := NewUserService(urMock)

		user, err := us.CreateUser(form)

		assert.NotNil(s.T(), err)
		assert.Empty(s.T(), user.PasswordHash)
		assert.Equal(s.T(), uint(0), user.ID)
	})
}
