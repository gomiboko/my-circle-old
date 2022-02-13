package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-contrib/sessions"
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
		uc := NewUserController(nil)

		values := url.Values{}
		values.Set("username", "test name")
		values.Set("email", validEmail)
		values.Set("password", validPassword)

		// URLエンコードで送信
		r := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(r)
		c.Request, _ = http.NewRequest(http.MethodPost, "/users", strings.NewReader(values.Encode()))
		c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		uc.Create(c)

		var res apiResponse
		json.Unmarshal(r.Body.Bytes(), &res)

		assert.Equal(s.T(), http.StatusBadRequest, r.Code)
		assert.Equal(s.T(), "不正なリクエストです", res.Message)
	})

	s.Run("不正な入力値の場合", func() {
		inputs := []forms.UserForm{
			// ユーザ名のチェックデータ
			{Email: validEmail, Password: validPassword, Username: ""},
			{Email: validEmail, Password: validPassword, Username: strings.Repeat("a", 46)},
			// メールアドレスのチェックデータ
			{Username: "test name", Password: validPassword, Email: ""},
			{Username: "test name", Password: validPassword, Email: "a"},
			{Username: "test name", Password: validPassword, Email: "😋@example.com"},
			{Username: "test name", Password: validPassword, Email: "@example.com"},
			// パスワードのチェックデータ
			{Username: "test name", Email: validEmail, Password: ""},
			{Username: "test name", Email: validEmail, Password: strings.Repeat("a", 7)},
			{Username: "test name", Email: validEmail, Password: strings.Repeat("a", 129)},
			{Username: "test name", Email: validEmail, Password: strings.Repeat("a", 7) + testutils.HalfWidthSpace},
			{Username: "test name", Email: validEmail, Password: strings.Repeat("a", 7) + testutils.FullWidthSpace},
			{Username: "test name", Email: validEmail, Password: strings.Repeat("a", 7) + "😋"},
			{Username: "test name", Email: validEmail, Password: strings.Repeat("a", 7) + "あ"},
		}

		uc := NewUserController(nil)

		for _, in := range inputs {
			reqBody, err := createRequestBodyStr(in)
			if err != nil {
				s.FailNow(err.Error())
			}
			r, c := createUserPostContext(reqBody)

			uc.Create(c)
			c.Writer.WriteHeaderNow()

			var res apiResponse
			json.Unmarshal(r.Body.Bytes(), &res)

			assert.Equal(s.T(), http.StatusBadRequest, r.Code)
			assert.Equal(s.T(), "不正なリクエストです", res.Message)
		}
	})

	s.Run("正常な入力値の場合", func() {
		inputs := []forms.UserForm{
			// ユーザ名のチェックデータ
			{Email: validEmail, Password: validPassword, Username: strings.Repeat("a", 1)},
			{Email: validEmail, Password: validPassword, Username: strings.Repeat("a", 45)},
			{Email: validEmail, Password: validPassword, Username: "にほんご"},
			{Email: validEmail, Password: validPassword, Username: "😋"},
			// メールアドレスのチェックデータ
			{Username: "test name", Password: validPassword, Email: createEmailAddress(emailMaxLength)},
			{Username: "test name", Password: validPassword, Email: "あ@example.com"},
			// パスワードのチェックデータ
			{Username: "test name", Email: validEmail, Password: strings.Repeat("a", 8)},
			{Username: "test name", Email: validEmail, Password: strings.Repeat("a", 128)},
			{Username: "test name", Email: validEmail, Password: "!@#$%^&*()-_=+[]{}\\|~;:'\",.<>/?`"},
			{Username: "test name", Email: validEmail, Password: "1234567890"},
			{Username: "test name", Email: validEmail, Password: "abcdefghijklmnopqrstuvwxyz"},
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

			reqBody, err := createRequestBodyStr(in)
			if err != nil {
				s.FailNow(err.Error())
			}
			r, c := createUserPostContext(reqBody)

			// sessions.sessionモック
			sessMock := new(mocks.SessionMock)
			sessMock.On("Set", mock.AnythingOfType("string"), mock.Anything)
			sessMock.On("Save").Return(nil)
			// sessions.Sessions(string, sessions.Store) と同様の処理を実行
			c.Set(sessions.DefaultKey, sessMock)

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

		form := forms.UserForm{Username: "test name", Email: validEmail, Password: validPassword}
		reqBody, err := createRequestBodyStr(form)
		if err != nil {
			s.FailNow(err.Error())
		}
		r, c := createUserPostContext(reqBody)

		// sessions.sessionモック
		sessMock := new(mocks.SessionMock)
		sessMock.On("Set", mock.AnythingOfType("string"), mock.Anything)
		sessMock.On("Save").Return(nil)
		// sessions.Sessions(string, sessions.Store) と同様の処理を実行
		c.Set(sessions.DefaultKey, sessMock)

		uc.Create(c)
		c.Writer.WriteHeaderNow()

		var res apiResponse
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

		form := forms.UserForm{Username: "test name", Email: validEmail, Password: validPassword}
		reqBody, err := createRequestBodyStr(form)
		if err != nil {
			s.FailNow(err.Error())
		}
		r, c := createUserPostContext(reqBody)

		// sessions.sessionモック
		sessMock := new(mocks.SessionMock)
		sessMock.On("Set", mock.AnythingOfType("string"), mock.Anything)
		sessMock.On("Save").Return(nil)
		// sessions.Sessions(string, sessions.Store) と同様の処理を実行
		c.Set(sessions.DefaultKey, sessMock)

		uc.Create(c)
		c.Writer.WriteHeaderNow()

		var res apiResponse
		json.Unmarshal(r.Body.Bytes(), &res)

		assert.Equal(s.T(), http.StatusInternalServerError, r.Code)
		assert.Equal(s.T(), "予期せぬエラーが発生しました", res.Message)
		sessMock.AssertNotCalled(s.T(), "Set", mock.AnythingOfType("string"), mock.AnythingOfType("uint"))
		sessMock.AssertNotCalled(s.T(), "Save")
	})
}

func createRequestBodyStr(obj interface{}) (string, error) {
	if j, err := json.Marshal(obj); err != nil {
		return "", err
	} else {
		return string(j), nil
	}
}

func createUserPostContext(reqBody string) (*httptest.ResponseRecorder, *gin.Context) {
	r := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(r)
	c.Request, _ = http.NewRequest(http.MethodPost, "/users", strings.NewReader(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	return r, c
}
