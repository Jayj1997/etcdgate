/*
 * @Author       : jayj
 * @Date         : 2021-06-24 14:27:19
 * @Description  :
 */
package res

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 返回200 无data
func Ok_(ctx *gin.Context) {
	Ok(ctx, OK, nil)
}

// 返回200但自定义msg code
func Ok__(ctx *gin.Context, msgCode int) {
	Ok(ctx, msgCode, nil)
}

// 返回200 自定义code data
func Ok(ctx *gin.Context, msgCode int, data interface{}) {
	ctx.JSON(http.StatusOK, ginH(msgCode, data))
}

// 无权限err
func Unauthorized(ctx *gin.Context, msgCode int) {
	ctx.JSON(http.StatusUnauthorized, ginH(msgCode, nil))
}

// 内部err
func InternalError(ctx *gin.Context) {
	ctx.JSON(http.StatusInternalServerError, ginH(InternalServerError, nil))
}

// 禁止访问err
func ForbiddenError(ctx *gin.Context, msgCode int) {
	ctx.JSON(http.StatusForbidden, ginH(msgCode, nil))
}

// 自定义 err
func Error(ctx *gin.Context, httpCode, msgCode int) {
	ctx.JSON(httpCode, ginH(msgCode, nil))
}

func ginH(msgCode int, data interface{}) gin.H {
	return gin.H{
		"code": msgCode,
		"msg":  GetMsg(msgCode),
		"data": data,
	}
}
