package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/danielMensah/bullhorn-sync/internal/broker"
	log "github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go waitForCancel(cancel)

	consumer, err := broker.NewKafkaConsumer("localhost:9092", "groupID")
	if err != nil {
		log.WithError(err).Fatal("creating new consumer")
	}
	defer consumer.Close()

	entities := make(chan *broker.EventWrapper)
	wg := &sync.WaitGroup{}

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case entity, ok := <-entities:
					if !ok {
						return
					}

					fmt.Println(entity.Event)
				}
			}
		}()
	}

	// subscribe to needed topics
	go consumer.Consume(ctx, "candidate", entities)
	go consumer.Consume(ctx, "company", entities)
	wg.Wait()
}

func waitForCancel(cancelFunc context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	cancelFunc()
	log.Info("terminate signal received, exiting...")
}
