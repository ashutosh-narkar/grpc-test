package server

import (
	"fmt"

	authz "github.com/envoyproxy/data-plane-api/api/auth"
	"golang.org/x/net/context"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/genproto/googleapis/rpc/status"
)

type authzServer struct{}

func NewServer() (*authzServer, error) {
	return &authzServer{}, nil
}

func (as *authzServer) Check(ctx context.Context, req *authz.CheckRequest) (*authz.CheckResponse, error) {
	fmt.Printf("Check Context %v\n", ctx)
	fmt.Printf("Check Request %v\n", req)

	// Check with OPA
	resp := authz.CheckResponse{Status: &status.Status{Code: int32(code.Code_OK)}}
	return &resp, nil
}
