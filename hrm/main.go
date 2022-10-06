package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	listen, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}

	grpcServer := grpc.NewServer()
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}