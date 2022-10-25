package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/gomiboko/my-circle/consts"
	"github.com/gomiboko/my-circle/controllers/mocks"
	"github.com/gomiboko/my-circle/forms"
	"github.com/gomiboko/my-circle/models"
	"github.com/gomiboko/my-circle/responses"
	"github.com/gomiboko/my-circle/testutils"
	"github.com/gomiboko/my-circle/validations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type homeInfoResponse struct {
	UserName    string
	UserIconUrl string
	Circles     []responses.Circle
}

type UserControllerTestSuite struct {
	suite.Suite
}

func (s *UserControllerTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	// カスタムバリデータの設定
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("notonlywhitespace", validations.NotOnlyWhitespace)
		v.RegisterValidation("password", validations.Password)
	}
}

func TestUserController(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}

func (s *UserControllerTestSuite) TestCreateUser() {
	const reqPath = "/users"

	s.Run("不正なリクエスト(URLエンコード)の場合", func() {
		usMock := new(mocks.UserServiceMock)
		usMock.On("CreateUser", mock.AnythingOfType("forms.UserForm")).Return(&models.User{}, testutils.ErrTest)

		uc := NewUserController(usMock)

		values := url.Values{}
		values.Set("username", testutils.ValidUserName)
		values.Set("email", testutils.ValidEmail)
		values.Set("password", testutils.ValidPassword)

		// URLエンコードで送信
		r := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(r)
		c.Request, _ = http.NewRequest(http.MethodPost, "/users", strings.NewReader(values.Encode()))
		c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		uc.Create(c)

		var res testutils.ApiErrorReponse
		json.Unmarshal(r.Body.Bytes(), &res)

		assert.Equal(s.T(), http.StatusBadRequest, r.Code)
		assert.Equal(s.T(), "不正なリクエストです", res.Message)
	})

	s.Run("不正な入力値の場合", func() {
		inputs := []forms.UserForm{
			// ユーザ名のチェックデータ
			{Email: testutils.ValidEmail, Password: testutils.ValidPassword, Username: ""},
			{Email: testutils.ValidEmail, Password: testutils.ValidPassword, Username: testutils.HalfWidthSpace},
			{Email: testutils.ValidEmail, Password: testutils.ValidPassword, Username: testutils.FullWidthSpace},
			{Email: testutils.ValidEmail, Password: testutils.ValidPassword, Username: strings.Repeat("a", testutils.UsernameMaxLength+1)},
			// メールアドレスのチェックデータ
			{Username: testutils.ValidUserName, Password: testutils.ValidPassword, Email: ""},
			{Username: testutils.ValidUserName, Password: testutils.ValidPassword, Email: testutils.HalfWidthSpace},
			{Username: testutils.ValidUserName, Password: testutils.ValidPassword, Email: testutils.FullWidthSpace},
			{Username: testutils.ValidUserName, Password: testutils.ValidPassword, Email: "a"},
			{Username: testutils.ValidUserName, Password: testutils.ValidPassword, Email: testutils.CreateEmailAddress(testutils.EmailMaxLength + 1)},
			// パスワードのチェックデータ
			{Username: testutils.ValidUserName, Email: testutils.ValidEmail, Password: ""},
			{Username: testutils.ValidUserName, Email: testutils.ValidEmail, Password: strings.Repeat("a", testutils.PasswordMinLength-1)},
			{Username: testutils.ValidUserName, Email: testutils.ValidEmail, Password: strings.Repeat("a", testutils.PasswordMaxLength+1)},
			{Username: testutils.ValidUserName, Email: testutils.ValidEmail, Password: testutils.ValidPassword + testutils.HalfWidthSpace},
			{Username: testutils.ValidUserName, Email: testutils.ValidEmail, Password: testutils.ValidPassword + testutils.FullWidthSpace},
			{Username: testutils.ValidUserName, Email: testutils.ValidEmail, Password: testutils.ValidPassword + testutils.FullWidthA},
			{Username: testutils.ValidUserName, Email: testutils.ValidEmail, Password: testutils.ValidPassword + "😋"},
			{Username: testutils.ValidUserName, Email: testutils.ValidEmail, Password: testutils.ValidPassword + "あ"},
		}

		usMock := new(mocks.UserServiceMock)
		usMock.On("CreateUser", mock.AnythingOfType("forms.UserForm")).Return(&models.User{}, testutils.ErrTest)

		uc := NewUserController(usMock)

		for _, in := range inputs {
			reqBody, err := testutils.CreateRequestBodyStr(in)
			if err != nil {
				s.FailNow(err.Error())
			}
			r, c := testutils.CreatePostContext(reqPath, reqBody)

			uc.Create(c)
			c.Writer.WriteHeaderNow()

			var res testutils.ApiErrorReponse
			json.Unmarshal(r.Body.Bytes(), &res)

			assert.Equal(s.T(), http.StatusBadRequest, r.Code)
			assert.Equal(s.T(), "不正なリクエストです", res.Message)
		}
	})

	s.Run("正常な入力値の場合", func() {
		inputs := []forms.UserForm{
			// ユーザ名のチェックデータ
			{Email: testutils.ValidEmail, Password: testutils.ValidPassword, Username: strings.Repeat("a", 1)},
			{Email: testutils.ValidEmail, Password: testutils.ValidPassword, Username: strings.Repeat("a", testutils.UsernameMaxLength)},
			// メールアドレスのチェックデータ
			{Username: testutils.ValidUserName, Password: testutils.ValidPassword, Email: testutils.CreateEmailAddress(testutils.EmailMaxLength)},
			// パスワードのチェックデータ
			{Username: testutils.ValidUserName, Email: testutils.ValidEmail, Password: strings.Repeat("a", testutils.PasswordMinLength)},
			{Username: testutils.ValidUserName, Email: testutils.ValidEmail, Password: strings.Repeat("a", testutils.PasswordMaxLength)},
			{Username: testutils.ValidUserName, Email: testutils.ValidEmail, Password: testutils.HalfWidthSymbol},
			{Username: testutils.ValidUserName, Email: testutils.ValidEmail, Password: "1234567890"},
			{Username: testutils.ValidUserName, Email: testutils.ValidEmail, Password: "abcdefghijklmnopqrstuvwxyz"},
		}

		for _, in := range inputs {
			user := models.User{
				ID:    1,
				Name:  in.Username,
				Email: in.Email,
			}

			usMock := new(mocks.UserServiceMock)
			usMock.On("CreateUser", mock.AnythingOfType("forms.UserForm")).Return(&user, nil)

			uc := NewUserController(usMock)

			reqBody, err := testutils.CreateRequestBodyStr(in)
			if err != nil {
				s.FailNow(err.Error())
			}
			r, c := testutils.CreatePostContext(reqPath, reqBody)

			// sessions.sessionモック
			sessMock := mocks.NewSessionMock()
			testutils.SetSessionMockToGin(c, sessMock)

			uc.Create(c)
			c.Writer.WriteHeaderNow()

			assert.Equal(s.T(), http.StatusCreated, r.Code)
			assert.Equal(s.T(), r.Body.Len(), 0)
			sessMock.AssertCalled(s.T(), "Set", mock.AnythingOfType("string"), mock.AnythingOfType("uint"))
			sessMock.AssertCalled(s.T(), "Save")
		}
	})

	s.Run("登録済みのメールアドレスだった場合", func() {
		user := models.User{Name: testutils.ValidUserName, Email: testutils.User1Email}

		usMock := new(mocks.UserServiceMock)
		usMock.On("CreateUser", mock.AnythingOfType("forms.UserForm")).Return(&user, testutils.ErrDuplicateEntry)

		uc := NewUserController(usMock)

		form := forms.UserForm{Username: testutils.ValidUserName, Email: testutils.ValidEmail, Password: testutils.ValidPassword}
		reqBody, err := testutils.CreateRequestBodyStr(form)
		if err != nil {
			s.FailNow(err.Error())
		}
		r, c := testutils.CreatePostContext(reqPath, reqBody)

		// sessions.sessionモック
		sessMock := mocks.NewSessionMock()
		testutils.SetSessionMockToGin(c, sessMock)

		uc.Create(c)
		c.Writer.WriteHeaderNow()

		var res testutils.ApiErrorReponse
		json.Unmarshal(r.Body.Bytes(), &res)

		assert.Equal(s.T(), http.StatusConflict, r.Code)
		assert.Equal(s.T(), "登録済みのメールアドレスです", res.Message)
		sessMock.AssertNotCalled(s.T(), "Set", mock.AnythingOfType("string"), mock.AnythingOfType("uint"))
		sessMock.AssertNotCalled(s.T(), "Save")
	})

	s.Run("予期せぬエラーが発生した場合", func() {
		user := models.User{Name: testutils.ValidUserName, Email: testutils.UnregisteredEmail}

		usMock := new(mocks.UserServiceMock)
		usMock.On("CreateUser", mock.AnythingOfType("forms.UserForm")).Return(&user, testutils.ErrTest)

		uc := NewUserController(usMock)

		form := forms.UserForm{Username: testutils.ValidUserName, Email: testutils.ValidEmail, Password: testutils.ValidPassword}
		reqBody, err := testutils.CreateRequestBodyStr(form)
		if err != nil {
			s.FailNow(err.Error())
		}
		r, c := testutils.CreatePostContext(reqPath, reqBody)

		// sessions.sessionモック
		sessMock := mocks.NewSessionMock()
		testutils.SetSessionMockToGin(c, sessMock)

		uc.Create(c)
		c.Writer.WriteHeaderNow()

		var res testutils.ApiErrorReponse
		json.Unmarshal(r.Body.Bytes(), &res)

		assert.Equal(s.T(), http.StatusInternalServerError, r.Code)
		assert.Equal(s.T(), testutils.UnexpectedErrMsg, res.Message)
		sessMock.AssertNotCalled(s.T(), "Set", mock.AnythingOfType("string"), mock.AnythingOfType("uint"))
		sessMock.AssertNotCalled(s.T(), "Save")
	})
}

