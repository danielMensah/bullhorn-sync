package worker

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/danielMensah/bullhorn-sync-poc/internal/bullhorn"
	"github.com/danielMensah/bullhorn-sync-poc/internal/eventbus"
	"github.com/danielMensah/bullhorn-sync-poc/internal/queue"
	log "github.com/sirupsen/logrus"
)

type Worker struct {
	bus      *eventbus.Client
	bullhorn *bullhorn.Client
}

func New(bus *eventbus.Client, bullhorn *bullhorn.Client) *Worker {
	return &Worker{bus: bus, bullhorn: bullhorn}
}

func (w *Worker) Run(ctx context.Context, messages <-chan queue.Message, wg *sync.WaitGroup) {
	defer wg.Done()

	for message := range messages {
		logger := log.WithField("eventBridge_id", message.Event.ID)
		event := bullhorn.Event{}

		if err := json.Unmarshal(message.Event.Detail, &event); err != nil {
			logger.WithError(err).Error("decoding event detail")
			continue
		}

		logger = logger.WithFields(log.Fields{
			"entityID":        event.EntityId,
			"entityName":      event.EntityName,
			"entityEventType": event.EntityEventType,
		})
		logger.Info("received event")

		record, err := w.bullhorn.FetchEntityChanges(event)
		if err != nil {
			logger.WithError(err).Error("fetching entity changes")
			continue
		}

		// TODO send data to database

		err = w.bus.SendProcessedEvent(ctx, eventbus.ProcessedSource, record)
		if err != nil {
			logger.WithError(err).Error("sending processed event")
			continue
		}

		if err := message.Delete(); err != nil {
			logger.WithError(err).Error("deleting message during worker")
		}

		logger.WithField("duration", time.Since(message.Event.Time).Milliseconds())
		logger.Info("duration recorded")
	}
}
