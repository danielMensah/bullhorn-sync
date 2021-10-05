package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/danielMensah/bullhorn-sync-poc/internal/bullhorn"
	log "github.com/sirupsen/logrus"
)

// Worker is responsible for doing the work to migrate data
type Worker struct {
	eventbusSvc *eventbridge.EventBridge
}

// New creates a new worker
func New(svc *eventbridge.EventBridge) *Worker {
	return &Worker{eventbusSvc: svc}
}

// Run is responsible for processing jobs
func (w *Worker) Run(ctx context.Context, records <-chan bullhorn.Record, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			break
		case r, ok := <-records:
			if !ok {
				return
			}

			logger := log.WithFields(log.Fields{
				"entityName":      r.EntityName,
				"entityId":        r.EntityId,
				"entityEventType": r.EntityEventType,
			})
			logger.Info("processing record")

			payload, err := json.Marshal(r)
			if err != nil {
				logger.WithError(err).Error("error marshalling data")
				return
			}

			res, err := w.eventbusSvc.PutEventsWithContext(ctx, &eventbridge.PutEventsInput{
				Entries: []*eventbridge.PutEventsRequestEntry{
					{
						EventBusName: aws.String("bullhorn-sync"),
						Source:       aws.String("sync.ingestion"),
						DetailType:   aws.String(r.EntityEventType),
						Detail:       aws.String(string(payload)),
						Time:         &r.EventTimestamp,
					},
				},
			})

			if err != nil {
				logger.WithError(err).Error("ingestion error")
				return
			}

			fmt.Println(res)
			//logger.Infof("job completed: %s", messageId)
		}
	}
}
