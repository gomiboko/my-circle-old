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
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "不正なリクエスト",
		})
		return
	}

	// ログイン認証
	result, err := ac.as.Authenticate(form.Email, form.Password)

	if err != nil {
		log.Print(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "予期せぬエラー",
		})
		return
	}

	// 認証失敗
	if !result {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "認証エラー",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "logged in",
	})
}

func (ac AuthController) Logout(c *gin.Context) {
	// TODO:
	c.JSON(http.StatusOK, gin.H{
		"msg": "logged out",
	})
}
