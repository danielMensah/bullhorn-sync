package poller

import (
	"time"

	"github.com/danielMensah/bullhorn-sync-poc/internal/broker"
	"github.com/danielMensah/bullhorn-sync-poc/internal/bullhorn"
	log "github.com/sirupsen/logrus"
)

type Poller struct {
	bh   bullhorn.Bullhorn
	Done chan bool
}

func New(bhClient bullhorn.Bullhorn) *Poller {
	return &Poller{
		bh:   bhClient,
		Done: make(chan bool),
	}
}

func (p *Poller) Run(events chan<- *broker.EventWrapper) {
	for {
		select {
		case <-p.Done:
			close(events)
			return
		default:
			fetchedEvents, err := p.bh.GetEvents()
			if err != nil {
				log.WithError(err).Error("getting entities")
			}

			events <- &broker.EventWrapper{
				Topic: "poll_event",
				Data:  fetchedEvents,
			}

			time.Sleep(10 * time.Second)
		}
	}
}
