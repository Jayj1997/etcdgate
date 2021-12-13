# ETCD gate

A lightweight dashboard for ETCD, determined to build a dynamic configuration service in a centralized web panel, allows you to configure across environments

frontend use vue + arco design, backend use etcd go client interface

## RUN in docker

if you **not familiar** with docker, git clone project, goto [docker-compose](https://github.com/Jayj1997/etcdgate/blob/base/docker-compose.yaml), and run

``` bash
  docker-compose up -d
```

if you **already have a etcd compose file**, add this in your compose file, make sure etcd-gate and etcd in a same network

```docker
etcd-gate:
    image: ifisjayj/etcdgate:0.0.1
    container_name: etcd-gate
    # restart: always
    ports:
      - 8080:8080
    networks:
      - etcd-net # make sure etcd-gate and etcd in a same network
    environment:
      - ADDR=http://etcd1:2379       # etcd address
      - PORT=8080                    # etcd-gate listen port
      - AUTH=false                   # enable auth
      - ROOT=root                    # root username, ignore it if !auth
      - PWD=root                     # root password, ignore it if !auth
      - TIMEOUT=5                    # timeout per request
      - TLS=false                    # enable tls
      - CA=                          #
      - CERT=                        #
      - KEYFILE=                     #
      - SEPARATOR=/                  # root separator
      - GIN_MODE=release             # gin mode, set debug for debug
```

if you'd like to enable auth (through etcd-gate or not), remember set those environment:

```docker
- AUTH=true                    # enable auth
- ROOT=root                    # root username
- PWD=root                     # root password
```

then view <http://localhost:8080/ui>

## RUN code

``` bash
# run application
go mod tidy
go run main.go

# run just frontend
yarn serve

# build to update frontend
yarn build
```

then view <http://localhost:8080/ui>

## Auth enable

if you'd like to enable auth, you must pass auth argument along with address, etcd-gate will try to open etcd auth and create a root account, if **root** and **pwd** is not provide, it will default to root:root

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
* [x] container
* [ ] test
* [ ] v2 support

## Contribute

contribute & PR

we desperately need a frontend. of course, if you'd like to contribute, you can contact me at ifisjayj@gmail.com

## difficulty

we'd like to write a powerful configuration panel that user manage their key-value、permission、 history here, and read config through etcd(instead of etcd-gate). however, it's kind of tricky to let user listen changes or make canary publish easily, just by interact with etcd api. for now, what I can imagine is write another module, cooperate this to use (**lightweight?**). any idea?
