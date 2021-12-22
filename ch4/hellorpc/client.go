package main

import (
	"fmt"
	"log"
	"net/rpc"
)

// 由于都在main包， 实际运行期间是无法引用的
//const HelloServiceName = "path/to/pkg.HelloService"
//
//type HelloServiceInterface = interface {
//	Hello(request string, reply *string) error
//}

type HelloServiceClient struct {
	*rpc.Client
}

func (p *HelloServiceClient) Hello(request string, reply *string) error {
	return p.Client.Call(HelloServiceName+".Hello", request, reply)
}

var _ HelloServiceInterface = (*HelloServiceClient)(nil)

func DialHelloService(network, address string) (*HelloServiceClient, error) {
	c, err := rpc.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return &HelloServiceClient{Client: c}, nil
}

func main() {
	client, err := DialHelloService("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	var reply string
	err = client.Hello("hello", &reply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply)
}
