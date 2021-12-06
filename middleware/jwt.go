/*
 * @Author       : jayj
 * @Date         : 2021-06-24 09:37:46
 * @Description  :
 */
package middleware

import (
	"etcdgate/utils"
	"etcdgate/utils/res"
	"time"

	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var claims *utils.Claims
		var err error

		code := res.OK

		token := ctx.GetHeader("Authorization")
		if token == "" {
			res.Unauthorized(ctx, res.TokenInvalid)

			ctx.Abort()
			return
		}

		claims, err = utils.ParseToken(token)
		if err != nil {
			res.Unauthorized(ctx, res.TokenInvalid)
			ctx.Abort()
			return
		}

		if time.Now().Unix() > claims.ExpiresAt {
			code = res.TokenExpired
		} else if err != nil {
			code = res.TokenInvalid
		}

		if code == res.TokenExpired {
			res.Ok__(ctx, res.TokenExpired)
		} else if code != res.OK {

			res.Unauthorized(ctx, code)

			ctx.Abort()
			return
		}

		ctx.Set("address", claims.Address)
		ctx.Set("username", claims.Username)
		ctx.Set("password", claims.Password)
		ctx.Next()
	}
}
