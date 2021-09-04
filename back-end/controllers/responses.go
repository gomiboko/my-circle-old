package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// メッセージを含むレスポンスボディを生成する
func messageResponseBody(message string) gin.H {
	return gin.H{"message": message}
}

func responseBody400BadRequest() (int, gin.H) {
	return http.StatusBadRequest, messageResponseBody("不正なリクエストです")
}

func responseBody500UnexpectedError() (int, gin.H) {
	return http.StatusInternalServerError, messageResponseBody("予期せぬエラーが発生しました")
}
