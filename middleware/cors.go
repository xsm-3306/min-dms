package middleware

import (
	"min-dms/common"
	"min-dms/response"
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

//jwt token中间件
func JwtAuthMiddle() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("Authorization") //采用 Bearer
		token = token[7:]

		if len(token) <= 7 {
			msg := "the auth header is empty"
			data := gin.H{}
			response.Failed(ctx, data, msg)
			ctx.Abort()
			return
		}

		calims, err := common.ParseToken(token)
		if err != nil {
			msg := "the token is invalid"

			//response.Failed(ctx, data, msg)
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"err": err,
				"msg": msg,
			})
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//解析出来结果之后，保存username，以供上下文使用
		ctx.Set("username", calims.User.Username)
		ctx.Next()

	}
}
