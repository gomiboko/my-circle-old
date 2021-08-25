package server

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gomiboko/my-circle/controllers"
	"github.com/gomiboko/my-circle/db"
	"github.com/gomiboko/my-circle/repositories"
	"github.com/gomiboko/my-circle/services"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	cfg := cors.DefaultConfig()
	cfg.AllowOrigins = []string{os.Getenv("FRONTEND_ORIGIN")}
	r.Use(cors.New(cfg))

	ur := repositories.NewUserRepository(db.GetDB())
	ac := controllers.NewAuthController(services.NewAuthService(ur))
	r.POST("/login", ac.Login)
	r.GET("/logout", ac.Logout)

	return r
}
