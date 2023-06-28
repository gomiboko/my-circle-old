package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const maxRequestBodySize = 1024 * 1024

func RequestBodySizeLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		// リクエストボディサイズの上限を1MiBに設定
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxRequestBodySize)

		c.Next()
	}
}
