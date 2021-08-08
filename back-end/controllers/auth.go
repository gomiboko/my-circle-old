package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gomiboko/my-circle/repositories"
)

type Auth struct {
	ur repositories.UserRepository
}

func NewAuth(ur repositories.UserRepository) *Auth {
	a := &Auth{
		ur: ur,
	}
	return a
}

type LoginForm struct {
	Email    string `json:"email"    binding:"required,email,max=254"`
	Password string `json:"password" binding:"required,alphanum,min=8,max=64"`
}

func (a Auth) Login(c *gin.Context) {
	var form LoginForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "入力エラー",
		})
		return
	}

	// ログイン認証
	user, err := a.ur.GetUser(form.Email, form.Password)

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

func (a Auth) Logout(c *gin.Context) {
	// TODO:
	c.JSON(http.StatusOK, gin.H{
		"msg": "logged out",
	})
}
