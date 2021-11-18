/*
 * @Author       : jayj
 * @Date         : 2021-11-13 19:36:26
 * @Description  :
 */
package routes

import (
	"confcenter/service"
	"confcenter/utils/res"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoadGin(serviceV3 *service.EtcdV3Service) *gin.Engine {

	g := gin.Default()

	v3 := g.Group("/v3")
	addV3Route(v3, serviceV3)

	addFrontend(g)

	g.NoRoute(func(ctx *gin.Context) {
		res.Error(ctx, http.StatusNotFound, res.UrlNotFound)
	})

	return g
}
