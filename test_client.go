package main

import (
	"log"
	"time"

	"github.com/grpc-test/testapp1"
	"github.com/grpc-test/testapp2"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051" // grpc proxy address
	defaultName = "Ash"
)

func main() {
	// Set up a connection to the gRPC proxy server
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Test Server 1
	c := testapp1.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &testapp1.Person{Name: defaultName})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetGreeting())

	// Test Server 2
	d := testapp2.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	s, err := d.SayHello(ctx, &testapp2.Person{Name: defaultName})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", s.GetGreeting())
}
