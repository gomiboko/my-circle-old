package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Auth struct{}

type LoginForm struct {
	Email    string `json:"email"    binding:"required,email,max=254"`
	Password string `json:"password" binding:"required,alphanum,min=8,max=64"`
}

func (a Auth) Login(c *gin.Context) {
	var form LoginForm
	err := c.BindJSON(&form)

	if err != nil {
		return
	}

	// TODO: ログイン処理

	c.JSON(http.StatusOK, gin.H{
		"msg": "logged in",
	})
}

func (a Auth) Logout(c *gin.Context) {
	// TODO:
	c.JSON(http.StatusOK, gin.H{
		"msg": "logged out",
	})
}
