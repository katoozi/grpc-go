package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/katoozi/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet Function Was invoked with: %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	result := &greetpb.GreetResponse{
		Result: fmt.Sprintf("Hello %s\n", firstName),
	}
	return result, nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serv: %v", err)
	}
}
