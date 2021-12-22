package main

import (
	"adgop/ch4/ch4.2/hello.pb"
	"log"
	"net"
	"net/rpc"
)

type HelloService struct{}

func (p *HelloService) Hello(in *hello.String, out *hello.String) error {
	out.Value = "Hello:" + in.GetValue()
	return nil
}

func main() {
	server := rpc.NewServer()
	hello.RegisterHelloService(server, new(HelloService))

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}
		go rpc.ServeConn(conn)
	}
}
