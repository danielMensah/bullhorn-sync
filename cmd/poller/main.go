package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/danielMensah/bullhorn-sync-poc/internal/broker"
	"github.com/danielMensah/bullhorn-sync-poc/internal/bullhorn"
	"github.com/danielMensah/bullhorn-sync-poc/internal/config"
	"github.com/danielMensah/bullhorn-sync-poc/internal/poller"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	cfg, err := config.LoadConfig()
	if err != nil {
		log.WithError(err).Fatal("getting configs")
	}

	ctx, cancel := context.WithCancel(context.Background())

	publisher, err := broker.NewKafkaPublisher(ctx, cfg.KafkaAddress)
	if err != nil {
		log.WithError(err).Fatal("creating new kafka publisher")
	}

	bhClient, err := bullhorn.New(ctx, cfg)
	if err != nil {
		log.WithError(err).Fatal("new bullhorn client")
	}
	p := poller.New(bhClient)

	events := make(chan *broker.EventWrapper)
	wg := &sync.WaitGroup{}

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go publisher.Publish(ctx, events, wg)
	}

	p.Run(events)
	wg.Wait()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c

		close(p.Done)

		log.Info("terminate signal received, exiting...")
		cancel()
	}()
}
