/*
 * @Author       : jayj
 * @Date         : 2021-12-06 13:56:02
 * @Description  : etcd user related
 * @LastEditors  : jayj
 * @LastEditTime : 2021-12-07 14:16:30
 */
package service

import "context"

// User use to make connection
type User struct {
	Username string // enabled when IsAuth=true
	Password string // enabled when IsAuth=true
	Address  string // etcd address
}

// Users get all users, root only
func (e *EtcdV3Service) Users(user *User) (interface{}, error) {

	// e.Mu.RLock()
	// defer e.Mu.RUnlock()

	rootCli, err := e.connect(user)
	if err != nil {
		return nil, whichError(err)
	}
	defer rootCli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), e.DialTimeout)

	userList, err := rootCli.UserList(ctx)
	cancel()
	if err != nil {
		return nil, whichError(err)
	}

	return userList, nil
}

// get a detailed information of a user (role detail)
func (e *EtcdV3Service) User(user *User, name string) ([]string, error) {

	rootCli, err := e.connect(user)
	if err != nil {
		return nil, whichError(err)
	}
	defer rootCli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), e.DialTimeout)

	userInfo, err := rootCli.UserGet(ctx, name)
	cancel()
	if err != nil {
		return nil, whichError(err)
	}

	return userInfo.Roles, nil
}

// UserAdd adds a user
func (e *EtcdV3Service) UserAdd(user *User, name, pwd string) (interface{}, error) {

	rootCli, err := e.connect(user)
	if err != nil {
		return nil, whichError(err)
	}
	defer rootCli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), e.DialTimeout)

	resp, err := rootCli.UserAdd(ctx, name, pwd)
	cancel()
	if err != nil {
		return nil, whichError(err)
	}

	return resp, err
}

// UserDelete delete a user
func (e *EtcdV3Service) UserDelete(user *User, name string) (interface{}, error) {

	e.Mu.Lock()
	defer e.Mu.Unlock()

	rootCli, err := e.connect(user)
	if err != nil {
		return nil, whichError(err)
	}
	defer rootCli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), e.DialTimeout)

	resp, err := rootCli.UserDelete(ctx, name)
	cancel()
	if err != nil {
		return nil, whichError(err)
	}

	return resp, nil
}

// UserGrant grant user a role
func (e *EtcdV3Service) UserGrant(user *User, name, role string) (interface{}, error) {

	rootCli, err := e.connect(user)
	if err != nil {
		return nil, whichError(err)
	}
	defer rootCli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), e.DialTimeout)

	resp, err := rootCli.UserGrantRole(ctx, name, role)
	cancel()
	if err != nil {
		return nil, whichError(err)
	}

	return resp, nil
}

// UserRevoke revoke user a role
func (e *EtcdV3Service) UserRevoke(user *User, name, role string) (interface{}, error) {

	e.Mu.Lock()
	defer e.Mu.Unlock()

	rootCli, err := e.connect(user)
	if err != nil {
		return nil, whichError(err)
	}
	defer rootCli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), e.DialTimeout)

	resp, err := rootCli.UserRevokeRole(ctx, name, role)
	cancel()
	if err != nil {
		return nil, whichError(err)
	}

	return resp, nil
}
