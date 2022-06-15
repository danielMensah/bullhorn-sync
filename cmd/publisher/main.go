package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/danielMensah/bullhorn-sync-poc/internal/bullhorn"
	"github.com/danielMensah/bullhorn-sync-poc/internal/config"
	"github.com/danielMensah/bullhorn-sync-poc/internal/kafka"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	cfg, err := config.LoadConfig()
	if err != nil {
		log.WithError(err).Fatal("getting configs")
	}

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

	err = messenger.Pub(events)

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
