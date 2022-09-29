package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//前端请求跨域中间件
func CrosMiddle() gin.HandlerFunc {
	return func(ct *gin.Context) {
		method := ct.Request.Method

		ct.Header("Access-Control-Allow-Origin", "*")
		ct.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, Accept-Language, Content-Language")
		ct.Header("Access-Control-Allow-Methods", "OPTIONS,HEAD,POST,GET,PUT")
		ct.Header("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type")
		ct.Header("Access-Control-Allow-Credentials", "TRUE")
		ct.Header("X-Powered-By", "3.2.1")
		if method == "OPTIONS" {
			//ct.JSON(http.StatusOK, "OPTIONS OK")
			ct.AbortWithStatus(http.StatusNoContent)
		}
	}
}
