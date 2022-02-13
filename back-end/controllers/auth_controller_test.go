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

const (
	validEmail        = "foo@example.com"
	validPassword     = "password123"
	emailMaxLength    = 254
	passwordMinLength = 8
	passwordMaxLength = 64
)

type apiResponse struct {
	Message string
}

// AuthControllerテストスイート
type AuthControllerTestSuite struct {
	suite.Suite
}

func (s *AuthControllerTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

// AuthServiceモック
type authServiceMock struct {
	mock.Mock
}

func (m *authServiceMock) Authenticate(email string, password string) (*uint, error) {
	args := m.Called(email, password)
	return args.Get(0).(*uint), args.Error(1)
}

func TestAuthController(t *testing.T) {
	suite.Run(t, new(AuthControllerTestSuite))
}

func (s *AuthControllerTestSuite) TestLogin() {
	s.Run("不正なリクエスト(URLエンコード)の場合", func() {
		ac := NewAuthController(nil)

		values := url.Values{}
		values.Set("email", validEmail)
		values.Add("password", validPassword)

		// URLエンコードで送信
		r := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(r)
		c.Request, _ = http.NewRequest(http.MethodPost, "/login", strings.NewReader(values.Encode()))
		c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		ac.Login(c)

		var res apiResponse
		json.Unmarshal(r.Body.Bytes(), &res)

		assert.Equal(s.T(), http.StatusBadRequest, r.Code)
		assert.Equal(s.T(), "不正なリクエストです", res.Message)
	})

	s.Run("不正な入力値の場合", func() {
		inputs := []forms.LoginForm{
			// メールアドレスのチェックデータ
			{Password: validPassword, Email: ""},
			{Password: validPassword, Email: "isNotEmail"},
			{Password: validPassword, Email: createEmailAddress(emailMaxLength + 1)},
			// パスワードのチェックデータ
			{Email: validEmail, Password: ""},
			{Email: validEmail, Password: strings.Repeat("a", passwordMinLength-1)},
			{Email: validEmail, Password: strings.Repeat("a", passwordMaxLength+1)},
			{Email: validEmail, Password: "にほんごぱすわーど"},
		}

		var userID *uint = nil
		asMock := new(authServiceMock)
		asMock.On("Authenticate", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(userID, nil)

		ac := NewAuthController(asMock)

		for _, in := range inputs {
			reqBody, err := testutils.CreateRequestBodyStr(in)
			if err != nil {
				s.FailNow(err.Error())
			}
			r, c := createLoginPostContext(reqBody)

			ac.Login(c)

			var res apiResponse
			json.Unmarshal(r.Body.Bytes(), &res)

			assert.Equal(s.T(), http.StatusUnauthorized, r.Code)
			assert.Equal(s.T(), "メールアドレスまたはパスワードが違います", res.Message)
		}
	})

	s.Run("正常な入力値の場合", func() {
		inputs := []forms.LoginForm{
			// メールアドレスのチェックデータ
			{Password: validPassword, Email: createEmailAddress(emailMaxLength)},
			{Password: validPassword, Email: "にほんご@example.com"},
			// パスワードのチェックデータ
			{Email: validEmail, Password: strings.Repeat("a", passwordMinLength)},
			{Email: validEmail, Password: strings.Repeat("a", passwordMaxLength)},
		}

		// AuthServiceモック
		userID := uint(1)
		asMock := new(authServiceMock)
		asMock.On("Authenticate", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(&userID, nil)

		ac := NewAuthController(asMock)

		for _, in := range inputs {
			reqBody, err := testutils.CreateRequestBodyStr(in)
			if err != nil {
				s.FailNow(err.Error())
			}
			r, c := createLoginPostContext(reqBody)

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
		asMock := new(authServiceMock)
		asMock.On("Authenticate", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(userID, nil)

		ac := NewAuthController(asMock)

		reqBody, err := testutils.CreateRequestBodyStr(forms.LoginForm{Email: validEmail, Password: validPassword})
		if err != nil {
			s.FailNow(err.Error())
		}
		r, c := createLoginPostContext(reqBody)

		ac.Login(c)

		var res apiResponse
		json.Unmarshal(r.Body.Bytes(), &res)

		assert.Equal(s.T(), http.StatusUnauthorized, r.Code)
		assert.Equal(s.T(), "メールアドレスまたはパスワードが違います", res.Message)
	})

	s.Run("予期せぬエラーが発生した場合", func() {
		var userID *uint = nil
		asMock := new(authServiceMock)
		asMock.On("Authenticate", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(userID, testutils.ErrTest)

		ac := NewAuthController(asMock)

		reqBody, err := testutils.CreateRequestBodyStr(forms.LoginForm{Email: validEmail, Password: validPassword})
		if err != nil {
			s.FailNow(err.Error())
		}
		r, c := createLoginPostContext(reqBody)

		ac.Login(c)

		var res apiResponse
		json.Unmarshal(r.Body.Bytes(), &res)

		assert.Equal(s.T(), r.Code, http.StatusInternalServerError)
		assert.Equal(s.T(), "予期せぬエラーが発生しました", res.Message)
	})
}

func (s *AuthControllerTestSuite) TestLogout() {
	ac := NewAuthController(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/logout", nil)

	// sessions.sessionモック
	sessMock := mocks.NewSessionMock()
	testutils.SetSessionMockToGin(c, sessMock)

	ac.Logout(c)

	var res apiResponse
	json.Unmarshal(w.Body.Bytes(), &res)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	assert.Equal(s.T(), "logged out", res.Message)
	sessMock.AssertCalled(s.T(), "Save")
	sessMock.AssertCalled(s.T(), "Clear")
}

// 指定の長さのメールアドレスを生成する
func createEmailAddress(length int) string {
	return strings.Repeat("a", length-len("@example.com")) + "@example.com"
}

// ログインリクエストのGinコンテキストを生成する
func createLoginPostContext(reqBody string) (*httptest.ResponseRecorder, *gin.Context) {
	r := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(r)
	c.Request, _ = http.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	return r, c
}
