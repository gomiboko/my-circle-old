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
	db             *gorm.DB
}

func (s *UserRepositoryTestSuite) SetupSuite() {
	// テストデータの読み込み準備
	fixtures, err := testutils.GetFixtures("../testutils/fixtures")
	if err != nil {
		s.FailNow(err.Error())
	}
	s.fixtures = fixtures

	// UserRepositoryの準備
	s.db, err = testutils.GetDB()
	if err != nil {
		s.FailNow(err.Error())
	}
	s.userRepository = NewUserRepository(s.db)
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
		assert.Equal(s.T(), testutils.User1RowVersion, user.RowVersion)
	})

	s.Run("存在しないメールアドレスの場合", func() {
		user, err := s.userRepository.Get(testutils.ValidEmail)

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

	s.Run("正常なデータの場合", func() {
		user := &models.User{
			Name:         "user",
			Email:        testutils.ValidEmail,
			PasswordHash: string(hash),
			IconUrl:      testutils.ValidUrl,
		}
		err := s.userRepository.Create(user)

		assert.Nil(s.T(), err)

		var createdData = models.User{}
		result := s.db.Where(&models.User{Email: testutils.ValidEmail}).First(&createdData)
		assert.Nil(s.T(), result.Error)

		assert.Greater(s.T(), createdData.ID, uint(0))
		assert.Equal(s.T(), "user", createdData.Name)
		assert.Equal(s.T(), testutils.ValidEmail, createdData.Email)
		assert.Equal(s.T(), string(hash), createdData.PasswordHash)
		assert.Equal(s.T(), testutils.ValidUrl, createdData.IconUrl)
		assert.Equal(s.T(), uint(1), createdData.RowVersion)
		assert.Zero(s.T(), len(createdData.Circles))
	})

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

		// ユーザ情報の検証(SELECTしている項目)
		assert.Equal(s.T(), testutils.User1ID, user.ID)
		assert.Equal(s.T(), testutils.User1Name, user.Name)
		assert.Equal(s.T(), testutils.User1Email, user.Email)
		assert.Equal(s.T(), testutils.User1CreatedAt, user.CreatedAt)
		assert.Equal(s.T(), testutils.User1UpdatedAt, user.UpdatedAt)
		// ユーザ情報の検証(SELECTしていない項目)
		assert.Empty(s.T(), user.PasswordHash)
		assert.Empty(s.T(), user.IconUrl)
		assert.Equal(s.T(), uint(0), user.RowVersion)

		// サークル情報の検証(SELECTしている項目)
		assert.Equal(s.T(), 1, len(user.Circles))
		assert.Equal(s.T(), testutils.Circle1ID, user.Circles[0].ID)
		assert.Equal(s.T(), testutils.Circle1Name, user.Circles[0].Name)
		assert.Equal(s.T(), testutils.Circle1IconUrl, user.Circles[0].IconUrl)
		// サークル情報の検証(SELECTしていない項目)
		assert.Equal(s.T(), time.Time{}, user.Circles[0].CreatedAt)
		assert.Equal(s.T(), time.Time{}, user.Circles[0].UpdatedAt)
		assert.Equal(s.T(), uint(0), user.Circles[0].RowVersion)
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

		assert.Equal(s.T(), 4, len(user.Circles))

		// 所属サークルが名前の昇順、IDの昇順で取得されていること
		circle1st := user.Circles[0]
		circle2nd := user.Circles[1]
		circle3rd := user.Circles[2]
		circle4th := user.Circles[3]
		assert.Equal(s.T(), "Circle03", circle1st.Name)
		assert.Equal(s.T(), "Circle03", circle2nd.Name)
		assert.Less(s.T(), circle1st.ID, circle2nd.ID)
		assert.Equal(s.T(), "Circle1", circle3rd.Name)
		assert.Equal(s.T(), "Circle2", circle4th.Name)
	})
}
