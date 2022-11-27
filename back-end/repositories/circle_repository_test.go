package repositories

import (
	"errors"
	"strings"
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/gomiboko/my-circle/consts"
	"github.com/gomiboko/my-circle/models"
	"github.com/gomiboko/my-circle/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type CircleRepositoryTestSuite struct {
	suite.Suite
	fixtures         *testfixtures.Loader
	circleRepository CircleRepository
	db               *gorm.DB
}

func (s *CircleRepositoryTestSuite) SetupSuite() {
	// テストデータの読み込み準備
	fixtures, err := testutils.GetFixtures("../testutils/fixtures")
	if err != nil {
		s.FailNow(err.Error())
	}
	s.fixtures = fixtures

	// CircleRepositoryの準備
	s.db, err = testutils.GetDB()
	if err != nil {
		s.FailNow(err.Error())
	}
	s.circleRepository = NewCircleRepository(s.db)
}

func TestCircleRepository(t *testing.T) {
	suite.Run(t, new(CircleRepositoryTestSuite))
}

func (s *CircleRepositoryTestSuite) TestCreateWithTran() {
	s.Run("トランザクションをコミットした場合", func() {
		if err := s.fixtures.Load(); err != nil {
			s.FailNow(err.Error())
		}

		const committedCircleName = "committed circle"
		circle := models.Circle{
			Name:    committedCircleName,
			IconUrl: testutils.ValidUrl,
		}

		tx := s.circleRepository.BeginTransaction()
		if err := s.circleRepository.CreateWithTran(tx, &circle); err != nil {
			tx.Rollback()
			s.FailNow(err.Error())
		}
		tx.Commit()

		var circles []models.Circle
		if err := s.db.Where("name = ?", committedCircleName).Find(&circles).Error; err != nil {
			s.FailNow(err.Error())
		}

		assert.Equal(s.T(), 1, len(circles))
		assert.Equal(s.T(), uint(1), circles[0].RowVersion)
		assert.Equal(s.T(), testutils.ValidUrl, circles[0].IconUrl)
	})

	s.Run("トランザクションをロールバックした場合", func() {
		if err := s.fixtures.Load(); err != nil {
			s.FailNow(err.Error())
		}

		const rolledBackCircleName = "rolled back circle"
		circle := models.Circle{
			Name:    rolledBackCircleName,
			IconUrl: testutils.ValidUrl,
		}

		tx := s.circleRepository.BeginTransaction()
		if err := s.circleRepository.CreateWithTran(tx, &circle); err != nil {
			tx.Rollback()
			s.FailNow(err.Error())
		}
		tx.Rollback()

		var circles []models.Circle
		if err := s.db.Where("name = ?", rolledBackCircleName).Find(&circles).Error; err != nil {
			s.FailNow(err.Error())
		}

		assert.Zero(s.T(), len(circles))
	})
}

func (s *CircleRepositoryTestSuite) TestUpdate() {
	s.Run("正常なデータの場合", func() {
		if err := s.fixtures.Load(); err != nil {
			s.FailNow(err.Error())
		}

		// 更新対象のデータを取得
		targetData := models.Circle{}
		if err := s.db.First(&targetData, testutils.Circle1ID).Error; err != nil {
			s.FailNow(err.Error())
		}

		assert.Equal(s.T(), testutils.Circle1Name, targetData.Name)
		assert.Equal(s.T(), testutils.Circle1IconUrl, targetData.IconUrl)
		assert.Equal(s.T(), testutils.Circle1RowVersion, targetData.RowVersion)

		const updatedCircleName = "updated circle name"
		targetData.Name = updatedCircleName
		targetData.IconUrl = testutils.ValidUrl
		if err := s.circleRepository.Update(&targetData); err != nil {
			s.FailNow(err.Error())
		}

		// 更新後のデータを取得
		updatedData := models.Circle{}
		if err := s.db.First(&updatedData, testutils.Circle1ID).Error; err != nil {
			s.FailNow(err.Error())
		}

		assert.Equal(s.T(), updatedCircleName, updatedData.Name)
		assert.Equal(s.T(), testutils.ValidUrl, updatedData.IconUrl)
		assert.Equal(s.T(), testutils.Circle1RowVersion+1, updatedData.RowVersion)
	})

	s.Run("楽観ロックエラーとなった場合", func() {
		if err := s.fixtures.Load(); err != nil {
			s.FailNow(err.Error())
		}

		// 更新対象データを取得
		targetData := models.Circle{}
		if err := s.db.First(&targetData, testutils.Circle1ID).Error; err != nil {
			s.FailNow(err.Error())
		}

		assert.Equal(s.T(), testutils.Circle1Name, targetData.Name)
		assert.Equal(s.T(), testutils.Circle1IconUrl, targetData.IconUrl)
		assert.Equal(s.T(), testutils.Circle1RowVersion, targetData.RowVersion)

		// 更新対象のデータが別の処理で先に更新される
		const updatedNameByOthers = "updated by others"
		const updatedUrlByOthers = "https://example.com/updated/by/others"
		copiedTargetData := targetData
		copiedTargetData.Name = updatedNameByOthers
		copiedTargetData.IconUrl = updatedUrlByOthers
		if err := s.circleRepository.Update(&copiedTargetData); err != nil {
			s.FailNow(err.Error())
		}

		updatedDataByOthers := models.Circle{}
		if err := s.db.First(&updatedDataByOthers, testutils.Circle1ID).Error; err != nil {
			s.FailNow(err.Error())
		}

		assert.Equal(s.T(), updatedNameByOthers, updatedDataByOthers.Name)
		assert.Equal(s.T(), updatedUrlByOthers, updatedDataByOthers.IconUrl)

		// 別の処理で更新されてしまっているデータを更新
		targetData.Name = "updated circle name"
		targetData.IconUrl = testutils.ValidUrl
		err := s.circleRepository.Update(&targetData)

		// 更新処理後のデータを取得
		updatedData := models.Circle{}
		if err := s.db.First(&updatedData, testutils.Circle1ID).Error; err != nil {
			s.FailNow(err.Error())
		}

		// 楽観ロックエラーが発生すること、後から実行した更新処理が反映されていないことを確認
		assert.True(s.T(), errors.Is(err, consts.ErrUpdatedByOthers))
		assert.Equal(s.T(), updatedNameByOthers, updatedData.Name)
		assert.Equal(s.T(), updatedUrlByOthers, updatedData.IconUrl)
	})

	s.Run("予期せぬエラーが発生した場合", func() {
		if err := s.fixtures.Load(); err != nil {
			s.FailNow(err.Error())
		}

		targetData := models.Circle{}
		if err := s.db.First(&targetData, testutils.Circle1ID).Error; err != nil {
			s.FailNow(err.Error())
		}

		assert.Equal(s.T(), testutils.Circle1Name, targetData.Name)

		// 桁数オーバーとなる値を設定して更新
		targetData.Name = strings.Repeat("a", testutils.CircleNameMaxLength+1)
		err := s.circleRepository.Update(&targetData)

		// 更新処理後のデータを取得
		updatedData := models.Circle{}
		if err := s.db.First(&updatedData, testutils.Circle1ID).Error; err != nil {
			s.FailNow(err.Error())
		}

		assert.NotNil(s.T(), err)
		assert.False(s.T(), errors.Is(err, consts.ErrUpdatedByOthers))
		assert.Equal(s.T(), testutils.Circle1Name, updatedData.Name)
	})
}
