package middleware

import (
	"min-dms/common"
	"min-dms/response"
	"min-dms/userhandler"
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
func JwtAuthMiddle(corUh *userhandler.Userhandler) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("Authorization") //采用 Bearer

		if len(token) <= 7 {
			msg := "the auth header is empty"
			data := gin.H{}
			response.Failed(ctx, data, msg)
			ctx.Abort()
			return
		}

		token = token[7:]
		//token解析之前，先在token可用数据集里面查询一次，这样可以实现toekn注销的功能
		//var corUh userhandler.Userhandler
		sql := "select id from user_authtoken_log where is_deleted=0 and token_str=?"
		result, err := corUh.UserService.Db.GetRows(sql, token)
		//log.Println(result)
		if err != nil {
			msg := "fail to get token info"
			data := gin.H{
				"err": err,
			}
			response.Response(ctx, http.StatusUnauthorized, 400, data, msg)
			ctx.Abort()
			return
		} else {
			if len(result) == 0 {
				msg := "fail to vertify token availability"
				data := gin.H{}
				response.Response(ctx, http.StatusUnauthorized, 400, data, msg)
				ctx.Abort()
				return
			}
		}
		//解析token
		calims, err := common.ParseToken(token)
		if err != nil {
			msg := "the token is invalid"
			data := gin.H{
				"err": err.Error(),
			}
			//response.Failed(ctx, data, msg)
			response.Response(ctx, http.StatusUnauthorized, 400, data, msg)
			ctx.Abort()
			return
		}
		//解析出来结果之后，保存username，以供上下文使用
		ctx.Set("username", calims.User.Username)
		ctx.Next()

	}
}
