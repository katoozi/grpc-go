package main

import (
	"context"
	"fmt"
	"log"

	"github.com/katoozi/grpc-go-course/blog/blogpb"
	"google.golang.org/grpc"
)

func main() {
	// set log flags for show more info in logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error while connect to server: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	// create new blog and after that reading blog
	// first CreateBlog and then ReadBlog rpc will called
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

	fmt.Println("ReadBlog Section is start...")
	reqReadBlog := &blogpb.ReadBlogRequest{
		BlogId: resp.GetBlog().GetId(),
	}

	resReadBlog, err := c.ReadBlog(context.Background(), reqReadBlog)
	if err != nil {
		log.Fatalf("Error while reading blog: %v\n", err)
	}
	fmt.Printf("ReadBlog rpc result from server: %v\n", resReadBlog)

	// read blog with wrong id for testing
	reqReadBlog = &blogpb.ReadBlogRequest{
		BlogId: "5d2731ccb815173786273e1f",
	}

	resReadBlog, err = c.ReadBlog(context.Background(), reqReadBlog)
	if err != nil {
		log.Fatalf("Error while reading blog: %v\n", err)
	}
	fmt.Printf("ReadBlog rpc result from server: %v\n", resReadBlog)
}
