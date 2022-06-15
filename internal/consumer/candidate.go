package consumer

import (
	"github.com/danielMensah/bullhorn-sync-poc/internal/bullhorn"
	"github.com/danielMensah/bullhorn-sync-poc/internal/kafka"
)

type CandidateConsumer struct {
}

func New(bhClient bullhorn.Client, queue kafka.Subscriber) {

}
