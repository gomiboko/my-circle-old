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

// AuthControllerテストスイート
type AuthControllerTestSuite struct {
	suite.Suite
}

func (s *AuthControllerTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func TestAuthController(t *testing.T) {
	suite.Run(t, new(AuthControllerTestSuite))
}

func (s *AuthControllerTestSuite) TestLogin() {
	const reqPath = "/login"

	s.Run("不正なリクエスト(URLエンコード)の場合", func() {
		var userID *uint = nil
		asMock := new(mocks.AuthServiceMock)
		asMock.On("Authenticate", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(userID, nil)

		ac := NewAuthController(asMock)

		values := url.Values{}
		values.Set("email", testutils.ValidEmail)
		values.Add("password", testutils.ValidPassword)

		// URLエンコードで送信
		r := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(r)
		c.Request, _ = http.NewRequest(http.MethodPost, reqPath, strings.NewReader(values.Encode()))
		c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		ac.Login(c)

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
		asMock := new(mocks.AuthServiceMock)
		asMock.On("Authenticate", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(userID, nil)

		ac := NewAuthController(asMock)

		for _, in := range inputs {
			reqBody, err := testutils.CreateRequestBodyStr(in)
			if err != nil {
				s.FailNow(err.Error())
			}
			r, c := testutils.CreatePostContext(reqPath, reqBody)

			ac.Login(c)

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

		// AuthServiceモック
		userID := uint(1)
		asMock := new(mocks.AuthServiceMock)
		asMock.On("Authenticate", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(&userID, nil)

		ac := NewAuthController(asMock)

		for _, in := range inputs {
			reqBody, err := testutils.CreateRequestBodyStr(in)
			if err != nil {
				s.FailNow(err.Error())
			}
			r, c := testutils.CreatePostContext(reqPath, reqBody)

			// sessions.sessionモック
			sessMock := mocks.NewSessionMock()
			testutils.SetSessionMockToGin(c, sessMock)

			ac.Login(c)
			c.Writer.WriteHeaderNow()

			assert.Equal(s.T(), http.StatusCreated, r.Code)
			assert.Equal(s.T(), 0, r.Body.Len())
			sessMock.AssertCalled(s.T(), "Set", mock.AnythingOfType("string"), mock.AnythingOfType("uint"))
			sessMock.AssertCalled(s.T(), "Save")
		}
	})

	s.Run("認証失敗の場合", func() {
		var userID *uint = nil
		asMock := new(mocks.AuthServiceMock)
		asMock.On("Authenticate", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(userID, nil)

		ac := NewAuthController(asMock)

		form := forms.LoginForm{Email: testutils.ValidEmail, Password: testutils.ValidPassword}
		reqBody, err := testutils.CreateRequestBodyStr(form)
		if err != nil {
			s.FailNow(err.Error())
		}
		r, c := testutils.CreatePostContext(reqPath, reqBody)

		ac.Login(c)

		var res testutils.ApiErrorReponse
		json.Unmarshal(r.Body.Bytes(), &res)

		assert.Equal(s.T(), http.StatusUnauthorized, r.Code)
		assert.Equal(s.T(), "メールアドレスまたはパスワードが違います", res.Message)
	})

	s.Run("予期せぬエラーが発生した場合", func() {
		var userID *uint = nil
		asMock := new(mocks.AuthServiceMock)
		asMock.On("Authenticate", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(userID, testutils.ErrTest)

		ac := NewAuthController(asMock)

		form := forms.LoginForm{Email: testutils.ValidEmail, Password: testutils.ValidPassword}
		reqBody, err := testutils.CreateRequestBodyStr(form)
		if err != nil {
			s.FailNow(err.Error())
		}
		r, c := testutils.CreatePostContext(reqPath, reqBody)

		ac.Login(c)

		var res testutils.ApiErrorReponse
		json.Unmarshal(r.Body.Bytes(), &res)

		assert.Equal(s.T(), r.Code, http.StatusInternalServerError)
		assert.Equal(s.T(), testutils.UnexpectedErrMsg, res.Message)
	})
}

func (s *AuthControllerTestSuite) TestLogout() {
	var userID *uint = nil
	asMock := new(mocks.AuthServiceMock)
	asMock.On("Authenticate", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(userID, nil)

	ac := NewAuthController(asMock)

	r, c := testutils.CreateGetContext("/logout")

	// sessions.sessionモック
	sessMock := mocks.NewSessionMock()
	testutils.SetSessionMockToGin(c, sessMock)

	ac.Logout(c)

	var res testutils.ApiErrorReponse
	json.Unmarshal(r.Body.Bytes(), &res)

	assert.Equal(s.T(), http.StatusOK, r.Code)
	assert.Equal(s.T(), "logged out", res.Message)
	sessMock.AssertCalled(s.T(), "Save")
	sessMock.AssertCalled(s.T(), "Clear")
}
