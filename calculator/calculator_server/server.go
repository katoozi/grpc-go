package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"time"

	"github.com/katoozi/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

// return list of primes less than N
func sieveOfEratosthenes(N int) (primes []int) {
	b := make([]bool, N)
	for i := 2; i < N; i++ {
		if b[i] == true {
			continue
		}
		primes = append(primes, i)
		for k := i * i; k < N; k += i {
			b[k] = true
		}
	}
	return
}

func (*server) PrimeNumberComposition(req *calculatorpb.PrimeNumberCompositionRequest, stream calculatorpb.CalculateService_PrimeNumberCompositionServer) error {
	fmt.Printf("PrimeNumberComposition Service Invoked: %v\n", req)
	primes := sieveOfEratosthenes(int(req.GetNumber()))
	for _, p := range primes {
		res := &calculatorpb.PrimeNumberCompositionResponse{
			Result: int32(p),
		}
		stream.Send(res)
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (*server) ComputeAvg(stream calculatorpb.CalculateService_ComputeAvgServer) error {
	fmt.Println("ComputeAvg Service Invoked")

	var result int32
	var count int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&calculatorpb.ComputeAvgResponse{
				Result: result / count,
			})
		}
		if err != nil {
			log.Fatalf("Error while close stream: %v", err)
		}
		count++
		result += req.GetNumber()
	}
}

func (*server) FindMaximum(stream calculatorpb.CalculateService_FindMaximumServer) error {
	fmt.Println("FindMaximum Service inviked")

	var previusNumbers int32
	count := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
			return nil
		}
		number := req.GetNumber()
		count++
		if count == 1 {
			err := stream.Send(&calculatorpb.FindMaximumResponse{
				Result: number,
			})
			if err != nil {
				log.Fatalf("Error while send through stream: %v", err)
				return err
			}
		}
		numberResult := float64(number)
		if numberResult > float64(previusNumbers) {
			previusNumbers = int32(number)
		}
		err = stream.Send(&calculatorpb.FindMaximumResponse{
			Result: previusNumbers,
		})
		if err != nil {
			log.Fatalf("Error while send through stream: %v", err)
			return err
		}
	}
}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Printf("Square Root RPC invoked with: %v\n", req)

	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Received Negetive number: %v\n",
			number,
		)
	}
	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
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
