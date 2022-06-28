package wpool

//import (
//	"context"
//	"sync"
//
//	"github.com/danielMensah/bullhorn-sync-poc/internal/bullhorn"
//	"github.com/danielMensah/bullhorn-sync-poc/internal/kafka"
//	pb "github.com/danielMensah/bullhorn-sync-poc/internal/proto"
//	"github.com/danielMensah/bullhorn-sync-poc/internal/storage"
//	log "github.com/sirupsen/logrus"
//)
//
//type Worker struct {
//	bullhorn bullhorn.Bullhorn
//	consumer kafka.Consumer
//	storage  storage.Storage
//}
//
//func New(bullhorn bullhorn.Bullhorn, consumer kafka.Consumer, storage storage.Storage) *Worker {
//	return &Worker{
//		bullhorn: bullhorn,
//		consumer: consumer,
//		storage:  storage,
//	}
//}
//
//func (w *Worker) Run(ctx context.Context, event <-chan *pb.Event, wg *sync.WaitGroup) {
//	defer wg.Done()
//
//	for {
//		select {
//		case <-ctx.Done():
//			return
//		case e, ok := <-event:
//			if !ok {
//				return
//			}
//
//			_, err := w.bullhorn.FetchEntityChanges(e)
//			if err != nil {
//				log.WithError(err).Error("failed to fetch entity changes")
//				continue
//			}
//
//			// TODO: send to cassandra
//		}
//	}
//}
