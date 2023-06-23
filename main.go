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

func main() {
	conf := config.ReadConfigFile("./conf/config.yaml")
	svrconf := conf.Server
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", svrconf.IP, svrconf.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("listening success: %s:%d\n", svrconf.IP, svrconf.Port)
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	llm_mesh.RegisterChatCompletionServiceServer(grpcServer, server.NewChatCompletionServer(&conf))
	grpcServer.Serve(lis)
}
