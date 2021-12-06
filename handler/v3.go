/*
 * @Author       : jayj
 * @Date         : 2021-11-17 21:47:36
 * @Description  : etcd v3 handler
 */
package handler

import (
	"etcdgate/service"
	"etcdgate/utils"
	"etcdgate/utils/res"

	"github.com/gin-gonic/gin"
)

type EtcdV3 struct {
	s *service.EtcdV3Service
}

func CreateEtcdV3Handler(v3 *service.EtcdV3Service) *EtcdV3 {
	return &EtcdV3{s: v3}
}

// Auth you must Auth first to get token
// no matter with or without etcd auth enabled
// it's for other methods to get address
func (v3 *EtcdV3) Auth(ctx *gin.Context) {

	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	address := ctx.PostForm("address")

	if address == "" {
		res.InternalError_(ctx, "etcd address is required")
		return
	}

	user := &service.User{
		Username: username,
		Password: password,
		Address:  address,
	}

	if err := v3.s.Auth(user); err != nil {
		res.InternalError_(ctx, err.Error())
		return
	}

	token, err := utils.GenerateToken(address, username, password)
	if err != nil {
		res.InternalError_(ctx, err.Error())
		return
	}

	res.Ok(ctx, res.OK, token)
}

func (v3 *EtcdV3) Get(ctx *gin.Context) {
	user := getUser(ctx)

	key := ctx.PostForm("key")

	resp, err := v3.s.Get(user, key)
	if err != nil {
		res.InternalError_(ctx, err.Error())
		return
	}

	res.Ok(ctx, res.OK, resp)
}

func (v3 *EtcdV3) Put(ctx *gin.Context) {

	user := getUser(ctx)

	key := ctx.PostForm("key")
	val := ctx.PostForm("val")

	if err := v3.s.Put(user, key, val); err != nil {
		res.InternalError_(ctx, err.Error())
		return
	}

	res.Ok_(ctx)
}

func (v3 *EtcdV3) Del(ctx *gin.Context) {

	user := getUser(ctx)

	key := ctx.PostForm("key")
	isDir := ctx.PostForm("dir")

	if err := v3.s.Del(user, key, isDir == "true"); err != nil {
		res.InternalError_(ctx, err.Error())
	}

	res.Ok_(ctx)
}

func (v3 *EtcdV3) Directory(ctx *gin.Context) {

	user := getUser(ctx)

	path, err := v3.s.GetDirectory(user)
	if err != nil {
		res.InternalError_(ctx, err.Error())
		return
	}

	res.Ok(ctx, res.OK, path)
}

func getUser(ctx *gin.Context) *service.User {

	a, _ := ctx.Get("address")
	u, _ := ctx.Get("username")
	p, _ := ctx.Get("password")

	return &service.User{
		Address:  a.(string),
		Username: u.(string),
		Password: p.(string),
	}
}
