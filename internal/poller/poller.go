package poller

import (
	"context"
	"time"

	"github.com/danielMensah/bullhorn-sync-poc/internal/bullhorn"
	log "github.com/sirupsen/logrus"
)

type Poller struct {
	bh bullhorn.Bullhorn
}

func New(bhClient bullhorn.Bullhorn) *Poller {
	return &Poller{
		bh: bhClient,
	}
}

func (p *Poller) Run(ctx context.Context, events chan<- interface{}) {
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

			events <- fetchedEvents

			time.Sleep(10 * time.Second)
		}
	}
}
