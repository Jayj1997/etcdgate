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
func addV3Route(rg, rgWithAuth *gin.RouterGroup, v3Service *service.EtcdV3Service) {

	v3Handler := handler.CreateEtcdV3Handler(v3Service)
	rg.POST("/auth", v3Handler.Auth)

	// key-val related
	rgWithAuth.POST("/get", v3Handler.Get)
	rgWithAuth.POST("/put", v3Handler.Put)
	rgWithAuth.POST("/del", v3Handler.Del)
	rgWithAuth.POST("/directory", v3Handler.Directory) // get path by current permission

	// User related (root only)
	rgWithAuth.GET("/users", v3Handler.Users)                     // gets all users
	rgWithAuth.GET("/user/:name", v3Handler.User)                 // gets detailed information of a user
	rgWithAuth.POST("/user_add", v3Handler.UserAdd)               // adds a new user
	rgWithAuth.DELETE("/user_delete/:name", v3Handler.UserDelete) // deletes a user
	rgWithAuth.POST("/user_grant", v3Handler.UserGrant)           // grants a role to a user
	rgWithAuth.POST("/user_revoke", v3Handler.UserRevoke)         // revokes a role from a user

	// TODO following routes

	// Role related (root only)
	rgWithAuth.GET("roles")           // lists all roles
	rgWithAuth.GET("/role")           // gets detailed information of a role
	rgWithAuth.POST("/role_add")      // adds a new role
	rgWithAuth.DELETE("/role_delete") // deletes a role
	rgWithAuth.POST("/role_grant")    // grants a key to a role
	rgWithAuth.GET("/role_revoke")    // revokes a key from a role
}
