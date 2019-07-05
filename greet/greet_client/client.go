package main

import (
	"fmt"
	"log"

	"github.com/katoozi/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello Im Client.")

	// we dont have ssl that why we have to use this => grpc.WithInsecure()
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	fmt.Printf("Created Client: %f", c)
}
