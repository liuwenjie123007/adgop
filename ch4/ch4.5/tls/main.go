package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"net"
	"time"
)

var (
	port = ":5001"

	tlsDir        = "./tls-config"
	tlsServerName = "server.io"

	ca        = tlsDir + "/ca.crt"
	serverCrt = tlsDir + "/server.crt"
	serverKey = tlsDir + "/server.key"
	clientCrt = tlsDir + "/client.crt"
	clientKey = tlsDir + "/client.key"
)

type myGrpcServer struct{}

func (s *myGrpcServer) SayHello(ctx context.Context, in *HelloRequest) (*HelloReply, error) {
	return &HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	go startServer()
	time.Sleep(time.Second)

	doClientWork()
}

func startServer() {
	certificate, err := tls.LoadX509KeyPair(serverCrt, serverKey)
	if err != nil {
		log.Panicf("could not load server key pair: %s", err)
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(ca)
	if err != nil {
		log.Panicf("could not read ca certificate: %s", err)
	}

	// Append the client certificates from the CA
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Panic("failed to append client certs")
	}

	// Create the TLS credentials
	creds := credentials.NewTLS(&tls.Config{
		ClientAuth:   tls.RequireAnyClientCert, // this is optional!
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    certPool,
	})

	server := grpc.NewServer(grpc.Creds(creds))
	RegisterGreeterServer(server, new(myGrpcServer))

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Panicf("could not list on %s: %s", port, err)
	}

	if err := server.Serve(lis); err != nil {
		log.Panicf("grpc serve error: %s", err)
	}
}

func doClientWork() {
	certificate, err := tls.LoadX509KeyPair(clientCrt, clientKey)
	if err != nil {
		log.Panicf("could not load client key pair: %s", err)
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(ca)
	if err != nil {
		log.Panicf("could not read ca certificate: %s", err)
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Panicf("failed to append ca certs")
	}
	creds := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: false,         // this is required!
		ServerName:         tlsServerName, // this is required
		Certificates:       []tls.Certificate{certificate},
		RootCAs:            certPool,
	})

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
	log.Printf("doClientwork: %s", r.Message)
}
