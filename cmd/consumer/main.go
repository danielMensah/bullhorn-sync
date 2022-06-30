package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/danielMensah/bullhorn-sync-poc/internal/config"
	"github.com/danielMensah/bullhorn-sync-poc/internal/consumer"
	"github.com/danielMensah/bullhorn-sync-poc/internal/kafka"
	"github.com/danielMensah/bullhorn-sync-poc/internal/repository"
	"github.com/danielMensah/bullhorn-sync-poc/internal/wpool"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	ctx, cancel := context.WithCancel(context.Background())

	cfg, err := config.LoadConfig()
	if err != nil {
		log.WithError(err).Fatal("getting configs")
	}

	repo, err := repository.NewCassandraDB(cfg)
	if err != nil {
		log.WithError(err).Fatal("creating new repository")
	}

	kafkaConsumer := kafka.NewConsumer("candidate", cfg.KafkaAddress)
	entityConsumer := consumer.New(repo)

	pool := wpool.New(20)
	pool.Run(ctx)

	entities := make(chan *kafka.EventWrapper)
	wg := &sync.WaitGroup{}

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go pool.AddJob(ctx, entities, entityConsumer)
	}

	kafkaConsumer.Consume(ctx, entities)
	wg.Wait()
	close(entities)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c

		if err = kafkaConsumer.Close(); err != nil {
			log.WithError(err).Fatal("closing kafka consumer")
		}

		log.Info("terminate signal received, exiting...")
		cancel()
	}()
}
