/*
 * @Author       : jayj
 * @Date         : 2021-11-13 19:40:39
 * @Description  : etcd v3 funcs
 */
package routes

import (
	"confcenter/handler"
	"confcenter/service"

	"github.com/gin-gonic/gin"
)

//TODO should use session?
// addV3Route add etcd v3 route
func addV3Route(rg *gin.RouterGroup, v3Service *service.EtcdV3Service) {

	v3Handler := handler.CreateEtcdV3Handler(v3Service)
	rg.POST("/auth", v3Handler.Auth)
	rg.POST("/get", v3Handler.Get)
	rg.POST("/put", v3Handler.Put)
	rg.POST("/del", v3Handler.Del)
}
