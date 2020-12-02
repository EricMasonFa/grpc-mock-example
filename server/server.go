package main

import (
	"flag"
	"fmt"
	"grpc-mock-example/internal"
	"grpc-mock-example/pb"
	"log"
	"net"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 10000, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)

	s := &internal.SignalService{}

	pb.RegisterSignalServer(grpcServer, s)
	log.Println(fmt.Sprintf("server is listning on localhost:%d", *port))
	_ = grpcServer.Serve(lis)
}
