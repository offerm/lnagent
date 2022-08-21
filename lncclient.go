package lnagent

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/proxy"
	"google.golang.org/grpc"
	"net"
)

var socksAddress string

func torDialer(_ context.Context, addr string) (net.Conn, error) {
	dialer, err := proxy.SOCKS5("tcp", socksAddress, nil, nil)
	if err != nil {
		return nil, err
	}
	conn, err := dialer.Dial("tcp", addr)
	return conn, err
}

func GetLNCClientConn(config *Config) *grpc.ClientConn {

	RpcHostPort := fmt.Sprintf("%v:%v", config.Host, config.Port)

	opts := []grpc.DialOption{grpc.WithInsecure()}
	if config.TorActive {
		socksAddress = config.TorSocks
		opts = append(opts, grpc.WithContextDialer(torDialer))
	}
	conn, err := grpc.Dial(RpcHostPort, opts...)
	if err != nil {
		log.Panicf("unable to connect with the coordinator at %v - %v", RpcHostPort, err)
	}
	return conn
}
