package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gomiboko/my-circle/controllers/mocks"
	"github.com/gomiboko/my-circle/forms"
	"github.com/gomiboko/my-circle/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// SessionControllerテストスイート
type SessionControllerTestSuite struct {
	suite.Suite
}

func (s *SessionControllerTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func TestSessionController(t *testing.T) {
	suite.Run(t, new(SessionControllerTestSuite))
}

func (s *SessionControllerTestSuite) TestLogin() {
	const reqPath = "/login"

	s.Run("不正なリクエスト(URLエンコード)の場合", func() {
		var userID *uint = nil
		ssMock := new(mocks.SessionServiceMock)
		ssMock.On("Authenticate", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(userID, nil)

		sc := NewSessionController(ssMock)

		values := url.Values{}
		values.Set("email", testutils.ValidEmail)
		values.Add("password", testutils.ValidPassword)

		// URLエンコードで送信
		r := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(r)
		c.Request, _ = http.NewRequest(http.MethodPost, reqPath, strings.NewReader(values.Encode()))
		c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		sc.Create(c)

		var res testutils.ApiErrorReponse
		json.Unmarshal(r.Body.Bytes(), &res)

		assert.Equal(s.T(), http.StatusBadRequest, r.Code)
		assert.Equal(s.T(), "不正なリクエストです", res.Message)
	})

	s.Run("不正な入力値の場合", func() {
		inputs := []forms.LoginForm{
			// メールアドレスのチェックデータ
			{Password: testutils.ValidPassword, Email: ""},
			{Password: testutils.ValidPassword, Email: "isNotEmail"},
			{Password: testutils.ValidPassword, Email: testutils.CreateEmailAddress(testutils.EmailMaxLength + 1)},
			// パスワードのチェックデータ
			{Email: testutils.ValidEmail, Password: ""},
			{Email: testutils.ValidEmail, Password: strings.Repeat("a", testutils.PasswordMinLength-1)},
			{Email: testutils.ValidEmail, Password: strings.Repeat("a", testutils.PasswordMaxLength+1)},
			{Email: testutils.ValidEmail, Password: "にほんごぱすわーど"},
		}

		var userID *uint = nil
		ssMock := new(mocks.SessionServiceMock)
		ssMock.On("Authenticate", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(userID, nil)

		sc := NewSessionController(ssMock)

		for _, in := range inputs {
			reqBody, err := testutils.CreateRequestBodyStr(in)
			if err != nil {
				s.FailNow(err.Error())
			}
			r, c := testutils.CreatePostContext(reqPath, reqBody)

			sc.Create(c)

			var res testutils.ApiErrorReponse
			json.Unmarshal(r.Body.Bytes(), &res)

			assert.Equal(s.T(), http.StatusUnauthorized, r.Code)
			assert.Equal(s.T(), "メールアドレスまたはパスワードが違います", res.Message)
		}
	})

	s.Run("正常な入力値の場合", func() {
		inputs := []forms.LoginForm{
			// メールアドレスのチェックデータ
			{Password: testutils.ValidPassword, Email: testutils.CreateEmailAddress(testutils.EmailMaxLength)},
			{Password: testutils.ValidPassword, Email: "にほんご@example.com"},
			// パスワードのチェックデータ
			{Email: testutils.ValidEmail, Password: strings.Repeat("a", testutils.PasswordMinLength)},
			{Email: testutils.ValidEmail, Password: strings.Repeat("a", testutils.PasswordMaxLength)},
		}

		// SessionServiceモック
		userID := uint(1)
		ssMock := new(mocks.SessionServiceMock)
		ssMock.On("Authenticate", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(&userID, nil)

		sc := NewSessionController(ssMock)

		for _, in := range inputs {
			reqBody, err := testutils.CreateRequestBodyStr(in)
			if err != nil {
				s.FailNow(err.Error())
			}
			r, c := testutils.CreatePostContext(reqPath, reqBody)

			// sessions.sessionモック
			sessMock := mocks.NewSessionMock()
			testutils.SetSessionMockToGin(c, sessMock)

			sc.Create(c)
			c.Writer.WriteHeaderNow()

			assert.Equal(s.T(), http.StatusCreated, r.Code)
			assert.Equal(s.T(), 0, r.Body.Len())
			sessMock.AssertCalled(s.T(), "Set", mock.AnythingOfType("string"), mock.AnythingOfType("uint"))
			sessMock.AssertCalled(s.T(), "Save")
		}
	})

	s.Run("認証失敗の場合", func() {
		var userID *uint = nil
		ssMock := new(mocks.SessionServiceMock)
		ssMock.On("Authenticate", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(userID, nil)

		sc := NewSessionController(ssMock)

		form := forms.LoginForm{Email: testutils.ValidEmail, Password: testutils.ValidPassword}
		reqBody, err := testutils.CreateRequestBodyStr(form)
		if err != nil {
			s.FailNow(err.Error())
		}
		r, c := testutils.CreatePostContext(reqPath, reqBody)

		sc.Create(c)

		var res testutils.ApiErrorReponse
		json.Unmarshal(r.Body.Bytes(), &res)

		assert.Equal(s.T(), http.StatusUnauthorized, r.Code)
		assert.Equal(s.T(), "メールアドレスまたはパスワードが違います", res.Message)
	})

	s.Run("予期せぬエラーが発生した場合", func() {
		var userID *uint = nil
		ssMock := new(mocks.SessionServiceMock)
		ssMock.On("Authenticate", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(userID, testutils.ErrTest)

		sc := NewSessionController(ssMock)

		form := forms.LoginForm{Email: testutils.ValidEmail, Password: testutils.ValidPassword}
		reqBody, err := testutils.CreateRequestBodyStr(form)
		if err != nil {
			s.FailNow(err.Error())
		}
		r, c := testutils.CreatePostContext(reqPath, reqBody)

		sc.Create(c)

		var res testutils.ApiErrorReponse
		json.Unmarshal(r.Body.Bytes(), &res)

		assert.Equal(s.T(), http.StatusInternalServerError, r.Code)
		assert.Equal(s.T(), testutils.UnexpectedErrMsg, res.Message)
	})
}

func (s *SessionControllerTestSuite) TestLogout() {
	var userID *uint = nil
	ssMock := new(mocks.SessionServiceMock)
	ssMock.On("Authenticate", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(userID, nil)

	sc := NewSessionController(ssMock)

	r, c := testutils.CreateGetContext("/logout")

	// sessions.sessionモック
	sessMock := mocks.NewSessionMock()
	testutils.SetSessionMockToGin(c, sessMock)

	sc.Destroy(c)

	var res testutils.ApiErrorReponse
	json.Unmarshal(r.Body.Bytes(), &res)

	assert.Equal(s.T(), http.StatusOK, r.Code)
	sessMock.AssertCalled(s.T(), "Save")
	sessMock.AssertCalled(s.T(), "Clear")
}
