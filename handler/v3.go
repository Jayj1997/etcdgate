/*
 * @Author       : jayj
 * @Date         : 2021-11-17 21:47:36
 * @Description  : etcd v3 handler
 */
package handler

import (
	"confcenter/service"
	"confcenter/utils"
	"confcenter/utils/res"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
		logrus.Errorln("empty etcd address")
		res.InternalError_(ctx, "etcd address is required")
		return
	}

	user := &service.User{
		Username: username,
		Password: password,
		Address:  address,
	}

	if err := v3.s.Auth(user); err != nil {
		logrus.Errorln(err)
		res.InternalError_(ctx, err.Error())
		return
	}

	token, err := utils.GenerateToken(address, username, password)
	if err != nil {
		logrus.Errorln(err)
		res.InternalError_(ctx, err.Error())
		return
	}

	res.Ok(ctx, res.OK, token)
}

// Get
// pass rev to get target revision of value
func (v3 *EtcdV3) Get(ctx *gin.Context) {
	user := getUser(ctx)

	key := ctx.PostForm("key")

	rev, _ := strconv.Atoi(ctx.PostForm("rev"))

	resp, err := v3.s.Get(user, key, int64(rev))
	if err != nil {
		logrus.Errorln(err)
		res.InternalError_(ctx, err.Error())
		return
	}

	res.Ok(ctx, res.OK, resp)
}

func (v3 *EtcdV3) Put(ctx *gin.Context) {

	user := getUser(ctx)

	key := ctx.PostForm("key")
	val := ctx.PostForm("val")

	resp, err := v3.s.Put(user, key, val)
	if err != nil {
		logrus.Errorln(err)
		res.InternalError_(ctx, err.Error())
		return
	}

	res.Ok(ctx, res.OK, resp)
}

func (v3 *EtcdV3) Del(ctx *gin.Context) {

	user := getUser(ctx)

	key := ctx.PostForm("key")
	isDir := ctx.PostForm("dir")

	if err := v3.s.Del(user, key, isDir == "true"); err != nil {
		logrus.Errorln(err)
		res.InternalError_(ctx, err.Error())
	}

	res.Ok_(ctx)
}

func (v3 *EtcdV3) Directory(ctx *gin.Context) {

	user := getUser(ctx)

	path, err := v3.s.GetDirectory(user)
	if err != nil {
		logrus.Errorln(err)
		res.InternalError_(ctx, err.Error())
		return
	}

	res.Ok(ctx, res.OK, path)
}

func (v3 *EtcdV3) Users(ctx *gin.Context) {

	if !v3.s.IsRoot(getUser(ctx)) {
		res.Unauthorized(ctx, res.NotRoot)
	}

	userList, err := v3.s.Users()
	if err != nil {
		res.InternalError_(ctx, err.Error())
		return
	}

	res.Ok(ctx, res.OK, userList)
}

func (v3 *EtcdV3) User(ctx *gin.Context) {

	if !v3.s.IsRoot(getUser(ctx)) {
		res.Unauthorized(ctx, res.NotRoot)
	}

	name := ctx.Param("name")

	userInfo, err := v3.s.User(name)
	if err != nil {
		res.InternalError_(ctx, err.Error())
		return
	}

	res.Ok(ctx, res.OK, userInfo)
}

func (v3 *EtcdV3) UserAdd(ctx *gin.Context) {

	if !v3.s.IsRoot(getUser(ctx)) {
		res.Unauthorized(ctx, res.NotRoot)
	}

	name := ctx.PostForm("name")
	pwd := ctx.PostForm("pwd")

	if name == "" || pwd == "" {
		res.Error(ctx, http.StatusForbidden, res.ParamsInvalid)
		return
	}

	resp, err := v3.s.UserAdd(name, pwd)
	if err != nil {
		res.InternalError_(ctx, err.Error())
		return
	}

	res.Ok(ctx, res.OK, resp)
}

func (v3 *EtcdV3) UserDelete(ctx *gin.Context) {

	if !v3.s.IsRoot(getUser(ctx)) {
		res.Unauthorized(ctx, res.NotRoot)
	}

	name := ctx.Param("name")

	resp, err := v3.s.UserDelete(name)
	if err != nil {
		res.InternalError_(ctx, err.Error())
		return
	}

	res.Ok(ctx, res.OK, resp)
}

