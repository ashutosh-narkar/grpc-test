package utils

import (
	"fmt"
	"strings"

	authz "github.com/envoyproxy/data-plane-api/api/auth"
	"github.com/grpc-proxy/proxy"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

func GetDirector(config Config) func(context.Context, string) (context.Context, *grpc.ClientConn, error) {

	return func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		// Copy the inbound metadata explicitly.
		outCtx, _ := context.WithCancel(ctx)
		outCtx = metadata.NewOutgoingContext(outCtx, md.Copy())
		fmt.Println("XXXX OUT Context", outCtx)

		//		for _, backend := range config.Backends {
		//			if strings.HasPrefix(fullMethodName, backend.Filter) {
		//				if config.Verbose {
		//					fmt.Printf("Found: Path: %s > Backend: %s > Backend Name: %s\n", fullMethodName, backend.Backend, backend.BackendName)
		//				}
		//				conn, err := grpc.DialContext(ctx, backend.Backend, grpc.WithCodec(proxy.Codec()),
		//					grpc.WithInsecure())
		//				if err != nil {
		//					fmt.Println("Backend Dialing Error: ", err)
		//				}
		//				return outCtx, conn, err
		//			}
		//		}
		//		if config.Verbose {
		//			fmt.Println("Not found backend for path: ", fullMethodName)
		//		}

		if ok {
			// Decide on which backend to dial
			if strings.HasPrefix(fullMethodName, "/istio.mixer.v1.Mixer/Report") {
				fmt.Println("Report Call received")

				conn, err := grpc.DialContext(ctx, "istio-mixer.istio-system:9091", grpc.WithCodec(proxy.Codec()), grpc.WithInsecure())
				if err != nil {
					fmt.Println("Backend Dialing Error: ", err)
				}
				defer conn.Close()
				return outCtx, conn, err
			} else if strings.HasPrefix(fullMethodName, "/istio.mixer.v1.Mixer/Check") {
				fmt.Println("Check Call received")

				conn, err := grpc.DialContext(ctx, "localhost:50051", grpc.WithCodec(proxy.Codec()), grpc.WithInsecure())
				if err != nil {
					fmt.Println("Backend Dialing Error: ", err)
				}
				defer conn.Close()
				client := authz.NewAuthorizationClient(conn)
				req := authz.CheckRequest{
					Attributes: &authz.AttributeContext{
						Request: &authz.AttributeContext_Request{
							Http: &authz.AttributeContext_HTTPRequest{
								Method: "GET",
								Path:   "/test",
							},
						},
					},
				}

				resp, err := client.Check(ctx, &req)
				if err != nil {
					fmt.Errorf("Check Failed %v", err)
				}
				fmt.Printf("Check response:\n %v\n", resp)

				return outCtx, conn, err
			}
		}
		return nil, nil, grpc.Errorf(codes.Unimplemented, "Unknown method")
	}
}
