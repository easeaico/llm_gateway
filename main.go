package main

import (
	"fmt"
	"log"
	"net"

	"github.com/easeaico/llm_mesh/internal/config"
	"github.com/easeaico/llm_mesh/internal/server"
	"github.com/easeaico/llm_mesh/pkg/llm_mesh"
	"google.golang.org/grpc"
)

const (
	port = 5984
)

func main() {
	conf := config.ReadConfigFile("./conf/config.yaml")
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("listening success: %d\n", port)
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	llm_mesh.RegisterChatCompletionServiceServer(grpcServer, server.NewChatCompletionServer(conf))
	grpcServer.Serve(lis)
}
