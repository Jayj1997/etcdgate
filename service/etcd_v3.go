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
	Mu          sync.RWMutex
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

// getTtl
func getTtl(cli *clientv3.Client, lease int64) int64 {
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

func (e *EtcdV3Service) Get(user *User, key string) (interface{}, error) {
	e.Mu.RLock()
	defer e.Mu.RUnlock()

	cli, err := e.connect(user)
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), e.DialTimeout)
	defer cancel()

	resp, err := cli.Get(ctx, "key")
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
		"is_dir":          false,
		"create_revision": kv.CreateRevision,
		"mod_revision":    kv.ModRevision,
		"ttl":             getTtl(cli, kv.Lease),
	}

	return result, nil
}
func (e *EtcdV3Service) Put(user *User, key, val string) error {
	e.Mu.Lock()
	defer e.Mu.Unlock()

	cli, err := e.connect(user)
	if err != nil {
		return err
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), e.DialTimeout)
	defer cancel()

	kv := clientv3.NewKV(cli)

	// memory old key-val
	resp, err := kv.Put(ctx, key, val, clientv3.WithPrevKV())
	if err != nil {
		return err
	}

	fmt.Println(resp.PrevKv)

	return nil
}
func (e *EtcdV3Service) Del(user *User) error {
	return nil
}
func (e *EtcdV3Service) Path(user *User) error {
	return nil
}
