package services

import (
	"testing"

	"github.com/gomiboko/my-circle/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthServiceTestSuite struct {
	suite.Suite
}

type userRepositoryMock struct {
	mock.Mock
}

func (m *userRepositoryMock) GetUser(email string) (*models.User, error) {
	args := m.Called(email)
	return args.Get(0).(*models.User), args.Error(1)
}

func TestAuthService(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}

func (s *AuthServiceTestSuite) TestAuthenticate() {
	s.Run("認証OKの場合", func() {
		user := models.User{
			PasswordHash: "$2a$10$5zIf9lXlK6F7eaMB38uRSes9ecydTeW/xDA53zADvQjrmxA/Q/BsG",
		}
		urMock := new(userRepositoryMock)
		urMock.On("GetUser", mock.AnythingOfType("string")).Return(&user, nil)

		as := NewAuthService(urMock)

		authenticataed, err := as.Authenticate("user1@example.com", "password")

		assert.True(s.T(), authenticataed)
		assert.Nil(s.T(), err)
	})

	s.Run("存在しないユーザの場合", func() {
		urMock := new(userRepositoryMock)
		urMock.On("GetUser", mock.AnythingOfType("string")).Return(new(models.User), gorm.ErrRecordNotFound)

		as := NewAuthService(urMock)

		authenticated, err := as.Authenticate("not-exist@example.com", "password")

		assert.False(s.T(), authenticated)
		assert.Nil(s.T(), err)
	})

	s.Run("パスワードが間違っている場合", func() {
		// "password"のハッシュ値とは異なるハッシュ値のユーザ
		user := models.User{
			PasswordHash: "$2a$10$5zIf9lXlK6F7eaMB38uRSeAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		}
		urMock := new(userRepositoryMock)
		urMock.On("GetUser", mock.AnythingOfType("string")).Return(&user, nil)

		as := NewAuthService(urMock)

		// "password"で認証
		authenticated, err := as.Authenticate("user1@example.com", "password")

		assert.False(s.T(), authenticated)
		assert.Nil(s.T(), err)
	})

	s.Run("DBエラーの場合", func() {
		urMock := new(userRepositoryMock)
		urMock.On("GetUser", mock.AnythingOfType("string")).Return(new(models.User), gorm.ErrInvalidDB)

		as := NewAuthService(urMock)

		authenticated, err := as.Authenticate("user1@example.com", "password")

		assert.False(s.T(), authenticated)
		assert.EqualError(s.T(), err, gorm.ErrInvalidDB.Error())
	})

	s.Run("BCryptエラーの場合", func() {
		user := models.User{
			PasswordHash: "InvalidBCryptHash",
		}
		urMock := new(userRepositoryMock)
		urMock.On("GetUser", mock.AnythingOfType("string")).Return(&user, nil)

		as := NewAuthService(urMock)

		authenticataed, err := as.Authenticate("user1@example.com", "password")

		assert.False(s.T(), authenticataed)
		assert.EqualError(s.T(), err, bcrypt.ErrHashTooShort.Error())
	})
}
