package server

import (
	"errors"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/gomiboko/my-circle/controllers"
	"github.com/gomiboko/my-circle/db"
	"github.com/gomiboko/my-circle/repositories"
	"github.com/gomiboko/my-circle/services"
	"github.com/gomiboko/my-circle/validations"
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
	// (内部で許可されていないOriginの場合は処理を中断して403を返しているので、OriginチェックによるCSRF対策も兼ねる)
	cfg := cors.DefaultConfig()
	cfg.AllowOrigins = []string{os.Getenv("FRONTEND_ORIGIN")}
	cfg.AllowCredentials = true
	r.Use(cors.New(cfg))

	// カスタムバリデーションの登録
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("notonlywhitespace", validations.NotOnlyWhitespace)
		v.RegisterValidation("password", validations.Password)
	} else {
		return nil, errors.New("カスタムバリデーションの登録に失敗しました")
	}

	// ルーティング
	ur := repositories.NewUserRepository(db.GetDB())
	ac := controllers.NewAuthController(services.NewAuthService(ur))
	uc := controllers.NewUserController(services.NewUserService(ur))
	r.POST("/login", ac.Login)
	r.GET("/logout", ac.Logout)
	auth := r.Group("/auth")
	{
		auth.GET("/check", ac.IsAuthorized)
	}
	users := r.Group("/users")
	{
		users.POST("", uc.Create)
		users.GET("/me", uc.GetHomeInfo)
	}

	return r, nil
}
