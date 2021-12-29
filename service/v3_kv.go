/*
 * @Author       : jayj
 * @Date         : 2021-12-16 14:15:20
 * @Description  :
 * @LastEditors  : jayj
 * @LastEditTime : 2021-12-16 15:01:52
 */
package service

import (
	"context"
	"errors"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type directory struct {
	PermType       int   `json:"perm_type"`       // permission type 0 read 1 write 2 readwrite
	CreateRevision int64 `json:"create_revision"` // last creation version
	ModRevision    int64 `json:"mod_revision"`    // last modification version
	RemainingLease int64 `json:"remaining_lease"` // lease remaining seconds
	// value          string `json:"value,omitempty"` // necessary?
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

// Watch watch key change
func (e *EtcdV3Service) Watch(user *User, key string) error {
	cli, err := e.connect(user)
	if err != nil {
		return whichError(err)
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), e.DialTimeout)
	defer cancel()

	watCh := cli.Watch(ctx, key)
	for resp := range watCh {
		for _, ev := range resp.Events {
			fmt.Println(ev.Type)
			fmt.Println(ev.Kv.Key)
			fmt.Println(ev.Kv.Value)
			fmt.Println(ev.Kv.CreateRevision)
		}
	}

	return nil
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
