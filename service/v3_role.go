/*
 * @Author       : jayj
 * @Date         : 2021-12-06 13:56:34
 * @Description  : etcd role related
 * @LastEditors  : jayj
 * @LastEditTime : 2021-12-06 14:53:41
 */
package service

import (
	"context"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// Roles Get all roles
func (e *EtcdV3Service) Roles() (interface{}, error) {

	rootCli, err := e.connect(e.root)
	if err != nil {
		return nil, err
	}
	defer rootCli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), e.DialTimeout)

	roleList, err := rootCli.RoleList(ctx)
	cancel()
	if err != nil {
		return nil, whichError(err)
	}

	return roleList, nil
}

type Permissions struct {
	Key      string `json:"key"`
	RangeEnd string `json:"range_end"`
	PermType int    `json:"perm_type"`
}

func (e *EtcdV3Service) Role(roleName string) ([]Permissions, error) {

	rootCli, err := e.connect(e.root)
	if err != nil {
		return nil, whichError(err)
	}
	defer rootCli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), e.DialTimeout)

	role, err := rootCli.RoleGet(ctx, roleName)
	cancel()
	if err != nil {
		return nil, whichError(err)
	}

	perms := []Permissions{}

	for _, p := range role.Perm {
		perms = append(perms, Permissions{
			Key:      string(p.Key),
			RangeEnd: string(p.RangeEnd),
			PermType: int(p.PermType),
		})
	}

	return perms, nil
}

func (e *EtcdV3Service) RoleAdd(roleName string) (interface{}, error) {

	rootCli, err := e.connect(e.root)
	if err != nil {
		return nil, whichError(err)
	}
	defer rootCli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), e.DialTimeout)

	resp, err := rootCli.RoleAdd(ctx, roleName)
	cancel()
	if err != nil {
		return nil, whichError(err)
	}

	return resp, nil
}

func (e *EtcdV3Service) RoleDelete(roleName string) (interface{}, error) {

	e.Mu.Lock()
	defer e.Mu.Unlock()

	rootCli, err := e.connect(e.root)
	if err != nil {
		return nil, whichError(err)
	}
	defer rootCli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), e.DialTimeout)

	resp, err := rootCli.RoleDelete(ctx, roleName)
	cancel()
	if err != nil {
		return nil, whichError(err)
	}

	return resp, nil
}

// auth.pb.go
// const (
// 	READ      Permission_Type = 0
// 	WRITE     Permission_Type = 1
// 	READWRITE Permission_Type = 2
// )
// RoleGrant Grants a key to a role
// official says that [key, rangeEnd), but what I test is [key, rangeEnd]
func (e *EtcdV3Service) RoleGrant(roleName, key, rangeEnd string, permissionType int32) (interface{}, error) {

	rootCli, err := e.connect(e.root)
	if err != nil {
		return nil, whichError(err)
	}
	defer rootCli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), e.DialTimeout)

	resp, err := rootCli.RoleGrantPermission(ctx, roleName, key, rangeEnd, clientv3.PermissionType(permissionType))
	cancel()
	if err != nil {
		return nil, whichError(err)
	}

	return resp, nil
}

func (e *EtcdV3Service) RoleRevoke(roleName, key, rangeEnd string) (interface{}, error) {

	rootCli, err := e.connect(e.root)
	if err != nil {
		return nil, whichError(err)
	}
	defer rootCli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), e.DialTimeout)

	resp, err := rootCli.RoleRevokePermission(ctx, roleName, key, rangeEnd)
	cancel()
	if err != nil {
		return nil, whichError(err)
	}

	return resp, nil
}
