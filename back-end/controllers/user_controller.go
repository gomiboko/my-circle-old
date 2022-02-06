package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/gomiboko/my-circle/forms"
	"github.com/gomiboko/my-circle/services"
)

type UserController struct {
	us services.UserService
}

func NewUserController(us services.UserService) *UserController {
	uc := &UserController{
		us: us,
	}
	return uc
}

func (uc *UserController) Create(c *gin.Context) {
	// 入力チェック
	var form forms.UserForm
	if err := c.ShouldBindJSON(&form); err != nil {
		log.Print(err)
		c.JSON(responseBody400BadRequest())
		return
	}

	// ユーザ登録
	user, err := uc.us.CreateUser(form)
	if err != nil {
		log.Print(err)

		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			c.JSON(http.StatusConflict, messageResponseBody("登録済みのメールアドレスです"))
		} else {
			c.JSON(responseBody500UnexpectedError())
		}
		return
	}

	// ユーザ登録に成功した場合、ログイン状態にする
	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Save()

	c.JSON(http.StatusCreated, gin.H{"user": user})
}