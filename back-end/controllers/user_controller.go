package controllers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gomiboko/my-circle/db"
	"github.com/gomiboko/my-circle/forms"
	"github.com/gomiboko/my-circle/services"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(us services.UserService) *UserController {
	uc := &UserController{
		userService: us,
	}
	return uc
}

func (uc *UserController) Create(c *gin.Context) {
	// 入力チェック
	var form forms.UserForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(responseBody400BadRequest())
		return
	}

	// ユーザ登録
	user, err := uc.userService.CreateUser(form)
	if err != nil {
		if db.Is(err, db.ErrDuplicateEntry) {
			c.JSON(http.StatusConflict, MessageResponseBody("登録済みのメールアドレスです"))
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

func (uc *UserController) GetHomeInfo(c *gin.Context) {
	session := sessions.Default(c)
	userId := session.Get("user_id").(uint)

	user, err := uc.userService.GetHomeInfo(userId)

	if err != nil {
		c.JSON(responseBody500UnexpectedError())
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
