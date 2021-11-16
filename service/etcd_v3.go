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

// split connect related out
type EtcdV3 struct {
	Addrs       []string
	IsAuth      bool
	UseTls      bool
	Cert        string
	KeyFile     string
	CaFile      string
	acc         string
	pwd         string
	cli         *clientv3.Client
	DialTimeout time.Duration
	mu          sync.Mutex
}

// Make sure to close the client after using it
// If the client is not closed, the connection will have leaky goroutines.
// https://github.com/etcd-io/etcd/tree/main/client/v3#get-started
func (e *EtcdV3) Connect() error {

	// tls related
	var tlsConf *tls.Config
	var err error

	if e.UseTls {
		tlsInfo := transport.TLSInfo{
			CertFile:      e.Cert,
			KeyFile:       e.KeyFile,
			TrustedCAFile: e.CaFile,
		}

		tlsConf, err = tlsInfo.ClientConfig()
		if err != nil {
			return err
		}
	}

	conf := clientv3.Config{
		Endpoints:   e.Addrs,
		DialTimeout: e.DialTimeout, // is this necessary to configurate?
		TLS:         tlsConf,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	}

	if e.IsAuth {
		if e.acc == "" || e.pwd == "" {
			return errors.New("empty account or password")
		}

		conf.Username = e.acc
		conf.Password = e.pwd
	}

	e.cli, err = clientv3.New(conf)
	if err != nil {
		return err
	}

	return nil
}

func (e *EtcdV3) Get() error {
	return nil
}
func (e *EtcdV3) Put(key, val string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), e.DialTimeout)
	defer cancel()

	kv := clientv3.NewKV(e.cli)

	// 记录下原来的kv历史
	resp, err := kv.Put(ctx, key, val, clientv3.WithPrevKV())
	if err != nil {
		return err
	}

	fmt.Println(resp.PrevKv)

	return nil
}
func (e *EtcdV3) Del() error {
	return nil
}
func (e *EtcdV3) Path() error {
	return nil
}
