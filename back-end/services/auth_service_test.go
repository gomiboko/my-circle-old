package services

import (
	"testing"

	"github.com/gomiboko/my-circle/models"
	"github.com/gomiboko/my-circle/testutils"
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

func (m *userRepositoryMock) Get(email string) (*models.User, error) {
	args := m.Called(email)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *userRepositoryMock) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func TestAuthService(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}

func (s *AuthServiceTestSuite) TestAuthenticate() {
	s.Run("認証OKの場合", func() {
		user := models.User{
			ID:           1,
			PasswordHash: testutils.User1PasswordHash,
		}
		urMock := new(userRepositoryMock)
		urMock.On("Get", mock.AnythingOfType("string")).Return(&user, nil)

		as := NewAuthService(urMock)

		userID, err := as.Authenticate(testutils.User1Email, testutils.User1Password)

		assert.Equal(s.T(), uint(1), *userID)
		assert.Nil(s.T(), err)
	})

	s.Run("存在しないユーザの場合", func() {
		urMock := new(userRepositoryMock)
		urMock.On("Get", mock.AnythingOfType("string")).Return(new(models.User), gorm.ErrRecordNotFound)

		as := NewAuthService(urMock)

		userID, err := as.Authenticate(testutils.UnregisteredEmail, testutils.User1Password)

		assert.Nil(s.T(), userID)
		assert.Nil(s.T(), err)
	})

	s.Run("パスワードが間違っている場合", func() {
		// "password"のハッシュ値とは異なるハッシュ値のユーザ
		user := models.User{
			PasswordHash: "$2a$10$5zIf9lXlK6F7eaMB38uRSeAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		}
		urMock := new(userRepositoryMock)
		urMock.On("Get", mock.AnythingOfType("string")).Return(&user, nil)

		as := NewAuthService(urMock)

		// "password"で認証
		userID, err := as.Authenticate(testutils.User1Email, testutils.User1Password)

		assert.Nil(s.T(), userID)
		assert.Nil(s.T(), err)
	})

	s.Run("DBエラーの場合", func() {
		urMock := new(userRepositoryMock)
		urMock.On("Get", mock.AnythingOfType("string")).Return(new(models.User), gorm.ErrInvalidDB)

		as := NewAuthService(urMock)

		userID, err := as.Authenticate(testutils.User1Email, testutils.User1Password)

		assert.Nil(s.T(), userID)
		assert.EqualError(s.T(), err, gorm.ErrInvalidDB.Error())
	})

	s.Run("BCryptエラーの場合", func() {
		user := models.User{
			PasswordHash: "InvalidBCryptHash",
		}
		urMock := new(userRepositoryMock)
		urMock.On("Get", mock.AnythingOfType("string")).Return(&user, nil)

		as := NewAuthService(urMock)

		userID, err := as.Authenticate(testutils.User1Email, testutils.User1Password)

		assert.Nil(s.T(), userID)
		assert.EqualError(s.T(), err, bcrypt.ErrHashTooShort.Error())
	})
}
