// To generate the protobuf code, run below command from outside the "testapp2" dir
// protoc -I testapp2/ testapp2/hello.proto --go_out=plugins=grpc:testapp2

package main

import (
	"log"
	"net"

	"github.com/grpc-test/testapp2"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	listen = "localhost:3002"
)

// server is used to implement testapp1.GreeterServer
type server struct{}

// SayHello implements testapp1.GreeterServer
func (s *server) SayHello(ctx context.Context, in *testapp2.Person) (*testapp2.Greeting, error) {
	return &testapp2.Greeting{Greeting: "Hello from Test Server2 " + in.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", listen)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	testapp2.RegisterGreeterServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
