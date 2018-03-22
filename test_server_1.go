// To generate the protobuf code, run below command from outside the "testapp1" dir
// protoc -I testapp1/ testapp1/hello.proto --go_out=plugins=grpc:testapp1

package main

import (
	"log"
	"net"

	"github.com/grpc-test/testapp1"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	listen = "localhost:3001"
)

// server is used to implement testapp1.GreeterServer
type server struct{}

// SayHello implements testapp1.GreeterServer
func (s *server) SayHello(ctx context.Context, in *testapp1.Person) (*testapp1.Greeting, error) {
	return &testapp1.Greeting{Greeting: "Hello from Test Server1 " + in.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", listen)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	testapp1.RegisterGreeterServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
