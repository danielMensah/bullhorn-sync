package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/danielMensah/bullhorn-sync/internal/auth"
	"github.com/danielMensah/bullhorn-sync/internal/broker"
	"github.com/danielMensah/bullhorn-sync/internal/bullhorn"
	"github.com/danielMensah/bullhorn-sync/internal/config"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	cfg, err := config.LoadConfig()
	if err != nil {
		log.WithError(err).Fatal("loading config")
	}

	ctx, cancel := context.WithCancel(context.Background())
	go waitForCancel(cancel)

	client, err := auth.New(ctx, cfg.BullhornUsername(), cfg.BullhornPassword(), cfg.OauthConfig())
	if err != nil {
		log.WithError(err).Fatal("creating new auth")
	}

	bhClient := bullhorn.New(cfg.BullhornSubscriptionUrl(), cfg.BullhornEntityUrl(), client)

	consumer, err := broker.NewKafkaConsumer(cfg.KafkaAddress(), "groupID")
	if err != nil {
		log.WithError(err).Fatal("creating new consumer")
	}
	defer consumer.Close()

	producer, err := broker.NewKafkaProducer(cfg.KafkaAddress())
	if err != nil {
		log.WithError(err).Fatal("creating new producer")
	}
	defer producer.Close()

	go producer.MonitorEvents()

	entityChangesChan := make(chan *broker.EventWrapper)
	pollerChan := make(chan *broker.EventWrapper)

	wg := &sync.WaitGroup{}

	for i := 0; i < 20; i++ {
		wg.Add(2)
		go processEntity(ctx, bhClient, pollerChan, entityChangesChan, wg)
		go producer.Produce(ctx, entityChangesChan, wg)
	}

	consumer.Consume(ctx, "poller_events", pollerChan)
	wg.Wait()
}

func processEntity(
	ctx context.Context,
	bhClient bullhorn.Bullhorn,
	pollerChan <-chan *broker.EventWrapper,
	entityChangesChan chan<- *broker.EventWrapper,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			close(entityChangesChan)
			return
		case events, ok := <-pollerChan:
			if !ok {
				return
			}

			bullhornEvent := events.Event.([]bullhorn.Event)

			for _, event := range bullhornEvent {
				entity, err := bhClient.FetchEntityChanges(event)
				if err != nil {
					log.WithError(err).Error("fetching entity changes")
					return
				}

				entityChangesChan <- &broker.EventWrapper{
					Topic: strings.ToLower(entity.Name),
					Event: entity.Changes,
				}
			}
		}
	}
}

func waitForCancel(cancelFunc context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	cancelFunc()
	log.Info("terminate signal received, exiting...")
}
