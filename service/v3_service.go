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
	"strings"
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
	Mu          sync.RWMutex // key-val read/write lock
	root        *User
}

// User use to make connection
type User struct {
	Username string // enabled when IsAuth=true
	Password string // enabled when IsAuth=true
	Address  string // etcd address
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
// it will create a Root account
// if auth is enabled and root account are not create
func (e *EtcdV3Service) IfRootAccount(user, pwd, addr string) error {
	if e.IsAuth {
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

// TODO add a more
// IsRoot check if is root account
func (e *EtcdV3Service) IsRoot(user *User) bool {
	return user.Username == e.root.Username && user.Password == e.root.Password
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

// Auth test connection by current User{}
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
		return nil, err
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
		return nil, err
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
		return nil, err
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), e.DialTimeout)

	kv := clientv3.NewKV(cli)

	resp, err := kv.Put(ctx, key, val, clientv3.WithPrevKV())
	cancel()
	if err != nil {
		switch err {
		case context.Canceled:
			return resp, fmt.Errorf("ctx is canceled by another routine: %v", err)
		case context.DeadlineExceeded:
			return resp, fmt.Errorf("ctx is attached with a deadline is exceeded: %v", err)
		case rpctypes.ErrEmptyKey:
			return resp, fmt.Errorf("client-side error: %v", err)
		default:
			return resp, fmt.Errorf("bad cluster endpoints, which are not etcd servers: %v", err)
		}
	}

	// old key-val
	fmt.Println(resp.PrevKv)

	return resp, nil
}

// Del delete dir/* if isDir
func (e *EtcdV3Service) Del(user *User, key string, isDir bool) error {
	e.Mu.Lock()
	defer e.Mu.Unlock()

	var err error

	cli, err := e.connect(user)
	if err != nil {
		return err
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), e.DialTimeout)
	defer cancel()

	if isDir {
		// delete key/*
		_, err = cli.Delete(ctx, key+e.Separator, clientv3.WithPrefix())
	} else {
		_, err = cli.Delete(ctx, key)
	}

	return err
}

type Directory struct {
	IsNode   bool                 `json:"is_node"`
	Children map[string]Directory `json:"children,omitempty"`
}

func (e *EtcdV3Service) GetDirectory(user *User) (interface{}, error) {

	e.Mu.RLock()
	defer e.Mu.RUnlock()

	cli, err := e.connect(user)
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	all, err := cli.Get(context.Background(), e.Separator, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
	if err != nil {
		return nil, err
	}

	dir := map[string]Directory{
		e.Separator: {
			Children: map[string]Directory{},
			IsNode:   false,
		},
	}

	for _, key := range all.Kvs {

		var (
			exist    bool      = false
			isNode   bool      = false
			cur      Directory = dir[e.Separator]
			splitKey []string  = strings.Split(string(key.Key), e.Separator)
		)

		for index, val := range splitKey {

			// head
			if val == "" {
				continue
			}

			// last one
			// there shouldn't be just directory like
			// /exampleA/exampleB/
			if index == len(splitKey)-1 {
				isNode = true
			}

			if _, exist = cur.Children[val]; !exist {
				cur.Children[val] = Directory{
					Children: map[string]Directory{},
					IsNode:   isNode,
				}
			}

			cur = cur.Children[val]
		}
	}

	resp := map[string]interface{}{
		"total":     all.Count,
		"is_more":   all.More,
		"directory": dir,
	}

	return resp, err
}

// TODO add role/permission related lock
// Users get all users, root only
func (e *EtcdV3Service) Users() (interface{}, error) {

	// e.Mu.RLock()
	// defer e.Mu.RUnlock()

	rootCli, err := e.connect(e.root)
	if err != nil {
		return nil, err
	}
	defer rootCli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), e.DialTimeout)

	userList, err := rootCli.UserList(ctx)
	cancel()
	if err != nil {
		return nil, err
	}

	return userList, nil
}

// get a detailed information of a user (role detail)
func (e *EtcdV3Service) User(name string) (interface{}, error) {

	rootCli, err := e.connect(e.root)
	if err != nil {
		return nil, err
	}
	defer rootCli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), e.DialTimeout)

	userInfo, err := rootCli.UserGet(ctx, name)
	cancel()
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}

// UserAdd adds a user
func (e *EtcdV3Service) UserAdd(name, pwd string) (interface{}, error) {

	rootCli, err := e.connect(e.root)
	if err != nil {
		return nil, err
	}
	defer rootCli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), e.DialTimeout)

	resp, err := rootCli.UserAdd(ctx, name, pwd)
	cancel()
	if err != nil {
		return nil, err
	}

	return resp, err
}

// UserDelete delete a user
func (e *EtcdV3Service) UserDelete(name string) (interface{}, error) {

	rootCli, err := e.connect(e.root)
	if err != nil {
		return nil, err
	}
	defer rootCli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), e.DialTimeout)

	resp, err := rootCli.UserDelete(ctx, name)
	cancel()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// UserGrant grant user a role
func (e *EtcdV3Service) UserGrant(name, role string) (interface{}, error) {

	rootCli, err := e.connect(e.root)
	if err != nil {
		return nil, err
	}
	defer rootCli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), e.DialTimeout)

	resp, err := rootCli.UserGrantRole(ctx, name, role)
	cancel()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// UserRevoke revoke user a role
func (e *EtcdV3Service) UserRevoke(name, role string) (interface{}, error) {

	rootCli, err := e.connect(e.root)
	if err != nil {
		return nil, err
	}
	defer rootCli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), e.DialTimeout)

	resp, err := rootCli.UserRevokeRole(ctx, name, role)
	cancel()
	if err != nil {
		return nil, err
	}

	return resp, nil
}
