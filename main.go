/*
 * @Author       : jayj
 * @Date         : 2021-11-15 14:08:34
 * @Description  :
 */
package main

import (
	"etcdgate/routes"
	"etcdgate/service"
	"etcdgate/utils"
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	isTLS   = flag.Bool("tls", false, "enable TLS")
	ca      = flag.String("ca", "", "TLS trusted ca file position")
	cert    = flag.String("cert", "", "TLS cert file position")
	keyfile = flag.String("keyfile", "", "tls keyfile position")
	timeout = flag.Int("timeout", 5, "dial timeout, eg. 5")
	port    = flag.Int("port", 8070, "server listen port, eg. 8080")
	isAuth  = flag.Bool("auth", false, "is etcd auth enabled, enable etcd's auth if not")
	root    = flag.String("root", "root", "etcd root user, default root if not provide")
	pwd     = flag.String("pwd", "root", "etcd root pwd, default root if not provide")
	addr    = flag.String("addr", "127.0.0.1:2379", "etcd address, default 127.0.0.1:2379 if not provide")
	// separator = flag.String("separator", "/", "key separator")
)

func main() {

	flag.Parse()

	v3 := &service.EtcdV3Service{
		IsAuth:      *isAuth,
		IsTls:       *isTLS,
		CaFile:      *ca,
		Cert:        *cert,
		KeyFile:     *keyfile,
		Separator:   "/",
		DialTimeout: time.Duration(*timeout) * time.Second,
		Mu:          sync.RWMutex{},
	}

	err := v3.IfRootAccount(*root, *pwd, *addr)
	if err != nil {
		logrus.Errorln("try create root user failed, IGNORE it if etcd already have a root account, err: ", err.Error())
	}

	g := routes.LoadGin(v3)

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", *port),
		Handler:        g,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logrus.Errorln("server stop or failed, info: ", err)
		}
	}()

	utils.GracefullyShutdown(server)
}
