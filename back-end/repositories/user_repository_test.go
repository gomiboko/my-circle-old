package repositories

import (
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/gomiboko/my-circle/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	fixtures *testfixtures.Loader
	ur       UserRepository
}

func (s *UserRepositoryTestSuite) SetupSuite() {
	// テストデータの読み込み準備
	fixtures, err := testutils.GetFixtures("../testutils/fixtures")
	if err != nil {
		s.FailNow(err.Error())
	}
	s.fixtures = fixtures

	// UserRepositoryの準備
	db, err := testutils.GetDB()
	if err != nil {
		s.FailNow(err.Error())
	}
	s.ur = NewUserRepository(db)
}

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

func (s *UserRepositoryTestSuite) TestGetUser() {
	err := s.fixtures.Load()
	if err != nil {
		s.FailNow(err.Error())
	}

	s.Run("存在するメールアドレスの場合", func() {
		user, err := s.ur.GetUser("user1@example.com")
		if err != nil {
			s.FailNow(err.Error())
		}

		jst, _ := time.LoadLocation("Asia/Tokyo")
		createdAt := time.Date(2021, 8, 24, 12, 34, 56, 0, jst)
		updatedAt := time.Date(2021, 8, 25, 23, 45, 01, 0, jst)

		assert.Equal(s.T(), "user1", user.Name)
		assert.Equal(s.T(), "user1@example.com", user.Email)
		assert.Equal(s.T(), "$2a$10$5zIf9lXlK6F7eaMB38uRSes9ecydTeW/xDA53zADvQjrmxA/Q/BsG", user.PasswordHash)
		assert.Equal(s.T(), createdAt, user.CreatedAt)
		assert.Equal(s.T(), updatedAt, user.UpdatedAt)
	})

	s.Run("存在しないメールアドレス場合", func() {
		user, err := s.ur.GetUser("not-exist@example.com")
		if err != nil {
			s.FailNow(err.Error())
		}

		assert.Nil(s.T(), user)
	})
}
