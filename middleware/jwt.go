/*
 * @Author       : jayj
 * @Date         : 2021-06-24 09:37:46
 * @Description  :
 */
package middleware

import (
	"etcdgate/utils"
	"etcdgate/utils/res"

	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var claims *utils.Claims
		var err error

		token := ctx.GetHeader("Authorization")

		if token == "" {
			res.Unauthorized(ctx, res.TokenInvalid)

			ctx.Abort()
			return
		}

		claims, err = utils.ParseToken(token)
		if claims == nil {
			res.Unauthorized(ctx, res.TokenInvalid)
			ctx.Abort()
			return
		} else if err != nil {
			res.Unauthorized(ctx, res.TokenExpired)
			ctx.Abort()
			return
		}

		ctx.Set("address", claims.Address)
		ctx.Set("username", claims.Username)
		ctx.Set("password", claims.Password)
		ctx.Next()
	}
}
