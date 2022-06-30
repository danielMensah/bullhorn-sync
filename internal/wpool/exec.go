package wpool

import (
	"context"
	"sync"

	"github.com/danielMensah/bullhorn-sync-poc/internal/consumer"
	pb "github.com/danielMensah/bullhorn-sync-poc/internal/proto"
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

func (wp *Pool) AddJob(entities <-chan *pb.Entity, consumer *consumer.Consumer) {
	for entity := range entities {
		var execFn ExecutionFn

		switch entity.Name {
		case "Candidate":
			execFn = consumer.ConsumeCandidate
		case "Company":
			execFn = consumer.ConsumeCompany
		default:
			log.Errorf("unsupported entity type: %v", entity.EventType)
			continue
		}

		wp.jobs <- Job{
			ExecFn: execFn,
			Entity: entity,
		}
	}
}
