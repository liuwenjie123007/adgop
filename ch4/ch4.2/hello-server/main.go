package main

import "net/rpc"

func main() {
	server := rpc.NewServer()
	RegisterHelloService
}
