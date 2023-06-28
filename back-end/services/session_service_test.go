package services

import (
	"testing"

	"github.com/gomiboko/my-circle/models"
	"github.com/gomiboko/my-circle/services/mocks"
	"github.com/gomiboko/my-circle/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SessionServiceTestSuite struct {
	suite.Suite
}

func TestSessionService(t *testing.T) {
	suite.Run(t, new(SessionServiceTestSuite))
}

func (s *SessionServiceTestSuite) TestAuthenticate() {
	s.Run("認証OKの場合", func() {
		user := models.User{
			ID:           1,
			PasswordHash: testutils.User1PasswordHash,
		}
		urMock := new(mocks.UserRepositoryMock)
		urMock.On("Get", mock.AnythingOfType("string")).Return(&user, nil)

		ss := NewSessionService(urMock)

		userID, err := ss.Authenticate(testutils.User1Email, testutils.User1Password)

		assert.Equal(s.T(), uint(1), *userID)
		assert.Nil(s.T(), err)
	})

	s.Run("存在しないユーザの場合", func() {
		urMock := new(mocks.UserRepositoryMock)
		urMock.On("Get", mock.AnythingOfType("string")).Return(new(models.User), gorm.ErrRecordNotFound)

		ss := NewSessionService(urMock)

		userID, err := ss.Authenticate(testutils.ValidEmail, testutils.User1Password)

		assert.Nil(s.T(), userID)
		assert.Nil(s.T(), err)
	})

	s.Run("パスワードが間違っている場合", func() {
		// "password"のハッシュ値とは異なるハッシュ値のユーザ
		user := models.User{
			PasswordHash: "$2a$10$5zIf9lXlK6F7eaMB38uRSeAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		}
		urMock := new(mocks.UserRepositoryMock)
		urMock.On("Get", mock.AnythingOfType("string")).Return(&user, nil)

		ss := NewSessionService(urMock)

		// "password"で認証
		userID, err := ss.Authenticate(testutils.User1Email, testutils.User1Password)

		assert.Nil(s.T(), userID)
		assert.Nil(s.T(), err)
	})

	s.Run("DBエラーの場合", func() {
		urMock := new(mocks.UserRepositoryMock)
		urMock.On("Get", mock.AnythingOfType("string")).Return(new(models.User), gorm.ErrInvalidDB)

		ss := NewSessionService(urMock)

		userID, err := ss.Authenticate(testutils.User1Email, testutils.User1Password)

		assert.Nil(s.T(), userID)
		assert.EqualError(s.T(), err, gorm.ErrInvalidDB.Error())
	})

	s.Run("BCryptエラーの場合", func() {
		user := models.User{
			PasswordHash: "InvalidBCryptHash",
		}
		urMock := new(mocks.UserRepositoryMock)
		urMock.On("Get", mock.AnythingOfType("string")).Return(&user, nil)

		ss := NewSessionService(urMock)

		userID, err := ss.Authenticate(testutils.User1Email, testutils.User1Password)

		assert.Nil(s.T(), userID)
		assert.EqualError(s.T(), err, bcrypt.ErrHashTooShort.Error())
	})
}
