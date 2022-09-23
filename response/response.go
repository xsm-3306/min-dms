package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//response封装统一的返回格式

func Response(ctx *gin.Context, httpStatus int, code int, data gin.H, msg string) {
	ctx.JSON(httpStatus, gin.H{"code": code, "data": data, "msg": msg})
}

//成功httpstatus返回200，code返回200
func Success(ctx *gin.Context, data gin.H, msg string) {
	Response(ctx, http.StatusOK, 200, data, msg)
}

//失败httpstatus返回200，code返回400
func Failed(ctx *gin.Context, data gin.H, msg string) {
	Response(ctx, http.StatusOK, 400, data, msg)
}
