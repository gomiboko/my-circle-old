package server

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gomiboko/my-circle/controllers"
	"github.com/gomiboko/my-circle/repositories"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	cfg := cors.DefaultConfig()
	cfg.AllowOrigins = []string{os.Getenv("FRONTEND_ORIGIN")}
	r.Use(cors.New(cfg))

	ac := controllers.NewAuthController(repositories.NewUserRepository())
	r.POST("/login", ac.Login)
	r.GET("/logout", ac.Logout)

	return r
}
