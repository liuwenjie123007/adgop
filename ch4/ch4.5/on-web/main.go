package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net/http"
	"strings"
	"time"
)

const port = ":5010"

type myGrpcServer struct{}

func (s *myGrpcServer) SayHello(ctx context.Context, in *HelloRequest) (*HelloReply, error) {
	return &HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	go startServer()
	time.Sleep(time.Second)

	doClientWork()
}

func doClientWork() {
	creds, err := credentials.NewClientTLSFromFile("tls-config/server.crt", "server.grpc.io")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial("localhost"+port, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := NewGreeterClient(conn)

	r, err := c.SayHello(context.Background(), &HelloRequest{Name: "gopher"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("doClientWork: %s", r.Message)
}

func startServer() {
	creds, err := credentials.NewServerTLSFromFile("tls-config/server.crt", "tls-config/server.key")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	RegisterGreeterServer(grpcServer, new(myGrpcServer))

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(w, "hello")
	})

	http.ListenAndServeTLS(port, "tls-config/server.crt", "tls-config/server.key", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMinor == 2 && strings.Contains(r.Header.Get("content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			mux.ServeHTTP(w, r)
		}
	}))
}
