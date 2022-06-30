package wpool

import (
	"context"
	"sync"

	"github.com/danielMensah/bullhorn-sync-poc/internal/broker"
	"github.com/danielMensah/bullhorn-sync-poc/internal/bullhorn"
	"github.com/danielMensah/bullhorn-sync-poc/internal/consumer"
	log "github.com/sirupsen/logrus"
)

type Pool struct {
	workersCount int
	jobs         chan Job
	results      chan Result
	Done         chan struct{}
}

func New(workerCount int) *Pool {
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

func (wp *Pool) AddJob(ctx context.Context, events <-chan *broker.EventWrapper, consumer *consumer.Client, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-events:
			if !ok {
				return
			}

			bullhornEntity := event.Data.(*bullhorn.Entity)

			var execFn ExecutionFn

			switch bullhornEntity.Name {
			case "Candidate":
				execFn = consumer.ConsumeCandidate
			case "Company":
				execFn = consumer.ConsumeCompany
			default:
				log.Errorf("unsupported entity type: %v", event.Topic)
				continue
			}

			wp.jobs <- Job{
				ExecFn: execFn,
				Entity: bullhornEntity,
			}
		}
	}
}
