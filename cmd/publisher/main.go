package main

import (
	"net"

	"github.com/danielMensah/bullhorn-sync-poc/internal/config"
	pb "github.com/danielMensah/bullhorn-sync-poc/internal/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var address = "0.0.0.0:50051"

type Server struct {
	pb.PublisherServiceServer
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.WithError(err).Fatal("getting configs")
	}

	listener, err := net.Listen("tcp", cfg.RPCAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}

	log.Printf("gRPC server listening on %s\n", address)

	s := grpc.NewServer()
	pb.RegisterPublisherServiceServer(s, &Server{})

	if err = s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}
