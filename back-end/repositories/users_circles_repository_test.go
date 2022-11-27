package repositories

import (
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/gomiboko/my-circle/models"
	"github.com/gomiboko/my-circle/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type UsersCirclesRepositoryTestSuite struct {
	suite.Suite
	fixtures               *testfixtures.Loader
	usersCirclesRepository UsersCirclesRepository
	db                     *gorm.DB
}

func (s *UsersCirclesRepositoryTestSuite) SetupSuite() {
	// テストデータの読み込み準備
	fixtures, err := testutils.GetFixtures("../testutils/fixtures")
	if err != nil {
		s.FailNow(err.Error())
	}
	s.fixtures = fixtures

	// UsersCirclesRepositoryの準備
	s.db, err = testutils.GetDB()
	if err != nil {
		s.FailNow(err.Error())
	}
	s.usersCirclesRepository = NewUsersCirclesRepository(s.db)
}

func TestUsersCirclesRepository(t *testing.T) {
	suite.Run(t, new(UsersCirclesRepositoryTestSuite))
}

func (s *UsersCirclesRepositoryTestSuite) TestCreateWithTran() {
	s.Run("トランザクションをコミットした場合", func() {
		if err := s.fixtures.Load(); err != nil {
			s.FailNow(err.Error())
		}

		usersCircles := models.UsersCircles{
			UserID:   testutils.User2ID,
			CircleID: testutils.Circle1ID,
		}

		tx := s.usersCirclesRepository.BeginTransaction()
		if err := s.usersCirclesRepository.CreateWithTran(tx, &usersCircles); err != nil {
			tx.Rollback()
			s.FailNow(err.Error())
		}
		tx.Commit()

		var usersCirclesList []models.UsersCircles
		if err := s.db.Where("user_id = ? AND circle_id = ?", testutils.User2ID, testutils.Circle1ID).Find(&usersCirclesList).Error; err != nil {
			s.FailNow(err.Error())
		}

		assert.Equal(s.T(), 1, len(usersCirclesList))
	})

	s.Run("トランザクションをロールバックした場合", func() {
		if err := s.fixtures.Load(); err != nil {
			s.FailNow(err.Error())
		}

		usersCircles := models.UsersCircles{
			UserID:   testutils.User2ID,
			CircleID: testutils.Circle1ID,
		}

		tx := s.usersCirclesRepository.BeginTransaction()
		if err := s.usersCirclesRepository.CreateWithTran(tx, &usersCircles); err != nil {
			tx.Rollback()
			s.FailNow(err.Error())
		}
		tx.Rollback()

		var usersCirclesList []models.UsersCircles
		if err := s.db.Where("user_id = ? AND circle_id = ?", testutils.User2ID, testutils.Circle1ID).Find(&usersCirclesList).Error; err != nil {
			s.FailNow(err.Error())
		}

		assert.Zero(s.T(), len(usersCirclesList))
	})
}
