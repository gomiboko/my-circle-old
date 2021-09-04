package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gomiboko/my-circle/services"
)

type AuthController struct {
	as services.AuthService
}

func NewAuthController(as services.AuthService) *AuthController {
	ac := &AuthController{
		as: as,
	}
	return ac
}

type LoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (ac AuthController) Login(c *gin.Context) {
	var form LoginForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(responseBody400BadRequest())
		return
	}

	// ログイン認証
	result, err := ac.as.Authenticate(form.Email, form.Password)

	if err != nil {
		log.Print(err)

		c.JSON(responseBody500UnexpectedError())
		return
	}

	// 認証失敗
	if !result {
		c.JSON(http.StatusUnauthorized, messageResponseBody("メールアドレスまたはパスワードが違います"))
		return
	}

	c.Status(http.StatusCreated)
}

func (ac AuthController) Logout(c *gin.Context) {
	// TODO:
	c.JSON(http.StatusOK, messageResponseBody("logged out"))
}
