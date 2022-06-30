package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/danielMensah/bullhorn-sync-poc/internal/bullhorn"
	"github.com/danielMensah/bullhorn-sync-poc/internal/config"
	"github.com/danielMensah/bullhorn-sync-poc/internal/poller"
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

	ctx, cancel := context.WithCancel(context.Background())
	conn, err := grpc.DialContext(ctx, cfg.RPCAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v\n", err)
	}
	defer conn.Close()

	bhClient, err := bullhorn.New(ctx, cfg)
	if err != nil {
		log.WithError(err).Fatal("new bullhorn client")
	}
	p := poller.New(bhClient)

	records := make(chan []*pb.Event)
	wg := &sync.WaitGroup{}

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go sendToPublisher(ctx, conn, records, wg)
	}

	p.Run(records)
	wg.Wait()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c

		close(p.Done)

		log.Info("terminate signal received, exiting...")
		cancel()
	}()

	log.Info("Producer started")
}

func sendToPublisher(ctx context.Context, conn *grpc.ClientConn, records <-chan []*pb.Event, wg *sync.WaitGroup) {
	defer wg.Done()

	c := pb.NewPublisherServiceClient(conn)
	wrappedEvent := &pb.EventsWrapper{
		Events: <-records,
	}

	_, err := c.Publish(ctx, wrappedEvent)
	if err != nil {
		log.WithError(err).Fatal("sending wrapped events to publisher")
	}
}
