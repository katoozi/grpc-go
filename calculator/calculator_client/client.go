package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/katoozi/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculateServiceClient(cc)

	// do unary services
	// doUnary(c)

	// do server streaming rpc
	// doServerStreaming(c)

	// do client streaming
	// doClientStreaming(c)

	// fo BiDi streaming
	doBiDiStreaming(c)
}

func doBiDiStreaming(c calculatorpb.CalculateServiceClient) {
	fmt.Println("starting to do BiDi streaming rpc...")

	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while call FindMaximum: %v", err)
	}

	requestNumbers := []int32{1, 2, 3, 4, 5, 4, 6, 7, 5, 8, 9, 2}

	waitc := make(chan struct{})

	go func() {
		for _, number := range requestNumbers {
			fmt.Printf("Send Request: %v\n", number)
			err := stream.Send(&calculatorpb.FindMaximumRequest{
				Number: number,
			})
			if err != nil {
				log.Fatalf("Error while sending through stream: %v", err)
			}
			time.Sleep(1 * time.Second)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receive from stream: %v", err)
				break
			}
			fmt.Println("Response From Server. Grather Number: ", res.GetResult())
		}
		close(waitc)
	}()

	<-waitc

}

func doClientStreaming(c calculatorpb.CalculateServiceClient) {
	fmt.Println("starting to do client streaming rpc...")

	requests := []*calculatorpb.ComputeAvgRequest{
		&calculatorpb.ComputeAvgRequest{
			Number: 1,
		},
		&calculatorpb.ComputeAvgRequest{
			Number: 2,
		},
		&calculatorpb.ComputeAvgRequest{
			Number: 3,
		},
		&calculatorpb.ComputeAvgRequest{
			Number: 4,
		},
		&calculatorpb.ComputeAvgRequest{
			Number: 5,
		},
	}
	stream, err := c.ComputeAvg(context.Background())
	if err != nil {
		log.Fatalf("Error while calling ComputeAvg: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending Request: %v\n", req)
		err := stream.Send(req)
		if err != nil {
			log.Fatalf("Error while send request: %v", err)
		}
		time.Sleep(1 * time.Second)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while close and receive: %v", err)
	}
	fmt.Printf("Request From Server: %v\n", res)
}

func doServerStreaming(c calculatorpb.CalculateServiceClient) {
	fmt.Println("starting to do server streaming rpc...")

	req := &calculatorpb.PrimeNumberCompositionRequest{
		Number: 40,
	}
	stream, err := c.PrimeNumberComposition(context.Background(), req)
	if err != nil {
		log.Fatalf("Error While call PrimeNumberComposition: %v", err)
	}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while Read Stream: %v", err)
		}
		log.Printf("Response From Server. Prime Number: %d", msg.GetResult())
	}
}

func doUnary(c calculatorpb.CalculateServiceClient) {
	fmt.Println("starting to do unary rpc...")

	doSumCall(c)
	doMulCall(c)
	doDivCall(c)
	doSubCall(c)
}

func doSumCall(c calculatorpb.CalculateServiceClient) {
	request := &calculatorpb.SumRequest{
		FirstNumber:  10,
		SecondNumber: 20,
	}
	res, err := c.Sum(context.Background(), request)
	if err != nil {
		log.Fatalf("Error while calling Sum service: %v", err)
	}
	log.Println("Result: ", res.GetResult())
}

func doDivCall(c calculatorpb.CalculateServiceClient) {
	request := &calculatorpb.DivRequest{
		FirstNumber:  20,
		SecondNumber: 5,
	}
	res, err := c.Div(context.Background(), request)
	if err != nil {
		log.Fatalf("Error while calling Div service: %v", err)
	}
	log.Println("Result: ", res.GetResult())
}

func doMulCall(c calculatorpb.CalculateServiceClient) {
	request := &calculatorpb.MultiplyRequest{
		FirstNumber:  20,
		SecondNumber: 5,
	}
	res, err := c.Multiply(context.Background(), request)
	if err != nil {
		log.Fatalf("Error while calling Multiply service: %v", err)
	}
	log.Println("Result: ", res.GetResult())
}

func doSubCall(c calculatorpb.CalculateServiceClient) {
	request := &calculatorpb.SubRequest{
		FirstNumber:  20,
		SecondNumber: 5,
	}
	res, err := c.Sub(context.Background(), request)
	if err != nil {
		log.Fatalf("Error while calling Multiply service: %v", err)
	}
	log.Println("Result: ", res.GetResult())
}
