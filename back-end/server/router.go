package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gomiboko/my-circle/controllers"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	a := new(controllers.Auth)
	r.GET("/login", a.Login)
	r.GET("/logout", a.Logout)

	return r
}
