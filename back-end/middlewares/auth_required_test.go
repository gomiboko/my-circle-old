package middlewares

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gomiboko/my-circle/controllers/mocks"
	"github.com/gomiboko/my-circle/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AuthRequiredTestSuite struct {
	suite.Suite
}

func (s *AuthRequiredTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func TestAuthRequired(t *testing.T) {
	suite.Run(t, new(AuthRequiredTestSuite))
}

func (s *AuthRequiredTestSuite) TestAuthRequired() {
	s.Run("未ログインの場合", func() {
		r, c := testutils.CreateGetContext("/")

		sessMock := mocks.NewSessionMock()
		sessMock.On("Get", mock.AnythingOfType("string")).Return(nil)
		testutils.SetSessionMockToGin(c, sessMock)

		AuthRequired()(c)

		var res testutils.ApiErrorReponse
		json.Unmarshal(r.Body.Bytes(), &res)

		assert.Equal(s.T(), http.StatusUnauthorized, r.Code)
		assert.Equal(s.T(), "ログインしてください", res.Message)
	})

	s.Run("ログイン済みの場合", func() {
		r, c := testutils.CreateGetContext("/")

		sessMock := mocks.NewSessionMock()
		sessMock.On("Get", mock.AnythingOfType("string")).Return(interface{}(uint(1)))
		testutils.SetSessionMockToGin(c, sessMock)

		AuthRequired()(c)

		assert.NotEqual(s.T(), http.StatusUnauthorized, r.Code)
	})
}
