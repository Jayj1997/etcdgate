/*
 * @Author       : jayj
 * @Date         : 2021-11-13 19:40:39
 * @Description  : etcd v3 funcs
 */
package routes

import "github.com/gin-gonic/gin"

func addV3Route(rg *gin.RouterGroup) {

	rg.POST("/auth")
	rg.POST("/get")
	rg.POST("/put")
	rg.POST("/del")
}
