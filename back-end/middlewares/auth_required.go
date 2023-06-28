package middlewares

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gomiboko/my-circle/consts"
	"github.com/gomiboko/my-circle/utils"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		isLoggedIn := session.Get(consts.SessKeyUserID) != nil

		if !isLoggedIn {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.MessageResponseBody(consts.MsgNeedToLogin))
		}
	}
}
