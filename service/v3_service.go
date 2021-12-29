/*
 * @Author       : jayj
 * @Date         : 2021-11-13 20:23:22
 * @Description  : etcd interactive v3 func
 */
package service

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/pkg/transport"
	"google.golang.org/grpc"
)

type EtcdV3Service struct {
	IsAuth      bool
	IsTls       bool
	Cert        string
	KeyFile     string
	CaFile      string
	DialTimeout time.Duration
	Separator   string
	Mu          sync.RWMutex // read/write lock
	root        *User
}

// connect
// Make sure to close the client after using it
// If the client is not closed, the connection will have leaky goroutines.
// https://github.com/etcd-io/etcd/tree/main/client/v3#get-started
func (e *EtcdV3Service) connect(user *User) (*clientv3.Client, error) {

	// tls related
	var tlsConf *tls.Config
	var err error

	if e.IsTls {
		tlsInfo := transport.TLSInfo{
			CertFile:      e.Cert,
			KeyFile:       e.KeyFile,
			TrustedCAFile: e.CaFile,
		}

		tlsConf, err = tlsInfo.ClientConfig()
		if err != nil {
			return nil, err
		}
	}

	conf := clientv3.Config{
		Endpoints:   []string{user.Address},
		DialTimeout: e.DialTimeout, // is this necessary to configurate?
		TLS:         tlsConf,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	}

	if e.IsAuth {
		if user.Username == "" || user.Password == "" {
			return nil, errors.New("empty account or password")
		}

		conf.Username = user.Username
		conf.Password = user.Password
	}

	cli, err := clientv3.New(conf)
	if err != nil {
		return nil, err
	}

	return cli, nil
}

// why on earth i write this
// IfRootAccount
// will create a Root account
// if auth is enabled and root account are not create
func (e *EtcdV3Service) IfRootAccount(user, pwd, addr string) error {
	if e.IsAuth && user != "" && pwd != "" {
		e.root = &User{
			Address:  addr,
			Username: user,
			Password: pwd,
		}

		cli, err := clientv3.New(clientv3.Config{
			Endpoints:   []string{e.root.Address},
			DialTimeout: e.DialTimeout,
		})
		if err != nil {
			return err
		}
		defer cli.Close()

		if _, err := cli.RoleAdd(context.TODO(), "root"); err != nil {
			return err
		}

		if _, err := cli.UserAdd(context.TODO(), e.root.Username, e.root.Password); err != nil {
			return err
		}

		if _, err := cli.UserGrantRole(context.TODO(), e.root.Username, "root"); err != nil {
			return err
		}

		if _, err := cli.AuthEnable(context.TODO()); err != nil {
			return err
		}

		logrus.Infoln("root account CREATED")
	}

	return nil
}

// getTTL get lease time-to-live
func getTTL(cli *clientv3.Client, lease int64) int64 {
	if resp, err := cli.Lease.TimeToLive(context.Background(), clientv3.LeaseID(lease)); err != nil {
		return 0
	} else if resp.TTL < 0 {
		return 0
	} else {
		return resp.TTL
	}
}

// getPerms get permissions by current user
// use to non-root user
func (e *EtcdV3Service) getPerms(user *User) ([]Permissions, error) {

	roles, err := e.User(user, user.Username)
	if err != nil {
		return nil, err
	}

	perms := []Permissions{}

	for _, role := range roles {
		perm, err := e.Role(user, role)
		if err != nil {
			return perm, err
		}

		perms = append(perms, perm...)
	}

	return perms, nil
}

// TODO more precise
// IsRoot check if is root account
func (e *EtcdV3Service) IsRoot(user *User) bool {
	if !e.IsAuth {
		return true
	}

	return user.Username == e.root.Username && user.Password == e.root.Password
}

// return format error
func whichError(err error) error {
	switch err {
	case context.Canceled:
		return fmt.Errorf("ctx is canceled by another routine: %v", err)
	case context.DeadlineExceeded:
		return fmt.Errorf("ctx is attached with a deadline is exceeded: %v", err)
	case rpctypes.ErrEmptyKey:
		return fmt.Errorf("client-side error: %v", err)
	case rpctypes.ErrPermissionDenied:
		return fmt.Errorf("server-side error: %v", err)
	default:
		return fmt.Errorf("error occurred, this may caused by bad cluster endpoints, error: %v", err)
	}
}

// Auth test connection by given user
func (e *EtcdV3Service) Auth(user *User) error {
	_, err := e.connect(user)

	return err
}
