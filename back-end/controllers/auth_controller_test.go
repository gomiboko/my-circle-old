package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const validEmail = "foo@example.com"
const validPassword = "password123"
const emailMaxLength = 254
const passwordMinLength = 8
const passwordMaxLength = 64

type apiResponse struct {
	Msg string
}

type authServiceMock struct {
	mock.Mock
}

func (m *authServiceMock) Authenticate(email string, password string) (bool, error) {
	args := m.Called(email, password)
	return args.Bool(0), args.Error(1)
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}

func TestInvalidRequest(t *testing.T) {
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

	assert.Equal(t, http.StatusBadRequest, r.Code)
	assert.Equal(t, "不正なリクエスト", res.Msg)
}

func TestLoginWithInvalidValue(t *testing.T) {
	inputs := []struct {
		email    string
		password string
	}{
		// メールアドレスのチェックデータ
		{
			email:    "",
			password: validPassword,
		},
		{
			email:    "isNotEmail",
			password: validPassword,
		},
		{
			email:    createEmailAddress(emailMaxLength + 1),
			password: validPassword,
		},
		// パスワードのチェックデータ
		{
			email:    validEmail,
			password: "",
		},
		{
			email:    validEmail,
			password: strings.Repeat("a", passwordMinLength-1),
		},
		{
			email:    validEmail,
			password: strings.Repeat("a", passwordMaxLength+1),
		},
		{
			email:    validEmail,
			password: "にほんごぱすわーど",
		},
	}

	asMock := new(authServiceMock)
	asMock.On("Authenticate", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(false, nil)

	ac := NewAuthController(asMock)

	for _, in := range inputs {
		reqBody := createRequestBody(in.email, in.password)
		r, c := createLoginPostContext(reqBody)

		ac.Login(c)

		var res apiResponse
		json.Unmarshal(r.Body.Bytes(), &res)

		assert.Equal(t, http.StatusUnauthorized, r.Code)
		assert.Equal(t, "認証エラー", res.Msg)
	}
}

func TestLoginWithValidValue(t *testing.T) {
	inputs := []struct {
		email    string
		password string
	}{
		// メールアドレスのチェックデータ
		{
			email:    createEmailAddress(emailMaxLength),
			password: validPassword,
		},
		{
			email:    "にほんご@example.com",
			password: validPassword,
		},
		// パスワードのチェックデータ
		{
			email:    validEmail,
			password: strings.Repeat("a", passwordMinLength),
		},
		{
			email:    validEmail,
			password: strings.Repeat("a", passwordMaxLength),
		},
	}

	asMock := new(authServiceMock)
	asMock.On("Authenticate", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(true, nil)

	ac := NewAuthController(asMock)

	for _, in := range inputs {
		reqBody := createRequestBody(in.email, in.password)
		r, c := createLoginPostContext(reqBody)

		ac.Login(c)

		var res apiResponse
		json.Unmarshal(r.Body.Bytes(), &res)

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Equal(t, "logged in", res.Msg)
	}
}

func TestLoginAuthenticationFailed(t *testing.T) {
	asMock := new(authServiceMock)
	asMock.On("Authenticate", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(false, nil)

	ac := NewAuthController(asMock)

	reqBody := createRequestBody(validEmail, validPassword)
	r, c := createLoginPostContext(reqBody)

	ac.Login(c)

	var res apiResponse
	json.Unmarshal(r.Body.Bytes(), &res)

	assert.Equal(t, http.StatusUnauthorized, r.Code)
	assert.Equal(t, "認証エラー", res.Msg)
}

func TestLoginUnexpectedError(t *testing.T) {
	asMock := new(authServiceMock)
	asMock.On("Authenticate", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(false, errors.New("test exception"))

	ac := NewAuthController(asMock)

	reqBody := createRequestBody(validEmail, validPassword)
	r, c := createLoginPostContext(reqBody)

	ac.Login(c)

	var res apiResponse
	json.Unmarshal(r.Body.Bytes(), &res)

	assert.Equal(t, r.Code, http.StatusInternalServerError)
	assert.Equal(t, "予期せぬエラー", res.Msg)
}

func TestLogout(t *testing.T) {
	ac := NewAuthController(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/logout", nil)

	ac.Logout(c)

	var res apiResponse
	json.Unmarshal(w.Body.Bytes(), &res)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "logged out", res.Msg)
}

// 指定の長さのメールアドレスを生成する
func createEmailAddress(length int) string {
	return strings.Repeat("a", length-len("@example.com")) + "@example.com"
}

// ログインのリクエストボディ文字列を生成する
func createRequestBody(email string, password string) string {
	return `{"email":"` + email + `","password":"` + password + `"}`
}

// ログインリクエストのGinコンテキストを生成する
func createLoginPostContext(reqBody string) (*httptest.ResponseRecorder, *gin.Context) {
	r := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(r)
	c.Request, _ = http.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	return r, c
}
