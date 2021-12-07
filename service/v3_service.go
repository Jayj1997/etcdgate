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

// PS: client and etcdctl use different storage to store data
// so you can't use etcdctl to store and client to read
// Get
// rev pass a non-zero number to get target revision of value
func (e *EtcdV3Service) Get(user *User, key string, rev int64) (interface{}, error) {
	e.Mu.RLock()
	defer e.Mu.RUnlock()

	var (
		resp *clientv3.GetResponse
		err  error
	)

	cli, err := e.connect(user)
	if err != nil {
		return nil, whichError(err)
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), e.DialTimeout)

	if rev == 0 {
		resp, err = cli.Get(ctx, key)
	} else {
		resp, err = cli.Get(ctx, key, clientv3.WithRev(rev))
	}
	cancel()
	if err != nil {
		return nil, whichError(err)
	}

	if resp.Count == 0 {
		return nil, errors.New("empty result")
	}

	kv := resp.Kvs[0]

	result := map[string]interface{}{
		"key":             string(kv.Key),
		"value":           string(kv.Value),
		"create_revision": kv.CreateRevision,
		"mod_revision":    kv.ModRevision,
		"ttl":             getTTL(cli, kv.Lease),
	}

	return result, nil
}

func (e *EtcdV3Service) Put(user *User, key, val string) (*clientv3.PutResponse, error) {
	e.Mu.Lock()
	defer e.Mu.Unlock()

	cli, err := e.connect(user)
	if err != nil {
		return nil, whichError(err)
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), e.DialTimeout)

	kv := clientv3.NewKV(cli)

	resp, err := kv.Put(ctx, key, val, clientv3.WithPrevKV())
	cancel()
	if err != nil {
		return resp, whichError(err)
	}

	// old key-val
	fmt.Println(resp.PrevKv)

	return resp, nil
}

// Del delete key
// delete key* if isDir == true
func (e *EtcdV3Service) Del(user *User, key string, isDir bool) error {
	e.Mu.Lock()
	defer e.Mu.Unlock()

	var err error

	cli, err := e.connect(user)
	if err != nil {
		return whichError(err)
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), e.DialTimeout)
	defer cancel()

	if isDir {
		// delete key*
		_, err = cli.Delete(ctx, key, clientv3.WithPrefix())
	} else {
		_, err = cli.Delete(ctx, key)
	}

	if err != nil {
		return whichError(err)
	}

	return nil

}

type directory struct {
	PermType       int   `json:"perm_type"`       // permission type 0 read 1 write 2 readwrite
	CreateRevision int64 `json:"create_revision"` // last creation version
	ModRevision    int64 `json:"mod_revision"`    // last modification version
	RemainingLease int64 `json:"remaining_lease"` // lease remaining seconds
	// value          string `json:"value,omitempty"` // necessary?
}

// GetDirectory get permitted key directory
func (e *EtcdV3Service) GetDirectory(user *User) (interface{}, error) {

	e.Mu.RLock()
	defer e.Mu.RUnlock()

	// dir[key] = directory
	dir := map[string]*directory{}

	resp := map[string]interface{}{
		"total":     0,
		"is_more":   false,
		"directory": dir,
	}

	if e.IsRoot(user) { // if root account
		rootCli, err := e.connect(user)
		if err != nil {
			return nil, whichError(err)
		}
		defer rootCli.Close()

		all, err := rootCli.Get(context.Background(),
			e.Separator,
			clientv3.WithPrefix(),
			clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend)) // all keys
		if err != nil {
			return nil, whichError(err)
		}

		for _, key := range all.Kvs {
			dir[string(key.Key)] = &directory{
				PermType:       2, // readwrite
				CreateRevision: key.CreateRevision,
				ModRevision:    key.ModRevision,
				RemainingLease: getTTL(rootCli, key.Lease),
			}

		}

		resp["total"] = all.Count
		resp["is_more"] = all.More

		return resp, nil
	}

	// if not root
	perms, err := e.getPerms(user)
	if err != nil {
		return nil, err
	}

	cli, err := e.connect(user)
	if err != nil {
		return nil, whichError(err)
	}

	for _, perm := range perms {
		// traversal keys to get all keys with rangeEnd[if have)
		// e.g. key aa rangeEnd ab
		//      then get aa/a aa/b aa/c ....
		keys, err := cli.Get(context.Background(), perm.Key, clientv3.WithRange(perm.RangeEnd))
		if err != nil {
			return nil, whichError(err)
		}

		for _, key := range keys.Kvs {
			// why this
			// e.g. userRoles
			// role1 key:aa   rangeEnd: ab permType: read
			// role2 key aa/1              permType: write
			if oldKey, exist := dir[string(key.Key)]; exist {
				oldKey.PermType = comparePerm(oldKey.PermType, perm.PermType)
				continue
			}

			dir[string(key.Key)] = &directory{
				PermType:       perm.PermType,
				CreateRevision: key.CreateRevision,
				ModRevision:    key.ModRevision,
				RemainingLease: getTTL(cli, key.Lease),
			}
		}
	}

	resp["total"] = len(dir)

	return resp, nil
}

// comparePerm compare permission_type
// e.g. old = 0(read) new = 1(write) return 2(readwrite)
func comparePerm(old, new int) int {

	//  20 21

	// 00 11 22
	if old == new {
		return new
	}

	// 10 01
	if old+new == 1 {
		return 2
	}

	// 20 21
	if old > new {
		return old
	}

	return new
}
