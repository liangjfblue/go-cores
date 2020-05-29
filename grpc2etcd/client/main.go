package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc/resolver"

	"github.com/liangjfblue/go-cores/grpc2etcd/discovery"

	"github.com/liangjfblue/go-cores/grpc2etcd/proto"

	"google.golang.org/grpc"
)

var (
	num         int
	serviceName = "helloworld"
	etcdAddrs   = []string{"172.16.7.16:9002", "172.16.7.16:9004", "172.16.7.16:9006"}

	defaultMsg = "earth"
)

func init() {
	flag.IntVar(&num, "n", 20, "request num time")
}
func main() {
	//use etcd for resolver
	rl := discovery.NewEtcdBuilder(etcdAddrs, serviceName)
	resolver.Register(rl)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*60)
	defer cancel()

	// Set up a connection to the server.
	for i := 0; i < num; i++ {
		conn, err := grpc.DialContext(
			ctx,
			fmt.Sprintf("%s:///%s", rl.Scheme(), serviceName),
			grpc.WithInsecure(),
			grpc.WithBlock(),
		)
		if err != nil {
			log.Fatalf("did not connect: %v", err)
			return
		}

		c := proto.NewHelloServiceClient(conn)
		r, err := c.HelloWorld(context.TODO(), &proto.HelloReq{Name: defaultMsg})
		if err != nil {
			log.Fatalf("could not hello: %v", err)
			return
		}

		conn.Close()

		log.Println("msg: ", r.GetMsg())
		time.Sleep(time.Second * 1)
	}
}
