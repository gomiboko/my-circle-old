package server

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gomiboko/my-circle/controllers"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	cfg := cors.DefaultConfig()
	cfg.AllowOrigins = []string{os.Getenv("FRONTEND_ORIGIN")}
	r.Use(cors.New(cfg))

	a := new(controllers.Auth)
	r.GET("/login", a.Login)
	r.GET("/logout", a.Logout)

	return r
}
