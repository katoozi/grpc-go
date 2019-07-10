package main

import (
	"context"
	"fmt"
	"log"

	"github.com/katoozi/grpc-go-course/blog/blogpb"
	"google.golang.org/grpc"
)

func main() {

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error while connect to server: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	// create new blog
	doCreateBlog(c)
}

func doCreateBlog(c blogpb.BlogServiceClient) {
	fmt.Println("Starting the CreateBlog rcp...")

	request := &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			Title:    "Download From Google",
			AuthorId: "k2527806",
			Content:  "How can i do that ?",
		},
	}

	resp, err := c.CreateBlog(context.Background(), request)
	if err != nil {
		log.Fatalf("Error while calling CreateBlog: %v", err)
	}
	log.Printf("Result From Server: %v", resp)
}
