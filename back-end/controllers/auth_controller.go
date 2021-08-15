package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gomiboko/my-circle/repositories"
)

type AuthController struct {
	ur repositories.UserRepository
}

func NewAuthController(ur repositories.UserRepository) *AuthController {
	ac := &AuthController{
		ur: ur,
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
	user, err := ac.ur.GetUser(form.Email, form.Password)

	if err != nil {
		log.Print(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "予期せぬエラー",
		})
		return
	}

	// 認証失敗
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "認証エラー",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "logged in",
		"user": user,
	})
}

func (ac AuthController) Logout(c *gin.Context) {
	// TODO:
	c.JSON(http.StatusOK, gin.H{
		"msg": "logged out",
	})
}
