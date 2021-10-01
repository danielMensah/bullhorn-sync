package main

import (
	"context"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/danielMensah/bullhorn-sync-poc/internal/bullhorn"
	"github.com/danielMensah/bullhorn-sync-poc/internal/worker"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

const workerCount = 20

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	spec, exists := os.LookupEnv("CRON_SPEC")
	if !exists {
		log.Fatal("missing environment variable 'CRON_SPEC'")
	}

	c := cron.New()
	_, err := c.AddFunc(spec, func() {
		syncEntities(ctx)
	})
	if err != nil {
		log.WithError(err).Fatal("could not start cron job")
	}

	c.Run()
}

func syncEntities(ctx context.Context) {
	sess := session.Must(session.NewSession())
	svc := eventbridge.New(sess)

	bhClient := bullhorn.New("", "")
	events, err := bhClient.GetEvents()
	if err != nil {
		log.WithError(err).Error("getting entities")
	}

	records := make(chan bullhorn.Record)
	wg := &sync.WaitGroup{}

	wkr := worker.New(svc)
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go wkr.Run(ctx, records, wg)
	}

	for _, event := range events {
		go bhClient.FetchChanges(event, records)
	}

	log.Info("sync completed")
}
