/*
 * @Author       : jayj
 * @Date         : 2021-11-13 20:20:19
 * @Description  : etcd interactive interface
 */
package service

type EtcdService interface {
	Connect() error
	Get() error
	Put() error
	Del() error
	Path() error // get directory
}

// func NewEtcdServiceV3() EtcdService {
// 	return &EtcdV3{}
// }
