package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	googleGrpc "google.golang.org/grpc"

	"github.com/cronfy/trainer/internal/app/useCase/multiplytask"
	grpcs "github.com/cronfy/trainer/internal/grpcservice"
	pb "github.com/cronfy/trainer/internal/grpcservice/generated"
	"github.com/cronfy/trainer/internal/tools/random"
)

var port = flag.Int("port", 50051, "The server port")

func main() {
	log.Println("Starting gRPC service")

	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := googleGrpc.NewServer()
	pb.RegisterTrainerServer(server, grpcs.Build(multiplytask.New(random.New())))

	log.Println("Listening on port:", *port)
	err = server.Serve(listener)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	} else {
		log.Println("Done")
	}
}
