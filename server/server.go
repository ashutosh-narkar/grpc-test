package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	google_rpc "github.com/gogo/googleapis/google/rpc"
	"github.com/grpc-test/attribute"
	mixerpb "github.com/istio-api/mixer/v1"
	"golang.org/x/net/context"
)

type (
	authzServer struct {
		// Istio's global dictionary
		globalWordList []string
		globalDict     map[string]int32
	}
)

func NewServer() (*authzServer, error) {
	list := attribute.GlobalList()
	globalDict := make(map[string]int32, len(list))
	for i := 0; i < len(list); i++ {
		globalDict[list[i]] = int32(i)
	}

	return &authzServer{
		globalWordList: list,
		globalDict:     globalDict,
	}, nil
}

func (as *authzServer) Check(ctx context.Context, req *mixerpb.CheckRequest) (*mixerpb.CheckResponse, error) {
	// Get the http request info from Istio's attributes
	protoBag := attribute.NewProtoBag(&req.Attributes, as.globalDict, as.globalWordList)
	checkBag := attribute.GetMutableBag(protoBag)
	defer checkBag.Done()

	fmt.Printf("\nAttributes are:\n%+v\n", checkBag.DebugString())

	method, _ := checkBag.Get("request.method")
	path, _ := checkBag.Get("request.path")
	headers, _ := checkBag.Get("request.headers")
	headersMap, ok := headers.(attribute.StringMap)
	if !ok {
		fmt.Println("Could not convert empty interface to type \"attribute.StringMap\"")
	}

	auth, found := headersMap.Get("authorization")
	if !found {
		fmt.Println("Cound not find Auth info in request header")
		//TODO: What if auth info not provided ? Return ?
	}
	userAuth := strings.Split(auth, " ")[1]

	// OPA input
	fmt.Printf("Request to OPA:\nMethod %s\nPath %s\nUser %s\n", method, path, userAuth)
	input := make(map[string]string)
	input["method"] = method.(string)
	input["path"] = path.(string)
	input["auth"] = userAuth

	body := make(map[string]interface{})
	body["input"] = input

	inputBytes, oerr := json.Marshal(body)
	if oerr != nil {
		fmt.Errorf("JSON Encoding error %v", oerr)
	}

	// ask OPA for a policy decision
	opa_server_url := "http://localhost:8181/v1/data/istio/http/allow"
	resp, oerr := http.Post(opa_server_url, "application/json", bytes.NewBuffer(inputBytes))
	if oerr != nil {
		fmt.Errorf("HTTP request error %v", oerr)
	}

	// handle OPA response
	var result map[string]interface{}
	oerr = json.NewDecoder(resp.Body).Decode(&result)
	if oerr != nil {
		fmt.Errorf("JSON Decoding error is %v", oerr)
	}

	policy_result, ok := result["result"].(bool)
	if !ok {
		fmt.Errorf("Type assertion error")
	}

	var status int32
	if policy_result {
		log.Printf("OPA: Operation allowed\n")
		status = int32(google_rpc.OK)
	} else {
		log.Printf("OPA: Operation not allowed\n")
		status = int32(google_rpc.PERMISSION_DENIED)
	}

	response := &mixerpb.CheckResponse{
		Precondition: mixerpb.CheckResponse_PreconditionResult{
			Status: google_rpc.Status{Code: status},
		},
	}

	return response, nil
}

func (as *authzServer) Report(ctx context.Context, req *mixerpb.ReportRequest) (*mixerpb.ReportResponse, error) {
	return new(mixerpb.ReportResponse), nil
}
