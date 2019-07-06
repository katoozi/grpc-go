package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/katoozi/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Sum Service Invoked: %v\n", req)
	result := req.FirstNumber + req.SecondNumber
	return &calculatorpb.SumResponse{
		Result: result,
	}, nil
}

func (*server) Div(ctx context.Context, req *calculatorpb.DivRequest) (*calculatorpb.DivResponse, error) {
	fmt.Printf("Div Service Invoked: %v\n", req)
	result := req.FirstNumber / req.SecondNumber
	return &calculatorpb.DivResponse{
		Result: result,
	}, nil
}

func (*server) Multiply(ctx context.Context, req *calculatorpb.MultiplyRequest) (*calculatorpb.MultiplyResponse, error) {
	fmt.Printf("Multiply Service Invoked: %v\n", req)
	result := req.FirstNumber * req.SecondNumber
	return &calculatorpb.MultiplyResponse{
		Result: result,
	}, nil
}

func (*server) Sub(ctx context.Context, req *calculatorpb.SubRequest) (*calculatorpb.SubResponse, error) {
	fmt.Printf("Sub Service Invoked: %v\n", req)
	result := req.FirstNumber - req.SecondNumber
	return &calculatorpb.SubResponse{
		Result: result,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculateServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
