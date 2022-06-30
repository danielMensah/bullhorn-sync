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
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	cfg, err := config.LoadConfig()
	if err != nil {
		log.WithError(err).Fatal("loading config")
	}

	ctx, cancel := context.WithCancel(context.Background())

	consumer := kafka.NewConsumer("poller_events", cfg.KafkaAddress)

	publisher, err := kafka.NewPublisher(ctx, cfg.KafkaAddress)
	if err != nil {
		log.WithError(err).Fatal("creating new publisher")
	}

	bhClient, err := bullhorn.New(ctx, cfg)
	if err != nil {
		log.WithError(err).Fatal("creating new bullhorn")
	}

	publisherEvents := make(chan *kafka.EventWrapper)
	pollerEvents := make(chan *kafka.EventWrapper)

	wg := &sync.WaitGroup{}

	for i := 0; i < 20; i++ {
		wg.Add(2)
		go processEntity(ctx, bhClient, pollerEvents, publisherEvents, wg)
		go publisher.Publish(ctx, publisherEvents, wg)
	}

	consumer.Consume(ctx, pollerEvents)
	wg.Wait()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c

		log.Info("terminate signal received, exiting...")
		cancel()
	}()
}

func processEntity(
	ctx context.Context,
	bhClient bullhorn.Bullhorn,
	pollerEvents <-chan *kafka.EventWrapper,
	publisherEvents chan<- *kafka.EventWrapper, wg *sync.WaitGroup,
) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			close(publisherEvents)
			return
		case events, ok := <-pollerEvents:
			if !ok {
				return
			}

			bullhornEvent := events.Data.([]bullhorn.Event)

			for _, event := range bullhornEvent {
				entity, err := bhClient.FetchEntityChanges(event)
				if err != nil {
					log.WithError(err).Error("fetching entity changes")
				}

				publisherEvents <- &kafka.EventWrapper{
					Topic: entity.Name,
					Data:  entity.Changes,
				}
			}
		}
	}
}
