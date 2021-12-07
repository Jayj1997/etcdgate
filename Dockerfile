# multi build

FROM golang:1.17.3 AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY="https://goproxy.cn,direct"

WORKDIR /src

COPY . .

RUN go mod tidy && go build -o etcdgate main.go

FROM alpine:3.13.5

ENV TLS=false
ENV CA=""
ENV CERT=""
ENV KEYFILE=""
ENV TIMEOUT=5
ENV PORT=8080
ENV SEPARATOR="/"
ENV AUTH=false
ENV ROOT="root"
ENV PASSWORD="root"
ENV ADDR="127.0.0.1:2379"
ENV GIN_MODE=release

WORKDIR /src/etcdgate

COPY --from=builder /src/etcdgate /src/etcdgate/etcdgate

ADD ui ui

EXPOSE ${PORT}

ENTRYPOINT  ./etcdgate -tls=$TLS -ca=$CA -cert=$CERT -keyfile=$KEYFILE \
    -timeout=$TIMEOUT -port=$PORT -separator=$SEPARATOR -auth=$AUTH \
    -root=$ROOT -pwd=$PASSWORD -addr=$ADDR