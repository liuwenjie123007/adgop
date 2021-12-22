package main

import (
	"log"
	"net"
	"net/rpc"
)

const HelloServiceName = "path/to/pkg.HelloService"

type HelloService struct {
}

type HelloServiceInterface = interface {
	Hello(request string, reply *string) error
}

func (p *HelloService) Hello(request string, reply *string) error {
	*reply = "hello:" + request
	return nil
}

func RegisterHelloService(svc HelloServiceInterface) error {
	return rpc.RegisterName(HelloServiceName, svc)
}

func main() {
	RegisterHelloService(new(HelloService))
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}
		// 这个只会服务一次
		rpc.ServeConn(conn)
	}
}
