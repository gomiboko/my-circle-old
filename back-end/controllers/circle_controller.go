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
	"github.com/gomiboko/my-circle/validations"
)

type CircleController struct {
	circleService  services.CircleService
	storageService services.StorageService
}

func NewCircleController(cs services.CircleService, ss services.StorageService) *CircleController {
	return &CircleController{
		circleService:  cs,
		storageService: ss,
	}
}

func (cc *CircleController) Create(c *gin.Context) {
	// 入力チェック
	var form forms.CircleForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(utils.ResponseBody400BadRequest())
		return
	}

	// アイコンファイルの形式チェック
	if validations.IsNotAllowedIconFileFormat(form.CircleIconFile) {
		c.JSON(utils.ResponseBody400BadRequest())
		return
	}

	// アイコンファイルのサイズチェック
	if validations.IsOverMaxIconFileSize(form.CircleIconFile) {
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

	res := gin.H{consts.ResKeyCircle: circle}

	// アイコンファイルをアップロード
	if form.CircleIconFile != nil {
		key := utils.CreateHashedStorageKey(consts.StorageDirCircles, consts.StorageKeyPrefixCircleIcon, circle.ID)

		err = cc.storageService.Upload(key, form.CircleIconFile)
		if err != nil {
			log.Print(err.Error())

			res[consts.ResKeyMessage] = consts.MsgFailedToRegisterCircleIcon
			res[consts.ResKeyMessageType] = consts.MsgTypeWarn
		}
	}

	c.JSON(http.StatusCreated, res)
}
