package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/liangjfblue/go-cores/grpc2etcd/proto"

	"github.com/liangjfblue/go-cores/grpc2etcd/discovery"

	"google.golang.org/grpc"
)

var (
	servicePort int
	serviceIp   string
	serviceName = "helloworld"
	etcdAddrs   = []string{"172.16.7.16:9002", "172.16.7.16:9004", "172.16.7.16:9006"}
)

func init() {
	flag.IntVar(&servicePort, "port", 8090, "server port")
	flag.StringVar(&serviceIp, "ip", "127.0.0.1", "server ip")
}

//go:generate protoc -I ../proto --go_out=plugins=grpc:../proto ../proto/hello.proto

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", servicePort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	etcdRegister := discovery.NewRegister(etcdAddrs, -1)
	if err := etcdRegister.Register(ctx, discovery.ServiceDesc{
		ServiceName: serviceName,
		Host:        serviceIp,
		Port:        servicePort,
		TTL:         time.Second * 3,
	}); err != nil {
		log.Fatal(err)
		return
	}

	s := grpc.NewServer()

	//pb register service
	proto.RegisterHelloServiceServer(s, &Hello{})

	log.Println("server start")
	log.Fatal(s.Serve(lis).Error())
}

type Hello struct {
}

func (h *Hello) HelloWorld(ctx context.Context, in *proto.HelloReq) (*proto.HelloResp, error) {
	log.Println("client recv: ", in.Name)

	out := &proto.HelloResp{
		Code: 1,
		Msg:  "hello " + in.Name,
	}

	return out, nil
}
