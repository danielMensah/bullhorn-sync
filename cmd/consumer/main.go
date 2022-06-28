package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/danielMensah/bullhorn-sync-poc/internal/bullhorn"
	"github.com/danielMensah/bullhorn-sync-poc/internal/config"
	"github.com/danielMensah/bullhorn-sync-poc/internal/kafka"
	pb "github.com/danielMensah/bullhorn-sync-poc/internal/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	cfg, err := config.LoadConfig()
	if err != nil {
		log.WithError(err).Fatal("getting configs")
	}

	conn, err := grpc.Dial(cfg.RPCAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v\n", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	messenger, err := kafka.NewPublisher(ctx, cfg.KafkaAddress)
	if err != nil {
		log.Fatal(err)
	}

	bhClient, err := bullhorn.New(ctx, cfg)
	if err != nil {
		log.WithError(err).Fatal("new bullhorn client")
	}

	events, err := bhClient.GetEvents()
	if err != nil {
		log.WithError(err).Fatal("getting entities")
	}

	c := pb.NewPublisherServiceClient(conn)
	wrappedEvent := &pb.EventsWrapper{
		Events: events,
	}

	_, err = c.Publish(ctx, wrappedEvent)
	if err != nil {
		log.WithError(err).Fatal("sending wrapped events to publisher")
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Info("terminate signal received, exiting...")
		if err := messenger.Close(); err != nil {
			log.Error("shutting down:", err)
		}
		cancel()
	}()

	log.Info("Producer started")
}
