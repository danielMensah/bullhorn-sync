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

	producer, err := broker.NewKafkaProducer(cfg.KafkaAddress)
	if err != nil {
		log.WithError(err).Fatal("creating new kafka producer")
	}

	bhClient, err := bullhorn.New(ctx, cfg)
	if err != nil {
		log.WithError(err).Fatal("new bullhorn client")
	}
	p := poller.New(bhClient)

	events := make(chan *broker.EventWrapper)
	wg := &sync.WaitGroup{}

	for i := 0; i < 20; i++ {
		wg.Add(2)
		go producer.MonitorEvents(wg)
		go producer.Produce(ctx, events, wg)
	}

	p.Run(ctx, events)
	wg.Wait()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c

		cancel()
		log.Info("terminate signal received, exiting...")
	}()
}
