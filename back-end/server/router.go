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
	"github.com/gomiboko/my-circle/aws"
	"github.com/gomiboko/my-circle/consts"
	"github.com/gomiboko/my-circle/controllers"
	"github.com/gomiboko/my-circle/db"
	"github.com/gomiboko/my-circle/middlewares"
	"github.com/gomiboko/my-circle/repositories"
	"github.com/gomiboko/my-circle/services"
	"github.com/gomiboko/my-circle/validations"
)

const (
	pathV1    = "/v1"
	pathUsers = "/users"
)

func NewRouter() (*gin.Engine, error) {
	r := gin.Default()

	// ミドルウェアの設定
	err := setupMiddlewares(r)
	if err != nil {
		return nil, err
	}

	// カスタムバリデーションの設定
	err = setupCustomValidations()
	if err != nil {
		return nil, err
	}

	// ルーティングの設定
	setupRoutings(r)

	return r, nil
}

func setupMiddlewares(r *gin.Engine) error {
	isSecure, err := strconv.ParseBool(os.Getenv(consts.EnvCookieSecure))
	if err != nil {
		return err
	}

	// セッションの設定
	store := memstore.NewStore([]byte(os.Getenv(consts.EnvSessionAuthKey)))
	store.Options(sessions.Options{
		Secure:   isSecure,
		HttpOnly: true,
		MaxAge:   60 * 60 * 24 * 30,
	})
	r.Use(sessions.Sessions("sessionid", store))

	// CORSの設定
	// 内部で許可されていないOriginの場合は処理を中断して403を返しているので、OriginチェックによるCSRF対策も兼ねる。
	// Abortするミドルウェアを登録する場合、CORSミドルウェアより先に登録してしまうと
	// 本来CORSミドルウェアで設定されるはずだった「Access-Control-Allow-Origin」ヘッダが付与されずにレスポンスが返却されることになり、
	// クライアント側でCORSエラーが発生してしまうので注意。
	cfg := cors.DefaultConfig()
	cfg.AllowOrigins = []string{os.Getenv(consts.EnvFrontendOrigin)}
	cfg.AllowCredentials = true
	r.Use(cors.New(cfg))

	// リクエストボディサイズ上限の設定
	r.Use(middlewares.RequestBodySizeLimiter())

	return nil
}

func setupCustomValidations() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("notonlywhitespace", validations.NotOnlyWhitespace)
		v.RegisterValidation("password", validations.Password)
		return nil
	} else {
		return errors.New(consts.MsgFailedToRegisterValidations)
	}
}

func setupRoutings(r *gin.Engine) {
	storageService := services.NewS3Service(aws.GetConf())

	ur := repositories.NewUserRepository(db.GetDB())
	cr := repositories.NewCircleRepository(db.GetDB())
	ucr := repositories.NewUsersCirclesRepository(db.GetDB())

	sc := controllers.NewSessionController(services.NewSessionService(ur))
	uc := controllers.NewUserController(services.NewUserService(ur))
	cc := controllers.NewCircleController(services.NewCircleService(cr, ucr), storageService)

	// 認証が不要なエンドポイント
	v1 := r.Group(pathV1)
	sess := v1.Group("/sessions")
	{
		sess.POST("", sc.Create)
		sess.DELETE("", sc.Destroy)
	}
	v1.POST(pathUsers, uc.Create)

	// 認証が必要なエンドポイント
	authorized := r.Group("/", middlewares.AuthRequired())
	v1Auth := authorized.Group(pathV1)
	{
		users := v1Auth.Group(pathUsers)
		users.GET("/me", uc.GetHomeInfo)

		circles := v1Auth.Group("/circles")
		circles.POST("", cc.Create)
	}
}
