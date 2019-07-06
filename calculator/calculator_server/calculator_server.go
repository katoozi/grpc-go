package main

import (
	"context"
	"log"
	"net"

	"github.com/katoozi/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	result := req.FirstNumber + req.SecondNumber
	return &calculatorpb.SumResponse{
		Result: result,
	}, nil
}

func (s *server) Div(ctx context.Context, req *calculatorpb.DivRequest) (*calculatorpb.DivResponse, error) {
	result := req.FirstNumber / req.SecondNumber
	return &calculatorpb.DivResponse{
		Result: result,
	}, nil
}

func (s *server) Multiply(ctx context.Context, req *calculatorpb.MultiplyRequest) (*calculatorpb.MultiplyResponse, error) {
	result := req.FirstNumber * req.SecondNumber
	return &calculatorpb.MultiplyResponse{
		Result: result,
	}, nil
}

func (s *server) Sub(ctx context.Context, req *calculatorpb.SubRequest) (*calculatorpb.SubResponse, error) {
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
