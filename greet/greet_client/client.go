package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/katoozi/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	// we dont have ssl that why we have to use this => grpc.WithInsecure()
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer cc.Close()

	// Create CLient
	c := greetpb.NewGreetServiceClient(cc)

	// do unary rpc call
	// doUnary(c)

	doServerStreaming(c)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("starting to do server streaming rpc...")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "mohammad",
			LastName:  "katoozi",
		},
	}
	stream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while get response: %v", err)
	}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// we reached end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while read from stream: %v", err)
		}
		log.Printf("Response From Server: %v", msg.GetResult())

	}

}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("starting to do unary rpc...")
	request := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Mohammad",
			LastName:  "Katoozi",
		},
	}
	gResponse, err := c.Greet(context.Background(), request)
	if err != nil {
		log.Fatalf("error while calling Greeting: %v", err)
	}
	log.Printf("gRPC Response: %v", gResponse.Result)
}
