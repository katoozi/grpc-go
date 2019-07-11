package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/katoozi/grpc-go-course/blog/blogpb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var collection *mongo.Collection

type bLogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

type server struct{}

func (*server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	log.Printf("CreateBlog rpc invoked with: %v\n", req)

	blog := req.GetBlog()

	data := bLogItem{
		AuthorID: blog.GetAuthorId(),
		Content:  blog.GetContent(),
		Title:    blog.GetTitle(),
	}

	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal, "Cannot Parse objectid")
	}

	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       oid.Hex(),
			Content:  data.Content,
			Title:    data.Title,
			AuthorId: data.AuthorID,
		},
	}, nil
}

func (*server) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	log.Printf("ReadBlog rpc invoked with: %v\n", req)

	blogID := req.GetBlogId()

	oid, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Cannot Parse ID: %v", err)
	}
	data := &bLogItem{}
	filter := bson.M{"_id": oid}

	res := collection.FindOne(context.Background(), filter)
	if err = res.Decode(&data); err != nil {
		return nil, status.Errorf(codes.NotFound, "blog not found in db: %v", res)
	}

	return &blogpb.ReadBlogResponse{
		Blog: &blogpb.Blog{
			Id:       data.ID.Hex(),
			AuthorId: data.AuthorID,
			Title:    data.Title,
			Content:  data.Content,
		},
	}, nil
}

func (*server) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	fmt.Printf("UpdateBlog invoked with: %v\n", req)

	blog := req.GetBlog()
	oid, err := primitive.ObjectIDFromHex(blog.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Cannot parse id: %v", err)
	}

	data := &bLogItem{}
	filter := bson.M{"_id": oid}

	res := collection.FindOne(context.Background(), filter)
	if err = res.Decode(&data); err != nil {
		return nil, status.Errorf(codes.NotFound, "blog not found in db: %v", res)
	}

	// update internal struct
	data.AuthorID = blog.GetAuthorId()
	data.Content = blog.GetContent()
	data.Title = blog.GetTitle()

	_, updateErr := collection.ReplaceOne(context.Background(), filter, data)
	if updateErr != nil {
		return nil, status.Errorf(codes.Internal, "Cannot update: %v", updateErr)
	}

	return &blogpb.UpdateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       data.ID.Hex(),
			AuthorId: data.AuthorID,
			Title:    data.Title,
			Content:  data.Content,
		},
	}, nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("Start Blog Service. Listning...")

	// make mongodb configurations. connect to mongodb
	fmt.Println("Start Mongodb client connection...")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://@localhost:27017"))
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.TODO(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	collection = client.Database("mydb").Collection("blog")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		fmt.Println("Start Server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// wait for control + c to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// block until control + c
	<-ch
	fmt.Println("Stopping The Server")
	// stopping the grpc server
	s.Stop()
	fmt.Println("Closing the Listener")
	// closing the net package listner
	lis.Close()
	fmt.Println("Close Mongodb Client")
	client.Disconnect(context.TODO())
	fmt.Println("End of App.")
}
