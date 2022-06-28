package wpool

import (
	"context"
	"sync"

	"github.com/danielMensah/bullhorn-sync-poc/internal/bullhorn"
	pb "github.com/danielMensah/bullhorn-sync-poc/internal/proto"
	"github.com/danielMensah/bullhorn-sync-poc/internal/storage"
	log "github.com/sirupsen/logrus"
)

type Pool struct {
	workersCount int
	jobs         chan Job
	results      chan Result
	Done         chan struct{}
}

func New(bh bullhorn.Bullhorn, workerCount int) *Pool {
	return &Pool{
		workersCount: workerCount,
		jobs:         make(chan Job, workerCount),
		results:      make(chan Result, workerCount),
		Done:         make(chan struct{}),
	}
}

func (wp *Pool) Run(ctx context.Context) {
	var wg sync.WaitGroup

	for i := 0; i < wp.workersCount; i++ {
		wg.Add(1)
		go wp.worker(ctx, &wg, wp.jobs, wp.results)
	}

	wg.Wait()
	close(wp.Done)
	close(wp.results)
}

func (wp *Pool) worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan Job, results chan<- Result) {
	defer wg.Done()

	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}
			// fan-in job execution multiplexing results into the results channel
			results <- job.execute(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (wp *Pool) Results() <-chan Result {
	return wp.results
}

func (wp *Pool) AddJob(events <-chan *pb.Event, storage storage.Storage) {
	for event := range events {
		var execFn ExecutionFn

		switch event.EntityEventType {
		case pb.EventType_INSERTED:
			execFn = storage.Store
		case pb.EventType_UPDATED:
			execFn = storage.Update
		default:
			log.Errorf("unsupported event type: %v", event.EntityEventType)
			continue
		}

		wp.jobs <- Job{
			ExecFn: execFn,
			Event:  event,
		}
	}
}
