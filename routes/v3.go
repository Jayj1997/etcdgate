/*
 * @Author       : jayj
 * @Date         : 2021-11-13 19:40:39
 * @Description  : etcd v3 funcs
 */
package routes

import (
	"etcdgate/handler"
	"etcdgate/service"

	"github.com/gin-gonic/gin"
)

//TODO should use session?
// addV3Route add etcd v3 route
func addV3Route(rg, rgWithAuth *gin.RouterGroup, v3Service *service.EtcdV3Service) {

	v3Handler := handler.CreateEtcdV3Handler(v3Service)
	rg.POST("/auth", v3Handler.Auth)

	rgWithAuth.POST("/get", v3Handler.Get)
	rgWithAuth.POST("/put", v3Handler.Put)
	rgWithAuth.POST("/del", v3Handler.Del)
	rgWithAuth.POST("/directory", v3Handler.Directory) // get path by current permission

	// TODO following routes
	rgWithAuth.GET("/roles")          // get all roles (root only)
	rgWithAuth.GET("/permissions")    // get all permissions (root only)
	rgWithAuth.GET("/add_permission") // add permission to key (root only)
}
