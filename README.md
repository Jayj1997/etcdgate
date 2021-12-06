# ETCD configuration center

A lightweight dashboard for ETCD, determined to build a dynamic configuration service in a centralized web panel, allows you to configure across environments

frontend use vue + arco design, backend use etcd go client interface

## Get Started

``` bash
# run application
go mod tidy
go run main.go

# run just frontend
yarn serve

# build to update frontend
yarn build
```

## Auth enable

if you'd like to enable auth, you must pass auth argument along with address, confcenter will try to open etcd auth and create a root account, if **root** and **pwd** is not provide, it will default to root:root

e.g. ` go run main.go --auth=true --addr=exampleurl:2379 `

however, if you already have a root account and enabled auth through etcdctl, you need to pass those arguments:

* auth
* addr
* root
* pwd

e.g. ` go run main.go --auth=true --addr=exampleurl:2379 --root=exampleroot --pwd=examplepassword `

## todo

* [x] base v3 function
  * [x] get/put/del/directory
  * [x] auth
* [ ] frontend
  * [ ] ui
  * [ ] listen opened config
* [ ] finish v3 function
  * [ ] tls
  * [ ] history
  * [ ] namespace
  * [ ] rollback
  * [ ] listen change (maybe)
  * [ ] canary publish (maybe)
* [ ] container
* [ ] test
* [ ] v2 support

## difficulty

we'd like to write a powerful configuration panel that user manage their key-value、permission、 history here, and read config through etcd(instead of confcenter). however, it's kind of tricky to let user listen changes or make canary publish easily, just by interact with etcd api. for now, what I can imagine is write another module, cooperate this to use (**lightweight?**). any idea?
