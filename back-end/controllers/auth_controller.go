package controllers

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gomiboko/my-circle/forms"
	"github.com/gomiboko/my-circle/services"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(as services.AuthService) *AuthController {
	ac := &AuthController{
		authService: as,
	}
	return ac
}

func (ac AuthController) Login(c *gin.Context) {
	var form forms.LoginForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(responseBody400BadRequest())
		return
	}

	// ログイン認証
	userID, err := ac.authService.Authenticate(form.Email, form.Password)

	if err != nil {
		log.Print(err)

		c.JSON(responseBody500UnexpectedError())
		return
	}

	// 認証失敗
	if userID == nil {
		c.JSON(http.StatusUnauthorized, MessageResponseBody("メールアドレスまたはパスワードが違います"))
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", *userID)
	session.Save()

	c.Status(http.StatusCreated)
}

func (ac AuthController) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Status(http.StatusOK)
}

func (ac AuthController) IsAuthorized(c *gin.Context) {
	session := sessions.Default(c)
	isLoggedIn := session.Get("user_id") != nil

	c.JSON(http.StatusOK, gin.H{"isAuthorized": isLoggedIn})
}
