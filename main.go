package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/easeaico/llm_mesh/pkg/llm_mesh"
	"github.com/easeaico/llm_mesh/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	var confFile string
	flag.StringVar(&confFile, "f", "config.yaml", "配置文件路径")
	flag.Parse()

	cfg := server.NewConfig(confFile)

	svrconf := cfg.Server
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", svrconf.IP, svrconf.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("listening success: %s:%d\n", svrconf.IP, svrconf.Port)
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	reflection.Register(grpcServer)
	llm_mesh.RegisterChatCompletionServiceServer(grpcServer, server.NewChatCompletionServer(&cfg))

	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg := sync.WaitGroup{}
	wg.Add(1)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		defer signal.Stop(sigCh)

		s := <-sigCh
		log.Printf("got signal %v, attempting graceful shutdown", s)
		cancel()
		grpcServer.GracefulStop()
		wg.Done()
	}()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}

	wg.Wait()
	log.Println("llm mesh server shutdown")
}
