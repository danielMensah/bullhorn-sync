package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/danielMensah/bullhorn-sync/internal/auth"
	"github.com/danielMensah/bullhorn-sync/internal/broker"
	"github.com/danielMensah/bullhorn-sync/internal/bullhorn"
	"github.com/danielMensah/bullhorn-sync/internal/config"
	"github.com/danielMensah/bullhorn-sync/internal/poller"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	cfg, err := config.LoadConfig()
	if err != nil {
		log.WithError(err).Fatal("getting configs")
	}

	ctx, cancel := context.WithCancel(context.Background())
	go waitForCancel(cancel)

	producer, err := broker.NewKafkaProducer(cfg.KafkaAddress())
	if err != nil {
		log.WithError(err).Fatal("creating new kafka producer")
	}
	defer producer.Close()

	go producer.MonitorEvents()

	client, err := auth.New(ctx, cfg.BullhornUsername(), cfg.BullhornPassword(), cfg.OauthConfig())
	if err != nil {
		log.WithError(err).Fatal("creating new auth")
	}

	bhClient := bullhorn.New(cfg.BullhornSubscriptionUrl(), cfg.BullhornEntityUrl(), client)
	p := poller.New(bhClient)

	events := make(chan *broker.EventWrapper)
	wg := &sync.WaitGroup{}

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go producer.Produce(ctx, events, wg)
	}

	p.Run(ctx, events)
	wg.Wait()
}

func waitForCancel(cancelFunc context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	cancelFunc()
	log.Info("terminate signal received, exiting...")
}
