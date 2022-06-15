package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/danielMensah/bullhorn-sync-poc/internal/bullhorn"
	"github.com/danielMensah/bullhorn-sync-poc/internal/config"
	"github.com/danielMensah/bullhorn-sync-poc/internal/kafka"
	pb "github.com/danielMensah/bullhorn-sync-poc/internal/proto"
	"github.com/danielMensah/bullhorn-sync-poc/internal/worker"
	log "github.com/sirupsen/logrus"
)

const workers = 20

type queueT interface {
	Close() error
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	ctx, cancel := context.WithCancel(context.Background())

	cfg, err := config.LoadConfig()
	if err != nil {
		log.WithError(err).Fatal("getting configs")
	}

	bh, err := bullhorn.New(ctx, cfg)
	if err != nil {
		log.WithError(err).Fatal("new bullhorn client")
	}

	queue := kafka.NewSubscriber("candidate", cfg.KafkaAddress)
	wkr := worker.New(bh, queue)

	event := make(chan *pb.Event)
	wg := &sync.WaitGroup{}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go wkr.Run(ctx, event, wg)
	}

	go terminate(queue, cancel)

	queue.Sub(ctx, event)

	log.Info("Consumer started")
}

func terminate(queue queueT, cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c

	log.Info("terminate signal received, exiting...")
	if err := queue.Close(); err != nil {
		log.Error("shutting down queue:", err)
	}

	cancel()
}
