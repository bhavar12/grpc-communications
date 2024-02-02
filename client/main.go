package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "project-engage-grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address     = "localhost:8080"
	defaultName = "world"
)

func main() {
	con, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("client faild to connect to server %v", err)
	}

	defer con.Close()

	client := pb.NewGreeterClient(con)

	names := &pb.NamesList{
		Names: []string{"Akhil", "Alice", "Bob"},
	}

	callSayHello(client, defaultName)
	callSayHelloServerStream(client, names)
	callSayHelloClientStream(client, names)
	callSayHelloBidirctiona(client, names)
}

// unary rpc request
// normal rest api working req resp. no streaming
func callSayHello(client pb.GreeterClient, name string) {
	log.Printf("unary rpc request started")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}
	log.Printf("%s", res.Message)
	log.Printf("unary rpc request ended")
}

// server side streaming request
// here client send req and server will stream the resp
func callSayHelloServerStream(client pb.GreeterClient, names *pb.NamesList) {
	log.Printf("Server Streaming started")
	stream, err := client.SayHelloServerStreaming(context.Background(), names)
	if err != nil {
		log.Fatalf("Could not send names: %v", err)
	}

	for {
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while streaming %v", err)
		}
		log.Println(message)
	}

	log.Printf("Server Streaming ended")
}

// cleint side streaming request
// client send stream of req and server will stream single resp
func callSayHelloClientStream(client pb.GreeterClient, names *pb.NamesList) {
	log.Printf("Client Streaming started")
	stream, err := client.SayHelloClientStreaming(context.Background())
	if err != nil {
		log.Fatalf("Could not send names: %v", err)
	}

	for _, name := range names.Names {
		req := &pb.HelloRequest{
			Name: name,
		}
		if err := stream.Send(req); err != nil {
			log.Fatalf("Error while sending %v", err)
		}
		log.Printf("Sent request with name: %s", name)
		time.Sleep(2 * time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving %v", err)
	}
	log.Printf("%v", res.Messages)
	log.Printf("Client Streaming ended")
}

// Bi Directional Steaming
// streaming of req and resp eg. stock market
func callSayHelloBidirctiona(client pb.GreeterClient, names *pb.NamesList) {
	log.Printf("Bi Dirctional Streaming started")
	stream, err := client.SayHelloBidirectionalStreaming(context.Background())
	if err != nil {
		log.Fatalf("error while streaming %v", err)
	}
	waitChan := make(chan struct{})
	go func() {
		for {
			message, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while streaming %v", err)
			}

			log.Printf("Message.. %v", message)
		}
		close(waitChan)
	}()

	for _, val := range names.Names {
		req := &pb.HelloRequest{Name: val}
		if err := stream.Send(req); err != nil {
			log.Fatalf("error while streaming req %v", err)
		}
		// adding sleep to avoid overlapping
		time.Sleep(time.Second * 1)
	}

	stream.CloseSend()
	<-waitChan
	log.Printf("Bi Dirctional Streaming ended")
}
