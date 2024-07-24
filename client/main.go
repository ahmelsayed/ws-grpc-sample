package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/ahmelsayed/ws-grpc-sample/hello"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cc, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer cc.Close()

	client := pb.NewHelloServiceClient(cc)
	request := &pb.HelloRequest{Name: "Foo"}

	resp, err := client.Hello(context.Background(), request)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Receive grpc response >> %s ", resp.Message)
}
