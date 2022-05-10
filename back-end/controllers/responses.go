package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// メッセージを含むレスポンスボディを生成する
func MessageResponseBody(message string) gin.H {
	return gin.H{"message": message}
}

func responseBody400BadRequest() (int, gin.H) {
	return http.StatusBadRequest, MessageResponseBody("不正なリクエストです")
}

func responseBody500UnexpectedError() (int, gin.H) {
	return http.StatusInternalServerError, MessageResponseBody("予期せぬエラーが発生しました")
}
