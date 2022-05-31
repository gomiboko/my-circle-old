package controllers

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gomiboko/my-circle/consts"
	"github.com/gomiboko/my-circle/forms"
	"github.com/gomiboko/my-circle/services"
	"github.com/gomiboko/my-circle/utils"
)

type SessionController struct {
	sessionService services.SessionService
}

func NewSessionController(ss services.SessionService) *SessionController {
	sc := &SessionController{
		sessionService: ss,
	}
	return sc
}

func (sc SessionController) Create(c *gin.Context) {
	var form forms.LoginForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(utils.ResponseBody400BadRequest())
		return
	}

	// ログイン認証
	userID, err := sc.sessionService.Authenticate(form.Email, form.Password)

	if err != nil {
		log.Print(err)

		c.JSON(utils.ResponseBody500UnexpectedError())
		return
	}

	// 認証失敗
	if userID == nil {
		c.JSON(http.StatusUnauthorized, utils.MessageResponseBody(consts.MsgFailedToLogin))
		return
	}

	session := sessions.Default(c)
	session.Set(consts.SessKeyUserId, *userID)
	session.Save()

	c.Status(http.StatusCreated)
}

func (sc SessionController) Destroy(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Status(http.StatusOK)
}
