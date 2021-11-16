/*
 * @Author       : jayj
 * @Date         : 2021-11-13 19:36:26
 * @Description  :
 */
package routes

import (
	"confcenter/utils/res"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoadGin() *gin.Engine {

	g := gin.Default()

	loadRoutes(g)

	g.NoRoute(func(ctx *gin.Context) {
		res.Error(ctx, http.StatusNotFound, res.UrlNotFound)
	})

	return g
}

func loadRoutes(g *gin.Engine) {

	v3 := g.Group("/v3")

	addV3Route(v3)

	addFrontend(g)
}
