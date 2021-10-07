package main

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"github.com/danielMensah/bullhorn-sync-poc/internal/bullhorn"
	"github.com/danielMensah/bullhorn-sync-poc/internal/config"
	"github.com/danielMensah/bullhorn-sync-poc/internal/eventbus"
	"github.com/danielMensah/bullhorn-sync-poc/internal/queue"
	"github.com/danielMensah/bullhorn-sync-poc/internal/worker"
	log "github.com/sirupsen/logrus"
)

const workers = 20

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	cfg, err := config.New()
	if err != nil {
		log.WithError(err).Fatal("getting configs")
	}

	bh, err := bullhorn.New(ctx, cfg.BhConfig, cfg.Oauth2Config)
	if err != nil {
		log.WithError(err).Fatal("initializing bullhorn client")
	}

	messages := make(chan queue.Message)
	wg := &sync.WaitGroup{}

	bus := eventbus.New(cfg.Region)
	w := worker.New(bus, bh)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go w.Run(ctx, messages, wg)
	}

	sqs, err := queue.NewSqsConsumer(cfg.Region, messages)
	if err != nil {
		log.WithError(err).Fatal("creating new SQS consumer")
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		<-c

		cancel()
	}()

	sqs.Receive(ctx)
	close(messages)
	wg.Wait()

	log.Warn("shutting down")
}
