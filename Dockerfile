FROM golang:1.17.3 as build

ENV GO111MODULE on
ENV GOPROXY "https://goproxy.io"

RUN mkdir /opt/etcdgate
