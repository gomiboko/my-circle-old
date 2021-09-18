package server

import (
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/gomiboko/my-circle/controllers"
	"github.com/gomiboko/my-circle/db"
	"github.com/gomiboko/my-circle/repositories"
	"github.com/gomiboko/my-circle/services"
)

func NewRouter() (*gin.Engine, error) {
	r := gin.Default()

	isSecure, err := strconv.ParseBool(os.Getenv("COOKIE_SECURE"))
	if err != nil {
		return nil, err
	}

	// セッションの設定
	store := memstore.NewStore([]byte(os.Getenv("SESSION_AUTH_KEY")))
	store.Options(sessions.Options{
		Secure:   isSecure,
		HttpOnly: true,
		MaxAge:   60 * 60 * 24 * 30,
	})
	r.Use(sessions.Sessions("sessionid", store))

	// CORSの設定
	cfg := cors.DefaultConfig()
	cfg.AllowOrigins = []string{os.Getenv("FRONTEND_ORIGIN")}
	cfg.AllowCredentials = true
	r.Use(cors.New(cfg))

	// ルーティング
	ur := repositories.NewUserRepository(db.GetDB())
	ac := controllers.NewAuthController(services.NewAuthService(ur))
	r.POST("/login", ac.Login)
	r.GET("/logout", ac.Logout)

	return r, nil
}
