package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/katoozi/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	// do server streaming call
	// doServerStreaming(c)

	// do client streaming
	// doClientStreaming(c)

	// do BiDi streaming
	// doBiDiStreaming(c)

	// do unary with deadline
	doUnaryGreetWithDeadline(c, 5*time.Second) // should complete
	doUnaryGreetWithDeadline(c, 1*time.Second) // should timeout
}

func doUnaryGreetWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	fmt.Println("UnaryGreetWithDeadline RPC invoked")

	request := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Mohammad",
			LastName:  "Katoozi",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	gResponse, err := c.GreetWithDeadline(ctx, request)
	if err != nil {
		statusError, ok := status.FromError(err)
		if ok {
			if statusError.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit! Deadline was exceeded")
			} else {
				fmt.Printf("unexpected error: %v\n", statusError)
			}
		} else {

		}
		log.Fatalf("error while calling GreetWithDeadline: %v\n", err)
	}
	log.Printf("gRPC Response: %v", gResponse.Result)
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("starting to do BiDi streaming rpc...")

	// we create a stream by invoking the client
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while calling GreetEveryone rpc: %v\n", err)
	}

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "mohammad",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "ali",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "reza",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "ghasem",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "amin",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "mahyar",
			},
		},
	}

	waitc := make(chan struct{})

	// we send a bunch of messages to the client (go routine)
	go func() {
		// function to send a bunch of messages
		for _, request := range requests {
			fmt.Printf("Send Request: %v\n", request)
			err := stream.Send(request)
			if err != nil {
				log.Fatalf("Error while sending requests: %v", err)
			}
			time.Sleep(1 * time.Second)
		}
		stream.CloseSend()
	}()

	// we receive a bunch of messages from the server (go routine)
	go func() {
		// function to receive a bunch of messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while reading stream: %v", err)
				break
			}
			fmt.Printf("Response From Server: %v\n", res)
		}
		close(waitc)
	}()

	// block until everything in done
	<-waitc
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("starting to do client streaming rpc...")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "mohammad",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "ali",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "reza",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "ghasem",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "amin",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "mahyar",
			},
		},
	}
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while create stream: %v", err)
	}

	for _, request := range requests {
		fmt.Printf("Sending Request: %v\n", request)
		err := stream.Send(request)
		if err != nil {
			log.Fatalf("Error while sending requests; %v", err)
		}
		time.Sleep(1 * time.Second)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while closing stream:%v", err)
	}
	fmt.Printf("Result From Server: %v\n", res)
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
