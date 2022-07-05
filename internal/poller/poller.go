package poller

import (
	"context"
	"time"

	"github.com/danielMensah/bullhorn-sync/internal/broker"
	"github.com/danielMensah/bullhorn-sync/internal/bullhorn"
	log "github.com/sirupsen/logrus"
)

// Poller is a poller that fetches events from bullhorn and sends them to the broker
type Poller struct {
	bh bullhorn.Bullhorn
}

// New returns a new poller
func New(bhClient bullhorn.Bullhorn) *Poller {
	return &Poller{
		bh: bhClient,
	}
}

// Run fetches events from bullhorn and sends them back in the events channel
func (p *Poller) Run(ctx context.Context, events chan<- *broker.EventWrapper) {
	for {
		select {
		case <-ctx.Done():
			close(events)
			return
		default:
			fetchedEvents, err := p.bh.GetEvents()
			if err != nil {
				log.WithError(err).Error("getting entities")
			}

			events <- &broker.EventWrapper{
				Topic: "poller_events",
				Event: fetchedEvents,
			}

			time.Sleep(10 * time.Second)
		}
	}
}
