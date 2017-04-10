// Copyright © 2017 thingful

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/thingful/device-hub/proto"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := proto.NewHubClient(conn)

	r, err := c.PipeList(context.Background(), &proto.PipeListRequest{})

	fmt.Println(r.Pipes, err)

	for _, p := range r.Pipes {

		fmt.Println(p.State)
		fmt.Println(p.MessageStats)

	}
}
