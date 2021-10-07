package main

import (
	"context"

	"github.com/danielMensah/bullhorn-sync-poc/internal/bullhorn"
	"github.com/danielMensah/bullhorn-sync-poc/internal/config"
	"github.com/danielMensah/bullhorn-sync-poc/internal/eventbus"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.New()
	if err != nil {
		log.WithError(err).Fatal("getting configs")
	}

	c := cron.New()
	_, err = c.AddFunc(cfg.CronSpec, func() {
		ingest(ctx, cfg)
	})
	if err != nil {
		log.WithError(err).Fatal("could not start cron job")
	}

	c.Run()
}

func ingest(ctx context.Context, cfg config.Config) {
	bhClient, err := bullhorn.New(ctx, cfg.BhConfig, cfg.Oauth2Config)
	if err != nil {
		log.WithError(err).Fatal("new bullhorn client")
	}

	events, err := bhClient.GetEvents()
	if err != nil {
		log.WithError(err).Fatal("getting entities")
	}

	bus := eventbus.New(cfg.Region)
	err = bus.SendIngestedEvents(ctx, eventbus.IngestionSource, events)
	if err != nil {
		log.WithError(err).Error("bus")
	}

	log.Info("sync completed")
}
