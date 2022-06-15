package controllers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gomiboko/my-circle/consts"
	"github.com/gomiboko/my-circle/forms"
	"github.com/gomiboko/my-circle/services"
	"github.com/gomiboko/my-circle/utils"
)

type CircleController struct {
	circleService services.CircleService
}

func NewCircleController(cs services.CircleService) *CircleController {
	return &CircleController{
		circleService: cs,
	}
}

func (cc *CircleController) Create(c *gin.Context) {
	// 入力チェック
	var form forms.CircleForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(utils.ResponseBody400BadRequest())
		return
	}

	userId := sessions.Default(c).Get(consts.SessKeyUserId).(uint)

	// サークル登録
	circle, err := cc.circleService.CreateCircle(form, userId)
	if err != nil {
		c.JSON(utils.ResponseBody500UnexpectedError())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"circle": circle})
}
