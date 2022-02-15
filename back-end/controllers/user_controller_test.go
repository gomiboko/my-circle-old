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
	"github.com/gomiboko/my-circle/controllers/mocks"
	"github.com/gomiboko/my-circle/forms"
	"github.com/gomiboko/my-circle/models"
	"github.com/gomiboko/my-circle/testutils"
	"github.com/gomiboko/my-circle/validations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type usersPostResponse struct {
	User models.User
}

type UserControllerTestSuite struct {
	suite.Suite
}

func (s *UserControllerTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	// カスタムバリデータの設定
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("password", validations.Password)
	}
}

func TestUserController(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}

func (s *UserControllerTestSuite) TestCreateUser() {
	s.Run("不正なリクエスト(URLエンコード)の場合", func() {
		usMock := new(mocks.UserServiceMock)
		usMock.On("CreateUser", mock.AnythingOfType("forms.UserForm")).Return(&models.User{}, testutils.ErrTest)

		uc := NewUserController(usMock)

		values := url.Values{}
		values.Set("username", "test name")
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
			{Email: testutils.ValidEmail, Password: testutils.ValidPassword, Username: strings.Repeat("a", testutils.UsernameMaxLength+1)},
			// メールアドレスのチェックデータ
			{Username: "test name", Password: testutils.ValidPassword, Email: ""},
			{Username: "test name", Password: testutils.ValidPassword, Email: "a"},
			{Username: "test name", Password: testutils.ValidPassword, Email: "😋@example.com"},
			{Username: "test name", Password: testutils.ValidPassword, Email: "@example.com"},
			// パスワードのチェックデータ
			{Username: "test name", Email: testutils.ValidEmail, Password: ""},
			{Username: "test name", Email: testutils.ValidEmail, Password: strings.Repeat("a", testutils.PasswordMinLength-1)},
			{Username: "test name", Email: testutils.ValidEmail, Password: strings.Repeat("a", testutils.PasswordMaxLength+1)},
			{Username: "test name", Email: testutils.ValidEmail, Password: strings.Repeat("a", testutils.PasswordMinLength-1) + testutils.HalfWidthSpace},
			{Username: "test name", Email: testutils.ValidEmail, Password: strings.Repeat("a", testutils.PasswordMinLength-1) + testutils.FullWidthSpace},
			{Username: "test name", Email: testutils.ValidEmail, Password: strings.Repeat("a", testutils.PasswordMinLength-1) + "😋"},
			{Username: "test name", Email: testutils.ValidEmail, Password: strings.Repeat("a", testutils.PasswordMinLength-1) + "あ"},
		}

		usMock := new(mocks.UserServiceMock)
		usMock.On("CreateUser", mock.AnythingOfType("forms.UserForm")).Return(&models.User{}, testutils.ErrTest)

		uc := NewUserController(usMock)

		for _, in := range inputs {
			reqBody, err := testutils.CreateRequestBodyStr(in)
			if err != nil {
				s.FailNow(err.Error())
			}
			r, c := createUserPostContext(reqBody)

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
			{Email: testutils.ValidEmail, Password: testutils.ValidPassword, Username: "にほんご"},
			{Email: testutils.ValidEmail, Password: testutils.ValidPassword, Username: "😋"},
			// メールアドレスのチェックデータ
			{Username: "test name", Password: testutils.ValidPassword, Email: testutils.CreateEmailAddress(testutils.EmailMaxLength)},
			{Username: "test name", Password: testutils.ValidPassword, Email: "あ@example.com"},
			// パスワードのチェックデータ
			{Username: "test name", Email: testutils.ValidEmail, Password: strings.Repeat("a", testutils.PasswordMinLength)},
			{Username: "test name", Email: testutils.ValidEmail, Password: strings.Repeat("a", testutils.PasswordMaxLength)},
			{Username: "test name", Email: testutils.ValidEmail, Password: "!@#$%^&*()-_=+[]{}\\|~;:'\",.<>/?`"},
			{Username: "test name", Email: testutils.ValidEmail, Password: "1234567890"},
			{Username: "test name", Email: testutils.ValidEmail, Password: "abcdefghijklmnopqrstuvwxyz"},
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
			r, c := createUserPostContext(reqBody)

			// sessions.sessionモック
			sessMock := mocks.NewSessionMock()
			testutils.SetSessionMockToGin(c, sessMock)

			uc.Create(c)
			c.Writer.WriteHeaderNow()

			var res usersPostResponse
			json.Unmarshal(r.Body.Bytes(), &res)

			assert.Equal(s.T(), http.StatusCreated, r.Code)
			assert.Equal(s.T(), uint(1), res.User.ID)
			assert.Equal(s.T(), in.Username, res.User.Name)
			assert.Equal(s.T(), in.Email, res.User.Email)
			sessMock.AssertCalled(s.T(), "Set", mock.AnythingOfType("string"), mock.AnythingOfType("uint"))
			sessMock.AssertCalled(s.T(), "Save")
		}
	})

	s.Run("登録済みのメールアドレスだった場合", func() {
		user := models.User{Name: "test name", Email: testutils.User1Email}

		usMock := new(mocks.UserServiceMock)
		usMock.On("CreateUser", mock.AnythingOfType("forms.UserForm")).Return(&user, testutils.ErrDuplicateEntry)

		uc := NewUserController(usMock)

		form := forms.UserForm{Username: "test name", Email: testutils.ValidEmail, Password: testutils.ValidPassword}
		reqBody, err := testutils.CreateRequestBodyStr(form)
		if err != nil {
			s.FailNow(err.Error())
		}
		r, c := createUserPostContext(reqBody)

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
		user := models.User{Name: "test name", Email: testutils.UnregisteredEmail}

		usMock := new(mocks.UserServiceMock)
		usMock.On("CreateUser", mock.AnythingOfType("forms.UserForm")).Return(&user, testutils.ErrTest)

		uc := NewUserController(usMock)

		form := forms.UserForm{Username: "test name", Email: testutils.ValidEmail, Password: testutils.ValidPassword}
		reqBody, err := testutils.CreateRequestBodyStr(form)
		if err != nil {
			s.FailNow(err.Error())
		}
		r, c := createUserPostContext(reqBody)

		// sessions.sessionモック
		sessMock := mocks.NewSessionMock()
		testutils.SetSessionMockToGin(c, sessMock)

		uc.Create(c)
		c.Writer.WriteHeaderNow()

		var res testutils.ApiErrorReponse
		json.Unmarshal(r.Body.Bytes(), &res)

		assert.Equal(s.T(), http.StatusInternalServerError, r.Code)
		assert.Equal(s.T(), "予期せぬエラーが発生しました", res.Message)
		sessMock.AssertNotCalled(s.T(), "Set", mock.AnythingOfType("string"), mock.AnythingOfType("uint"))
		sessMock.AssertNotCalled(s.T(), "Save")
	})
}

func createUserPostContext(reqBody string) (*httptest.ResponseRecorder, *gin.Context) {
	r := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(r)
	c.Request, _ = http.NewRequest(http.MethodPost, "/users", strings.NewReader(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	return r, c
}
