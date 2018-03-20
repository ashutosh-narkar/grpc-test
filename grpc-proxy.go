// A gRPC reverse proxy

package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	authz "github.com/envoyproxy/data-plane-api/api/auth"
	"github.com/grpc-proxy/proxy"
	"github.com/grpc-test/server"
	"github.com/grpc-test/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	configurationFile := "./config.json"

	args := os.Args[1:]
	if len(args) > 0 {
		configurationFile = args[0]
	}

	config := utils.GetConfiguration(configurationFile)

	listen := ":50051"
	if config.Listen != "" {
		listen = config.Listen
	}

	lis, err := net.Listen("tcp", listen)
	defer lis.Close()

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Printf("Starting gRPC Proxy on %s\n", listen)

	grpcServer := GetGrpcServer(config)
	authzServer, err := server.NewServer()
	if err != nil {
		log.Fatalf("Unable to start server %v", err)
	}

	authz.RegisterAuthorizationServer(grpcServer, authzServer)
	reflection.Register(grpcServer)

	// Run gRPC server on separate goroutine
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Use a buffered channel so we don't miss any signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Block until a signal is received.
	s := <-c
	fmt.Println("Got signal:", s)
}

func GetGrpcServer(config utils.Config) *grpc.Server {
	var grpcOptions []grpc.ServerOption

	// Add custom codec and handler
	grpcOptions = append(grpcOptions, grpc.CustomCodec(proxy.Codec()),
		grpc.UnknownServiceHandler(proxy.TransparentHandler(utils.GetDirector(config))))

	return grpc.NewServer(grpcOptions...)
}