func (v3 *EtcdV3) UserGrant(ctx *gin.Context) {

	if !v3.s.IsRoot(getUser(ctx)) {
		res.Unauthorized(ctx, res.NotRoot)
	}

	name := ctx.PostForm("name")
	role := ctx.PostForm("role")

	resp, err := v3.s.UserGrant(name, role)
	if err != nil {
		res.InternalError_(ctx, err.Error())
		return
	}

	res.Ok(ctx, res.OK, resp)
}

func (v3 *EtcdV3) UserRevoke(ctx *gin.Context) {
	if !v3.s.IsRoot(getUser(ctx)) {
		res.Unauthorized(ctx, res.NotRoot)
	}

	name := ctx.PostForm("name")
	role := ctx.PostForm("role")

	resp, err := v3.s.UserRevoke(name, role)
	if err != nil {
		res.InternalError_(ctx, err.Error())
		return
	}

	res.Ok(ctx, res.OK, resp)
}

func (v3 *EtcdV3) Roles(ctx *gin.Context) {

	if !v3.s.IsRoot(getUser(ctx)) {
		res.Unauthorized(ctx, res.NotRoot)
	}

	resp, err := v3.s.Roles()
	if err != nil {
		res.InternalError_(ctx, err.Error())
		return
	}

	res.Ok(ctx, res.OK, resp)
}

func (v3 *EtcdV3) Role(ctx *gin.Context) {

	if !v3.s.IsRoot(getUser(ctx)) {
		res.Unauthorized(ctx, res.NotRoot)
	}

	name := ctx.Param("name")

	resp, err := v3.s.Role(name)
	if err != nil {
		res.InternalError_(ctx, err.Error())
		return
	}

	res.Ok(ctx, res.OK, resp)
}

// RoleAdd adds a new role to an etcd cluster.
func (v3 *EtcdV3) RoleAdd(ctx *gin.Context) {

	if !v3.s.IsRoot(getUser(ctx)) {
		res.Unauthorized(ctx, res.NotRoot)
	}

	name := ctx.Param("name")

	resp, err := v3.s.RoleAdd(name)
	if err != nil {
		res.InternalError_(ctx, err.Error())
		return
	}

	res.Ok(ctx, res.OK, resp)
}

// RoleDelete deletes a role.
func (v3 *EtcdV3) RoleDelete(ctx *gin.Context) {

	if !v3.s.IsRoot(getUser(ctx)) {
		res.Unauthorized(ctx, res.NotRoot)
	}

	name := ctx.Param("name")

	resp, err := v3.s.RoleDelete(name)
	if err != nil {
		res.InternalError_(ctx, err.Error())
		return
	}

	res.Ok(ctx, res.OK, resp)
}

func (v3 *EtcdV3) RoleGrant(ctx *gin.Context) {

	if !v3.s.IsRoot(getUser(ctx)) {
		res.Unauthorized(ctx, res.NotRoot)
	}

	roleName := ctx.PostForm("role_name")
	key := ctx.PostForm("key")
	rangeEnd := ctx.PostForm("range_end")
	permissionTypeStr := ctx.PostForm("type")

	permissionType, _ := strconv.Atoi(permissionTypeStr)

	resp, err := v3.s.RoleGrant(roleName, key, rangeEnd, int32(permissionType))
	if err != nil {
		res.InternalError_(ctx, err.Error())
		return
	}

	res.Ok(ctx, res.OK, resp)
}

func (v3 *EtcdV3) RoleRevoke(ctx *gin.Context) {

	if !v3.s.IsRoot(getUser(ctx)) {
		res.Unauthorized(ctx, res.NotRoot)
	}

	roleName := ctx.PostForm("role_name")
	key := ctx.PostForm("key")
	rangeEnd := ctx.PostForm("range_end")

	resp, err := v3.s.RoleRevoke(roleName, key, rangeEnd)
	if err != nil {
		res.InternalError_(ctx, err.Error())
		return
	}

	res.Ok(ctx, res.OK, resp)
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
