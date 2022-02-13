package repositories

import (
	"errors"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/gomiboko/my-circle/db"
	"github.com/gomiboko/my-circle/models"
	"github.com/gomiboko/my-circle/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	fixtures       *testfixtures.Loader
	userRepository UserRepository
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
	s.userRepository = NewUserRepository(db)
}

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

func (s *UserRepositoryTestSuite) TestGet() {
	if err := s.fixtures.Load(); err != nil {
		s.FailNow(err.Error())
	}

	s.Run("存在するメールアドレスの場合", func() {
		user, err := s.userRepository.Get(testutils.User1Email)
		if err != nil {
			s.FailNow(err.Error())
		}

		jst, _ := time.LoadLocation("Asia/Tokyo")
		createdAt := time.Date(2021, 8, 24, 12, 34, 56, 0, jst)
		updatedAt := time.Date(2021, 8, 25, 23, 45, 01, 0, jst)

		assert.Equal(s.T(), testutils.User1Name, user.Name)
		assert.Equal(s.T(), testutils.User1Email, user.Email)
		assert.Equal(s.T(), testutils.User1PasswordHash, user.PasswordHash)
		assert.Equal(s.T(), createdAt, user.CreatedAt)
		assert.Equal(s.T(), updatedAt, user.UpdatedAt)
	})

	s.Run("存在しないメールアドレス場合", func() {
		user, err := s.userRepository.Get(testutils.UnregisteredEmail)

		assert.True(s.T(), errors.Is(err, gorm.ErrRecordNotFound))
		assert.Equal(s.T(), models.User{}, *user)
	})
}

func (s *UserRepositoryTestSuite) TestCreate() {
	if err := s.fixtures.Load(); err != nil {
		s.FailNow(err.Error())
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	if err != nil {
		s.FailNow(err.Error())
	}

	s.Run("メールアドレスが重複する場合", func() {
		user := &models.User{
			Name:         "user",
			Email:        testutils.User1Email,
			PasswordHash: string(hash),
		}
		err := s.userRepository.Create(user)

		assert.NotNil(s.T(), err)
		assert.True(s.T(), db.Is(err, db.ErrDuplicateEntry))
	})

	s.Run("メールアドレスが重複しない場合", func() {
		user := &models.User{
			Name:         "user",
			Email:        testutils.UnregisteredEmail,
			PasswordHash: string(hash),
		}
		err := s.userRepository.Create(user)

		assert.Nil(s.T(), err)
		assert.Greater(s.T(), user.ID, uint(0))
	})
}
