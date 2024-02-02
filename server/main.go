package main

import (
	"context"
	"log"
	"net"
	pb "project-engage-grpc/proto"

	"io"
	"time"

	"google.golang.org/grpc"
)

type server struct {
	pb.GreeterServer
}

// define the address
const (
	address = ":8080"
)

func main() {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen server %v", err)
	}
	grpcSer := grpc.NewServer()
	pb.RegisterGreeterServer(grpcSer, &server{})
	log.Printf("Server started %v", lis.Addr())
	if err := grpcSer.Serve(lis); err != nil {
		log.Fatalf("Failed to start server %v", err)
	}
}

// unary communication
// will not stream . just send the resp
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello, I am " + in.Name}, nil
}

// server side streaming handler
// will get the single req and will stream the resp
func (s *server) SayHelloServerStreaming(req *pb.NamesList, stream pb.Greeter_SayHelloServerStreamingServer) error {
	log.Printf("Got request with names : %v", req.Names)
	for _, name := range req.Names {
		res := &pb.HelloReply{
			Message: "Hello " + name,
		}
		if err := stream.Send(res); err != nil {
			return err
		}
		// 2 second delay to simulate a long running process
		time.Sleep(2 * time.Second)
	}
	return nil
}

// client side streaming handler
// will get the stream of req and will send single resp
func (s *server) SayHelloClientStreaming(stream pb.Greeter_SayHelloClientStreamingServer) error {
	var messages []string
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.MessagesList{Messages: messages})
		}
		if err != nil {
			return err
		}
		log.Printf("Got request with name : %v", req.Name)
		messages = append(messages, "Hello "+req.Name)
	}
}

func (s *server) SayHelloBidirectionalStreaming(stream pb.Greeter_SayHelloBidirectionalStreamingServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		log.Printf("Got request with name : %v", req.Name)
		res := &pb.HelloReply{
			Message: "Hello.. " + req.Name,
		}

		if err := stream.Send(res); err != nil {
			return err
		}
	}
}
