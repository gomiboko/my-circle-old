package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gomiboko/my-circle/consts"
)

// メッセージを含むレスポンスボディを生成する
func MessageResponseBody(message string) gin.H {
	return gin.H{"message": message}
}

func ResponseBody400BadRequest() (int, gin.H) {
	return http.StatusBadRequest, MessageResponseBody(consts.Msg400BadRequest)
}

func ResponseBody500UnexpectedError() (int, gin.H) {
	return http.StatusInternalServerError, MessageResponseBody(consts.Msg500UnexpectedError)
}
