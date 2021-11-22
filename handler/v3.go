/*
 * @Author       : jayj
 * @Date         : 2021-11-17 21:47:36
 * @Description  : etcd v3 handler
 */
package handler

import (
	"confcenter/service"
	"confcenter/utils/res"

	"github.com/gin-gonic/gin"
)

type EtcdV3 struct {
	s *service.EtcdV3Service
}

func CreateEtcdV3Handler(v3 *service.EtcdV3Service) *EtcdV3 {
	return &EtcdV3{s: v3}
}

func (v3 *EtcdV3) Auth(ctx *gin.Context) {
}

func (v3 *EtcdV3) Get(ctx *gin.Context) {
	user := &service.User{
		Address: "192.168.110.162:2379",
	}

	key := ctx.PostForm("key")

	resp, err := v3.s.Get(user, key)
	if err != nil {
		res.InternalError_(ctx, err.Error())
		return
	}

	res.Ok(ctx, res.OK, resp)
}

func (v3 *EtcdV3) Put(ctx *gin.Context) {

	user := &service.User{
		Address: "192.168.110.162:2379",
	}

	key := ctx.PostForm("key")
	val := ctx.PostForm("val")

	if err := v3.s.Put(user, key, val); err != nil {
		res.InternalError_(ctx, err.Error())
		return
	}

	res.Ok_(ctx)
}

func (v3 *EtcdV3) Del(ctx *gin.Context) {

	user := &service.User{
		Address: "192.168.110.162:2379",
	}

	key := ctx.PostForm("key")

	if err := v3.s.Del(user, key); err != nil {
		res.InternalError_(ctx, err.Error())
		return
	}

	res.Ok_(ctx)
}

// func (v3 *EtcdV3) Path(ctx *gin.Context) {

// 	user := &service.User{
// 		Address: "192.168.110.162:2379",
// 	}

// 	resp, err := v3.s.Path(user)
// 	if err != nil {
// 		res.InternalError_(ctx, err.Error())
// 		return
// 	}

// 	res.Ok(ctx, res.OK, resp)
// }
