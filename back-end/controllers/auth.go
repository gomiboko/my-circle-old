package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Auth struct{}

func (a Auth) Login(c *gin.Context) {
	// TODO:
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
