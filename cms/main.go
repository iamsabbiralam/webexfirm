package main

import (
	"log"

	"google.golang.org/grpc"
)

func main() {
	_, err := grpc.Dial(":3000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %s", err)
	}
}