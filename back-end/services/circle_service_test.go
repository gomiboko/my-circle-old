package services

import (
	"errors"
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/gomiboko/my-circle/forms"
	"github.com/gomiboko/my-circle/models"
	"github.com/gomiboko/my-circle/repositories"
	"github.com/gomiboko/my-circle/services/mocks"
	"github.com/gomiboko/my-circle/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type CircleServiceTestSuite struct {
	suite.Suite
	fixtures               *testfixtures.Loader
	circleRepository       repositories.CircleRepository
	usersCirclesRepository repositories.UsersCirclesRepository
	db                     *gorm.DB
}

func (s *CircleServiceTestSuite) SetupSuite() {
	// テストデータの読み込み準備
	fixtures, err := testutils.GetFixtures("../testutils/fixtures")
	if err != nil {
		s.FailNow(err.Error())
	}
	s.fixtures = fixtures

	// トランザクションのテスト用にモックではないRepositoryの準備
	s.db, err = testutils.GetDB()
	if err != nil {
		s.FailNow(err.Error())
	}
	s.circleRepository = repositories.NewCircleRepository(s.db)
	s.usersCirclesRepository = repositories.NewUsersCirclesRepository(s.db)
}

func TestCircleService(t *testing.T) {
	suite.Run(t, new(CircleServiceTestSuite))
}

func (s *CircleServiceTestSuite) TestCreateCircle() {
	s.Run("circlesテーブル、users_circlesテーブルへの登録に成功した場合", func() {
		if err := s.fixtures.Load(); err != nil {
			s.FailNow(err.Error())
		}

		cs := NewCircleService(s.circleRepository, s.usersCirclesRepository)

		const testCircleName = "test circle"
		circleForm := forms.CircleForm{
			CircleName:     testCircleName,
			CircleIconFile: nil,
		}
		circle, err := cs.CreateCircle(circleForm, testutils.User1ID)

		require.Nil(s.T(), err)
		require.NotNil(s.T(), circle)

		// 登録した circles データを取得
		var circles []models.Circle
		err = s.db.Where("name = ?", testCircleName).Find(&circles).Error
		require.Nil(s.T(), err)

		assert.Equal(s.T(), 1, len(circles))

		// 登録した users_circles データを取得
		var usersCirclesList []models.UsersCircles
		err = s.db.Where("user_id = ? AND circle_id = ?", testutils.User1ID, circles[0].ID).Find(&usersCirclesList).Error
		require.Nil(s.T(), err)

		assert.Equal(s.T(), 1, len(usersCirclesList))
	})

	s.Run("circlesテーブルへの登録に失敗した場合", func() {
		if err := s.fixtures.Load(); err != nil {
			s.FailNow(err.Error())
		}

		// テスト結果の検証用にサークル登録処理前のデータを取得しておく
		var allUsersCirclesBefore []models.UsersCircles
		err := s.db.Find(&allUsersCirclesBefore).Error
		require.Nil(s.T(), err)

		crMock := mocks.NewCircleRepositoryMock(s.db)
		crMock.On("CreateWithTran", mock.AnythingOfType("*gorm.DB"), mock.AnythingOfType("*models.Circle")).Return(testutils.ErrTest)
		cs := NewCircleService(crMock, s.usersCirclesRepository)

		const testCircleName = "test circle"
		circleForm := forms.CircleForm{
			CircleName:     testCircleName,
			CircleIconFile: nil,
		}
		circle, err := cs.CreateCircle(circleForm, testutils.User1ID)

		assert.NotNil(s.T(), err)
		assert.True(s.T(), errors.Is(testutils.ErrTest, err))
		assert.Nil(s.T(), circle)

		// サークル登録処理後の users_circles データを取得
		var allUsersCirclesAfter []models.UsersCircles
		err = s.db.Find(&allUsersCirclesAfter).Error
		require.Nil(s.T(), err)

		// サークル登録処理後の circles データを取得
		var circles []models.Circle
		err = s.db.Where("name = ?", testCircleName).Find(&circles).Error
		require.Nil(s.T(), err)

		assert.Zero(s.T(), len(circles))
		assert.Equal(s.T(), len(allUsersCirclesBefore), len(allUsersCirclesAfter))
	})

	s.Run("circlesテーブルへの登録に成功し、users_circlesテーブルへの登録に失敗した場合", func() {
		if err := s.fixtures.Load(); err != nil {
			s.FailNow(err.Error())
		}

		// テスト結果の検証用にサークル登録処理前のデータ取得しておく
		var allUsersCirclesBefore []models.UsersCircles
		err := s.db.Find(&allUsersCirclesBefore).Error
		require.Nil(s.T(), err)

		ucrMock := mocks.NewUsersCirclesRepositoryMock(s.db)
		ucrMock.On("CreateWithTran", mock.AnythingOfType("*gorm.DB"), mock.AnythingOfType("*models.UsersCircles")).Return(testutils.ErrTest)
		cs := NewCircleService(s.circleRepository, ucrMock)

		const testCircleName = "test circle"
		circleForm := forms.CircleForm{
			CircleName:     testCircleName,
			CircleIconFile: nil,
		}
		circle, err := cs.CreateCircle(circleForm, testutils.User1ID)

		assert.NotNil(s.T(), err)
		assert.True(s.T(), errors.Is(testutils.ErrTest, err))
		assert.Nil(s.T(), circle)

		// サークル登録処理後の users_circles データを取得
		var allUsersCirclesAfter []models.UsersCircles
		err = s.db.Find(&allUsersCirclesAfter).Error
		require.Nil(s.T(), err)

		// サークル登録処理後の users_circles データを取得
		var circles []models.Circle
		err = s.db.Where("name = ?", testCircleName).Find(&circles).Error
		require.Nil(s.T(), err)

		assert.Zero(s.T(), len(circles))
		assert.Equal(s.T(), len(allUsersCirclesBefore), len(allUsersCirclesAfter))
	})
}
