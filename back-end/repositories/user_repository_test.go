package repositories

import (
	"errors"
	"testing"

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

		assert.Equal(s.T(), testutils.User1Name, user.Name)
		assert.Equal(s.T(), testutils.User1Email, user.Email)
		assert.Equal(s.T(), testutils.User1PasswordHash, user.PasswordHash)
		assert.Equal(s.T(), "", user.IconUrl)
		assert.Equal(s.T(), testutils.User1CreatedAt, user.CreatedAt)
		assert.Equal(s.T(), testutils.User1UpdatedAt, user.UpdatedAt)
		assert.Equal(s.T(), uint(0), user.RowVersion)
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
			IconUrl:      testutils.ValidUrl,
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
			IconUrl:      testutils.ValidUrl,
		}
		err := s.userRepository.Create(user)

		assert.Nil(s.T(), err)
		assert.Greater(s.T(), user.ID, uint(0))
	})
}

func (s *UserRepositoryTestSuite) TestGetHomeInfo() {
	if err := s.fixtures.Load(); err != nil {
		s.FailNow(err.Error())
	}

	s.Run("1つのサークルに所属しているユーザの場合", func() {
		user, err := s.userRepository.GetHomeInfo(testutils.User1ID)
		if err != nil {
			s.FailNow(err.Error())
		}

		assert.Equal(s.T(), testutils.User1Name, user.Name)
		assert.Equal(s.T(), testutils.User1Email, user.Email)
		assert.Empty(s.T(), user.PasswordHash)
		assert.Equal(s.T(), testutils.User1CreatedAt, user.CreatedAt)
		assert.Equal(s.T(), testutils.User1UpdatedAt, user.UpdatedAt)

		assert.Equal(s.T(), 1, len(user.Circles))
		assert.Equal(s.T(), testutils.Circle1ID, user.Circles[0].ID)
		assert.Equal(s.T(), testutils.Circle1Name, user.Circles[0].Name)
		assert.Equal(s.T(), time.Time{}, user.Circles[0].CreatedAt)
		assert.Equal(s.T(), time.Time{}, user.Circles[0].UpdatedAt)
	})

	s.Run("サークルに所属していないユーザの場合", func() {
		user, err := s.userRepository.GetHomeInfo(testutils.User2ID)
		if err != nil {
			s.FailNow(err.Error())
		}

		assert.Empty(s.T(), user.Circles)
	})

	s.Run("複数のサークルに所属しているユーザの場合", func() {
		user, err := s.userRepository.GetHomeInfo(testutils.User3ID)
		if err != nil {
			s.FailNow(err.Error())
		}

		assert.Equal(s.T(), 3, len(user.Circles))

		// 所属サークルが名前の昇順で取得されていること
		circle1st := user.Circles[0]
		circle2nd := user.Circles[1]
		circle3rd := user.Circles[2]
		assert.Equal(s.T(), "Circle03", circle1st.Name)
		assert.Equal(s.T(), "Circle1", circle2nd.Name)
		assert.Equal(s.T(), "Circle2", circle3rd.Name)
	})
}
