package main

import (
	"context"
	"fmt"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/katoozi/grpc-go-course/blog/blogpb"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/testdata"
)

func main() {
	// set log flags for show more info in logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	perRPC := oauth.NewOauthAccess(fetchToken())
	creds, err := credentials.NewClientTLSFromFile(testdata.Path("ca.pem"), "x.test.youtube.com")
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}
	opts := []grpc.DialOption{
		// In addition to the following grpc.DialOption, callers may also use
		// the grpc.CallOption grpc.PerRPCCredentials with the RPC invocation
		// itself.
		// See: https://godoc.org/google.golang.org/grpc#PerRPCCredentials
		grpc.WithPerRPCCredentials(perRPC),
		// grpc.WithInsecure(),
		grpc.WithTransportCredentials(creds),
	}
	cc, err := grpc.Dial("localhost:50051", opts...)
	if err != nil {
		log.Fatalf("Error while connect to server: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	// create new blog and after that reading blog
	// first CreateBlog and then ReadBlog rpc will called
	// doCreateBlog(c)

	// get list og blogs
	doListBlog(c)
}

// fetchToken simulates a token lookup and omits the details of proper token
// acquisition. For examples of how to acquire an OAuth2 token, see:
// https://godoc.org/golang.org/x/oauth2
func fetchToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken: "123456789",
	}
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
	// reqReadBlog = &blogpb.ReadBlogRequest{
	// 	BlogId: "5d2731ccb815173786273e1f",
	// }

	// resReadBlog, err = c.ReadBlog(context.Background(), reqReadBlog)
	// if err != nil {
	// 	log.Fatalf("Error while reading blog: %v\n", err)
	// }
	// fmt.Printf("ReadBlog rpc result from server: %v\n", resReadBlog)

	// updateBlog

	newBlog := &blogpb.UpdateBlogRequest{
		Blog: &blogpb.Blog{
			Id:       resp.GetBlog().GetId(),
			Title:    "Changes Title",
			AuthorId: "Change Author",
			Content:  "Change Content",
		},
	}

	updateResp, updateErr := c.UpdateBlog(context.Background(), newBlog)
	if updateErr != nil {
		log.Fatalf("Error From Server: %v", updateErr)
	}

	fmt.Printf("Blog was updated:%v\n", updateResp)

	// delete blog rpc
	deleteReq := &blogpb.DeleteBlogRequest{
		BlogId: updateResp.GetBlog().GetId(),
		// BlogId: "5d274a89aba540a5a277f86f", // wrong test for test errors
	}
	delResp, delErr := c.DeleteBlog(context.Background(), deleteReq)
	if delErr != nil {
		log.Fatalf("Error while DeleteBlog rpc: %v\n", delErr)
	}
	fmt.Printf("Object with %s id is deleted\n", delResp.GetBlogId())
}

func doListBlog(c blogpb.BlogServiceClient) {
	fmt.Println("Start ListBlog server stream rpc...")

	resp, err := c.ListBlog(context.Background(), &empty.Empty{})
	if err != nil {
		log.Fatalf("Error while calling ListBlog: %v\n", err)
	}
	// count := 0
	// for {
	// 	resp, err := stream.Recv()
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	if err != nil {
	// 		log.Fatalf("Error while reading stream: %v\n", err)
	// 	}
	// 	fmt.Printf("Response From Server %d: %v\n", count, resp.GetBlog())
	// 	count++
	// }
	fmt.Printf("Response From Server: %v\n", resp)
}
