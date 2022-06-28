package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/danielMensah/bullhorn-sync-poc/internal/bullhorn"
	"github.com/danielMensah/bullhorn-sync-poc/internal/config"
	"github.com/danielMensah/bullhorn-sync-poc/internal/kafka"
	pb "github.com/danielMensah/bullhorn-sync-poc/internal/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) Publish(ctx context.Context, eventWrapper *pb.EventsWrapper) (*emptypb.Empty, error) {
	log.SetFormatter(&log.JSONFormatter{})

	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("cannot load config: %w", err)
	}

	messenger, err := kafka.NewPublisher(ctx, cfg.KafkaAddress)
	if err != nil {
		return nil, fmt.Errorf("cannot create new publisher: %w", err)
	}

	bhClient, err := bullhorn.New(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot create new bullhorn instance: %w", err)
	}

	record := make(chan *bullhorn.Entity)
	wg := &sync.WaitGroup{}

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go messenger.Publish(ctx, record, wg)
	}

	for _, event := range eventWrapper.Events {
		entity, err := bhClient.FetchEntityChanges(event)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve entity changes: %w", err)
		}

		record <- entity
	}

	wg.Wait()
	close(record)

	return nil, nil
}
