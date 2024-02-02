# Basic Go gRPC Server and Client

This is a basic gRPC server and client written in Go. It is based on the

We have implemented a simple gRPC server and client with the following functionality:
- unary RPC
- server-side streaming RPC
- client-side streaming RPC
- bidirectional streaming RPC

# Setting up a gRPC-Go project
1. Create a new directory for your project and cd into it

```bash
mkdir project-engage-grpc
cd project-engage-grpc
mkdir client server proto
```

2. Initialize a Go module

```bash
go mod init project-engage-grpc

```

3. Download protobuf dependancy 

```bash
go get google.golang.org/protobuf/proto
```

4. Install protoc and "protoc" command works on your terminal
   ```Download the setup based on your OS from this link```
   https://github.com/protocolbuffers/protobuf/releases

   ```And setup the $PATH env varaible till bin```
   ```Verifi the protoc cmd working fine or not from your terminal```

5. Installing the gRPC Go plugin

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28

go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

```

6. Create the proto file with the required services and messages in the proto directory

7. Generate .pb.go files from the proto file

depending on what path you mention in your greet.proto file, you will either run this - 


```bash
protoc --go_out=. --go-grpc_out=. proto/greet.proto
```
8. Create the server and client directories and create the main.go files with necessary controllers and services


# Running the application

1. Run the server

```bash
go run server/main.go
```

2. Run the client

```bash
go run client/main.go
```