func (s *UserControllerTestSuite) TestGetHomeInfo() {
	s.Run("ログイン中のユーザ情報が取得できた場合", func() {
		circles := []responses.Circle{{ID: 1}}
		homeInfo := gin.H{
			consts.ResKeyUserName:    testutils.ValidUserName,
			consts.ResKeyUserIconUrl: testutils.ValidUrl,
			consts.ResKeyCircles:     circles,
		}
		usMock := new(mocks.UserServiceMock)
		usMock.On("GetHomeInfo", mock.AnythingOfType("uint")).Return(homeInfo, nil)

		uc := NewUserController(usMock)

		r, c := testutils.CreateGetContext("/users/me")

		sessMock := mocks.NewSessionMock()
		sessMock.On("Get", mock.AnythingOfType("string")).Return(interface{}(uint(2)))
		testutils.SetSessionMockToGin(c, sessMock)

		uc.GetHomeInfo(c)
		c.Writer.WriteHeaderNow()

		var res homeInfoResponse
		json.Unmarshal(r.Body.Bytes(), &res)

		assert.Equal(s.T(), http.StatusOK, r.Code)
		assert.Equal(s.T(), testutils.ValidUserName, res.UserName)
		assert.Equal(s.T(), testutils.ValidUrl, res.UserIconUrl)
		assert.Equal(s.T(), 1, len(res.Circles))
		assert.Equal(s.T(), uint(1), res.Circles[0].ID)
	})

	s.Run("予期せぬエラーが発生した場合", func() {
		usMock := new(mocks.UserServiceMock)
		usMock.On("GetHomeInfo", mock.AnythingOfType("uint")).Return(nil, testutils.ErrTest)

		uc := NewUserController(usMock)

		r, c := testutils.CreateGetContext("/users/me")

		sessMock := mocks.NewSessionMock()
		sessMock.On("Get", mock.AnythingOfType("string")).Return(interface{}(uint(1)))
		testutils.SetSessionMockToGin(c, sessMock)

		uc.GetHomeInfo(c)
		c.Writer.WriteHeaderNow()

		var res testutils.ApiErrorReponse
		json.Unmarshal(r.Body.Bytes(), &res)

		assert.Equal(s.T(), http.StatusInternalServerError, r.Code)
		assert.Equal(s.T(), "予期せぬエラーが発生しました", res.Message)
	})
}
